// implementation of bracket in go
// bracket is a concatenative programming language geared towards genetic
// programming
package main

import (
    //"errors"
    "fmt"
)


const cells = 100*1024*1024
const gcMargin  = cells - 24
const stackSize = 1024*1024

// Tagbits (from right  to left)
// local pointer, global pointer, Int, Prim, Symbol, Float
// global pointers not used yet, will become meta gene pool of gene bracket
// Bit 1 = 0 ->  Cons
//    Bit 4 = 0 --> pointer to cell on local heap
//    Bit 4 = 1 --> pointer to cell on global heap
// Bit 2 = 1 ->  Number
//    Bit 4 = 0 --> Int
//    Bit 4 = 1 --> Float
// Bit 3 = 1 ->  Symbol
//    Bit 4 = 0 --> assignable symbol
//    Bit 4 = 1 --> primitive

const tagType   = 15 // mask with bits 1111
const tagLocal  = 0  // bits 0000    cell on local heap
const tagGlobal = 8  // bits 1000    cell on global heap
const tagPrim   = 1  // bits 0001
const tagSymb   = 9  // bits 1001
const tagInt    = 3  // bits 0011
const tagFloat  = 11 // bits 1011

type value int

func boxCell(x int) value {return value(x<<4) }   // since taglocal=0
func boxGlobal(x int) value {return value(x<<4 | tagGlobal)}
func boxPrim(x int) value {return value(x<<4 | tagPrim)}
func boxSymb(x int) value {return value(x<<4 | tagSymb)}
func boxInt(x int)  value {return value(x<<4 | tagInt)}
//func box_float(x Int) Int { Int(reinterpret(Int32,x)) << 32 | tagFloat

func unbox(x value) int  {return int(x)>>4}   // remove all tags

// in contrast to C, here the pointer is
// just the heap index, that is, a number
func ptr(x value) int    {return int(x)>>4}   
//func unbox_float(x value) float = reinterpret(Float32, Int32(x>>32))

func isInt(x value)    bool {return (x & tagType == tagInt)}
func isFloat(x value)  bool {return (x & tagType == tagFloat)}
func isPrim(x value)   bool {return (x & tagType == tagPrim)}
func isSymb(x value)   bool {return (x & tagType == tagSymb)}
func isLocal(x value)  bool {return (x & tagType == tagLocal)}
func isGlobal(x value) bool {return (x & tagType == tagGlobal)}

func isCons(x value) bool {return (x & 1) == 0}
func isAtom(x value) bool {return !isCons(x)}
func isNumb(x value) bool {return (x & 3) == 3}
func isAbstractSymb(x value) bool {return (x & 3) == 1}

func isNil(x value) bool {return x == nill}
func isDef(x value) bool {return x != nill}

func (vm *Vm) isCons2(x value) bool {
    return isCons(x) && isCons(vm.cdr(x))
}
//isCons3(x,vm) = isCons(x) && isCons(cdr(x,vm)) && isCons(cddr(x,vm))


type cell struct {
    car  value
    cdr  value
}

const (  // bracket primitives
        nill value = iota<<4 | tagPrim   //  "nill" since "nil" already taken by golang
        dup
        drop
        swap
        cons
        car
        cdr
        eval
        rec
        def
        whl
        add
        sub
        mul
        div
        gt
        lt
        eq
        iff   // "iff" since "if" already taken by golang
        esc
        vesc
        val
        trace
        unbound
)

var primStr = map[value] string {
    cons:"cons", car:"car", cdr:"cdr", def:"def", dup:"dup", drop:"drop", 
    esc:"esc", eval:"eval", eq:"eq", iff:"if", 
    rec:"rec", swap:"swap", val:"val", vesc:"vesc", whl:"whl",
    add:"+", sub:"-", mul:"*", div:"/", gt:">", lt:"<",trace:"trace",
}

var str2prim = map[string] value {
    "cons":cons, "car":car, "cdr":cdr, "def":def, "dup":dup, "drop":drop, 
    "esc":esc, "eval":eval, "eq":eq, "if":iff, 
    "rec":rec, "swap":swap, "val":val, "vesc":vesc, "whl":whl,
    "add":add, "+":add, "sub":sub, "-":sub, "*":mul, "mul":mul, "/":div, "div":div,
    "gt":gt, ">":gt, "lt":lt, "<":lt, "trace":trace,
}

type stats struct { // some statistics about the running program
    nInst   int   // number of Instructions
    nRecur  int   // recursion depth
    nSteps  int   // no of performed programming steps
    extent  int   // exent of genome at birth (vm size)
}

type Vm struct {
    bra  value    // program, future of computation
    ket  value    // result, past of computation
    aux  value    // auxilliary stack (helps with stack shuffling)
    env  value    // environment
    root value    // anchor to save stuff for gc
    next int      // index to next entry on arena
    arena []cell  // memory arena to hold the cells
    brena []cell  // second arena, needed for copying gc
    stack []value  //
    stackIndex int
    needGc bool  // flag to indicate that heap space gets rare
    depth int     // current recursion depth
    stats stats   // some statistics about the running program
    trace int   //trace mode e: 0=no trace, 1=trace non-verbose, 3=verbose

}

func (vm *Vm) reset() {
    vm.next = 0
    vm.stats = stats{0,0,0,0}
    vm.bra = nill
    vm.ket = nill
    vm.aux = nill
    vm.env = vm.cons(nill,nill)
    vm.root = nill
    vm.stackIndex = 0
    vm.depth = 0
    vm.trace = 0
}

func init_vm() Vm {
    a := make([]cell, cells)
    b := make([]cell, cells)
    stack := make([]value, stackSize)
    stats := stats{0,0,0,0}
    vm := Vm{nill,nill,nill,nill,nill,0,a,b,stack,0,false,0,stats,0}
    vm.env = vm.cons(nill,nill)
    return vm 
}

//  garbage collector  *********************************
//  implement Cheney copying algorithm
//    Cheney :       non-recursive traversal of live-objects
func (vm *Vm) relocate(c value) value {
   if !isCons(c) {
       return c
   }
   indv := ptr(c)
   ah := vm.brena[indv]
   if ah.car == unbound {
       return ah.cdr
   }
   ind := vm.next
   vm.arena[ind]   = vm.brena[indv]
   vm.brena[indv] = cell{unbound, boxCell(ind)}
   vm.next += 1
   return boxCell(ind)
}

func (vm *Vm) gc() {
   fmt.Println("starting gc ************************************************")
   var c cell
   vm.brena, vm.arena = vm.arena, vm.brena
   finger :=  0
   vm.next = 0

   // scan root of every live object
   vm.bra = vm.relocate(vm.bra)
   vm.ket = vm.relocate(vm.ket)
   vm.root = vm.relocate(vm.root)
   //vm.aux = vm.relocate(vm.aux)
   vm.env = vm.relocate(vm.env)
   for i:=0; i<=vm.stackIndex; i++ { 
     vm.stack[i] = vm.relocate(vm.stack[i])
   }

   // scan remaining objects in arena (including objects added by this loop)
   for finger < vm.next {
      c = vm.arena[finger]
      vm.arena[finger] = cell{vm.relocate(c.car), vm.relocate(c.cdr)}
      finger += 1
  }

   //println("GC: live objects found: ", vm.next-1)

   if vm.next >= gcMargin {
       fmt.Println("Bracket GC, arena too small")
   }
   vm.needGc = false
   fmt.Println("GC finished")
}

// **********************


func (vm *Vm) cons(pcar, pcdr value) value {
   vm.next += 1
   if vm.next > gcMargin {
     vm.needGc = true
   }
   vm.arena[vm.next] = cell{pcar,pcdr}
   return boxCell(vm.next)  // return a boxed index
}

// pop top element from list
func (vm *Vm) pop(list, p *value) bool {
    if isCons(*list) {
        c := vm.arena[unbox(*list)]
        *p = c.car
        *list = c.cdr
        return true
    } else {  // pop from atom returns the atom and false
        *p = *list
        return false
    }
}

// Pop first two items in list 
func (vm *Vm) pop2(list, p1, p2 *value) bool {
    return vm.pop(list, p1) && vm.pop(list, p2);  
}

// unsafe, assumes p is a Cons
func (vm *Vm) car (p value) value {
    return vm.arena[p>>4].car
}

// unsafe, assumes p is a Cons
func (vm *Vm) cdr (p value) value {
    return vm.arena[p>>4].cdr
}

func (vm *Vm) caar(p value) value { 
    return vm.car(vm.car(p))
}
//func cadr(p Ptr, vm Vm) value { 
//    return car(cdr(p,vm).(Ptr),vm)
//}
//func cddr(p Ptr, vm Vm) value { 
//    return cdr(cdr(p,vm).(Ptr),vm)
//}

// reverse a list
// if list contained a dotted pair, reverse returns normal list
// but also a flag 
func (vm *Vm) reverse(list value) (value, bool) {
    var p value 
    l := nill
    for vm.pop(&list,&p) {
       l = vm.cons(p,l)
    }
    if isDef(list) {   // list contained a dotted pais
        l = vm.cons(list,l)
        return l, true
    } else {
        return l, false // list did not contain a dotted pair
    }
}

// just count the number of conses, ie dotted pair has length 1
func (vm *Vm) length(list value) int {
   n := 0
   for isCons(list) {
       n += 1
       list = vm.cdr(list)
   }
   return n
}

func (vm *Vm) isEqual(p1, p2 value) bool {
   if isCons(p1) && isCons(p2) { 
      return (vm.isEqual(vm.car(p1),vm.car(p2)) && 
              vm.isEqual(vm.cdr(p1),vm.cdr(p2)))
   } else { 
       return (p1 == p2)
   }
}


// stack functions
func (vm *Vm) pushStack(x value) {
    if vm.stackIndex == stackSize {
        panic("Vm stack overflow")
    }
    vm.stack[vm.stackIndex] = x
    vm.stackIndex++;
}

// we don't check for underflow because it should never occur
func (vm *Vm) popStack() value {
    //if vm.stackIndex == 0 {
    //    return nill, errors.New("Vm stack underflow")
    //}
    vm.stackIndex--
    x := vm.stack[vm.stackIndex]
    return x
}

func (vm *Vm) getStack() value {
    //if vm.stackIndex == 0 {
    //    return nill, errors.New("Vm stack underflow")
    //}
    x := vm.stack[vm.stackIndex]
    return x
}

func (vm *Vm) replaceStack(x value) {
    vm.stack[vm.stackIndex] = x
}

func (vm *Vm) newEnv(env value) value {
    return vm.cons(nill,env)
}

func (vm *Vm) findKey(key, env value) value {
    bnds := vm.car(vm.env)  //current frame (list of bindings) is on top of env
    for isCons(bnds) {
       bnd := vm.car(bnds)
       if vm.car(bnd) == key {
           return bnd
       } else {
          bnds = vm.cdr(bnds) 
       }
    }
    return nill
}
          
func (vm *Vm) boundvalue(key value) value { // lookup symbol.. 
    bnd := vm.binding(key, vm.env)
    if bnd == nill {
        return nill
    } else {
        return vm.cdr(bnd)
    }
}      


func (vm *Vm) binding(key, env value) value {
    bnd := vm.findKey(key, env)
    for isNil(env) {
       env = vm.cdr(env)
       if isNil(env) {
          return nill
       }
       bnd = vm.findKey(key, env)
    }
    return bnd
}


func (vm *Vm) bindKey(key,val value) {
    var b,f,frame value
    f = nill
    vm.pop(&vm.env,&frame)   // pop frame from env
    for vm.pop(&frame,&b) {  // search frame for key
       if vm.car(b) == key {
          break
       } else {
          f = vm.cons(b,f)   // store binding in aux frame f
       }
    }
    for vm.pop(&f,&b) {    // put other bindings back in place
       frame = vm.cons(b,frame)
    }
    b1 := vm.cons(key, val)  // and add new binding
    frame = vm.cons(b1, frame)
    vm.env = vm.cons(frame, vm.env)
}

func (vm *Vm) replaceKeyInFrame(key, val, frame value) (newFrame value) {
    f := nill
    var b value
    for vm.pop(&frame,&b) {  // search for key
       if vm.car(b) == key {
          b1 := vm.cons(key, val)  // new binding
          frame = vm.cons(b1, frame)
          break
       } else {
          f = vm.cons(b,f)
      }
    }
    if frame == nill {  // key not found
        return nill
    } 
    for vm.pop(&f,&b) { // put other bindings back in place
       frame = vm.cons(b,frame)
    }
    return frame
}

func (vm *Vm) replaceKey(key, val, env value) value {
    e := nill  // empty envirohnment
    var f value
    envSave := env
    for vm.pop(&env,&f) {  // search all frames in environment
       f1 := vm.replaceKeyInFrame(key, val, f)
       if isDef(f1) { // key found
          env = vm.cons(f1, env) 
          vm.printKet(env)
          break
       } else {
          e = vm.cons(f,e)
       }
    } 
    if isNil(env) { // key not found
        b := vm.cons(key,val)
        vm.pop(&envSave,&f)
        f := vm.cons(b,f)
        return vm.cons(f,env)
    } 
    
    for vm.pop(&e,&f) {  // put other frames back in place
      env = vm.cons(f, env)
    }
    return env
}

func istrue(x value) bool {
    switch {
    case isNumb(x):
          return unbox(x) != 0
    case isPrim(x):
          return x != nill
    default:
          return true
    }
}



// *******************************************
func (vm *Vm) fDup() {
   if isCons(vm.ket) {
       vm.ket = vm.cons(vm.car(vm.ket), vm.ket)
   }
}

func (vm *Vm) fDrop() {
   if isCons(vm.ket) {
       vm.ket = vm.cdr(vm.ket)
   }
}

func (vm *Vm) fSwap() {
   var a,b value 
   if vm.pop2(&vm.ket, &a, &b) {
       vm.ket = vm.cons(a,vm.ket)
       vm.ket = vm.cons(b,vm.ket)
   }
}

func (vm *Vm) fCons() {
    var p1,p2 value
    if vm.pop2(&vm.ket, &p1, &p2) {
      vm.ket = vm.cons(vm.cons(p1,p2),vm.ket)
    }
}

func (vm *Vm) fCar() {
    var head, p value
    if vm.pop(&vm.ket, &p) {
        if vm.pop(&p, &head) { // car a list
            vm.ket = vm.cons(p,vm.ket) // leave the rest of the list on the ket 
            vm.ket = vm.cons(head,vm.ket)
        } else {    // car a symbol
          val := vm.boundvalue(p)   // lookup symbol.. 
          if isCons(val) {
              vm.ket = vm.cons(vm.car(val),vm.ket)
          }
        }
    }
}

func (vm *Vm) fCdr() {
    var head,p value
    if vm.pop(&vm.ket, &p) {
        if vm.pop(&p, &head) { // cdr a list
            vm.ket = vm.cons(p,vm.ket) 
        } else {
          val := vm.boundvalue(p)   // lookup symbol.. 
          if isCons(val) {
              vm.ket = vm.cons(vm.cdr(val),vm.ket)
          } else {
              vm.ket = vm.cons(nill,vm.ket)  // at least leave a nill on ket
          }
        }
    }
}

type mathIntFunc func(int, int) int
func myAdd(x,y int) int {return x+y}
func mySub(x,y int) int {return x-y}
func myMul(x,y int) int {return x*y}
func myDiv(x,y int) int {return x/y}

func myGt(x,y int) int {
   if x>y {
       return 1
   } else  {
       return 0
   }
}
func myLt(x,y int) int {
   if y>x {
       return 1
   } else  {
       return 0
   }
}

func (vm *Vm) fMath(op mathIntFunc) {
    var n1,n2 value
    if vm.pop2(&vm.ket, &n1, &n2) {
      n3 := boxInt(op(unbox(n1), unbox(n2)))
      vm.ket = vm.cons(n3,vm.ket)
      //fmt.Println("add: ",n3)
  }
}

func (vm *Vm) fEq() {
   var p1, p2 value
   if vm.pop2(&vm.ket, &p1, &p2) {
       b := boxInt(0)
       if vm.isEqual(p1,p2) {b = boxInt(1)}
       vm.ket = vm.cons(b, vm.ket)
   }
}

func (vm *Vm) fIf() {
   var p, p1, p2 value
   if vm.pop2(&vm.ket, &p1, &p2) && vm.pop(&vm.ket,&p) {
       if istrue(p) {
          vm.ket = vm.cons(p1, vm.ket)
      } else {
          vm.ket = vm.cons(p2, vm.ket)
      }
   }
}


func (vm *Vm) fEsc() {
    var val value
    if vm.pop(&vm.bra, &val) { 
         vm.ket = vm.cons(val,vm.ket)
    }
}

func (vm *Vm) fVesc() {
    var val value
    if vm.pop(&vm.bra, &val) { 
         vm.ket = vm.cons(val,vm.ket)
         vm.fVal()
    }
}

func (vm *Vm) fVal() {
    var key value
    if vm.pop(&vm.ket, &key) { 
        if isCons(key) {
          vm.ket = vm.cons(key,vm.ket)
        } else {
          val := vm.boundvalue(key)       // lookup symbol.. 
          vm.ket = vm.cons(val,vm.ket)    // .. and place on ket
        }
    }
}

func (vm *Vm) fTrace() { // change trace mode
    var p value
    if vm.pop(&vm.ket,&p) {
        vm.trace = unbox(p)
    }
}

func (vm *Vm) fEval() {
    var op value
    if vm.pop(&vm.ket,&op){
        switch {
        case isCons(op):
            vm.evalCons(op)
        case isNil(op):
            return
        case isPrim(op):
             vm.evalPrim(op)
        case isSymb(op):
             vm.evalSymb(op) 
        default:   // eval a number
             vm.evalNumb(op)
        }
    }
}

func (vm *Vm) fRec() {
//anonymous recursion: replace bra of this scope by original value
    var b value
    if vm.pop(&vm.ket,&b) { // pop a boolean value
        if istrue(b) {
            vm.bra = vm.getStack()
        }
    }
}

func (vm *Vm) fDef() {
   var key, val value
   if vm.pop2(&vm.ket, &key, &val) {
       if !isCons(key) {
         vm.bindKey(key,val)  // bind key to val in top env-frame
       }
       //vm.printList(vm.env); fmt.Println()
       //vm.printList(vm.car(vm.env)); fmt.Println()
   }
}

func (vm *Vm) fWhl() {
   fmt.Println("Start whl ")
   var q,b value
   if vm.pop2(&vm.ket,&q, &b) {
      //vm->root = cons(vm->bra,vm->root,vm);
      //vm->root = cons(q,vm->root,vm);
       vm.root = vm.cons(vm.bra,vm.root)
       vm.root = vm.cons(q,vm.root)
       for istrue(b) {
         vm.bra = vm.car(vm.root)
         vm.evalBra()
         if !vm.pop(&vm.ket,&b){
             break
         }
       }
       vm.root = vm.cdr(vm.root)
       vm.bra = vm.car(vm.root)
       vm.root = vm.cdr(vm.root)
   }
}

func (vm *Vm) evalCons(op value) {
    if isCons(vm.bra) {
       vm.depth++
       vm.pushStack(vm.env)
       vm.pushStack(vm.bra)
       vm.pushStack(op)
       vm.env = vm.newEnv(vm.env)
    } else { // tail position
       vm.replaceStack(op)
    }
    vm.bra = op
}

func (vm *Vm) evalNumb(n value) {
    vm.ket = vm.cons(n,vm.ket)
}


func (vm *Vm) evalSymb(sym value) {
    bnd := vm.binding(sym,vm.env)
    if isNil(bnd) {
       vm.ket = vm.cons(nill,vm.ket)
    } else {
       val := vm.cdr(bnd)
       if isCons(val) {
          vm.evalCons(val)
       } else {
          vm.ket = vm.cons(val,vm.ket)
       }
    }

}

func (vm *Vm) evalPrim(p value) {
    switch p { 
    case dup:
        vm.fDup()
    case drop:
        vm.fDrop()
    case swap:
        vm.fSwap()
    case cons:
        vm.fCons()
    case car:
        vm.fCar()
    case cdr:
        vm.fCdr()
    case eval:
        vm.fEval()
    case rec:
        vm.fRec()
    case def:
        vm.fDef()
    case whl:
        vm.fWhl()
    case add:
        vm.fMath(myAdd)
    case sub:
        vm.fMath(mySub)
    case mul:
        vm.fMath(myMul)
    case div:
        vm.fMath(myDiv)
    case gt:
        vm.fMath(myGt)
    case lt:
        vm.fMath(myLt)
    case eq:
        vm.fEq()
    case iff:
        vm.fIf()
    case esc:
        vm.fEsc()
    case vesc:
        vm.fVesc()
    case val:
        vm.fVal()
    case trace:
        vm.fTrace()
    default:
        fmt.Println("Error: unknown primitive")
        
    }
}

func (vm *Vm) evalBra() {
      //fmt.Println("Start eval ")
    startingDepth := vm.depth
    vm.pushStack(vm.bra)
    var e value
    for {
        if vm.trace > 0 { 
            //fmt.Println("trace")
            //vm.printBra(vm.bra)
            vm.printKet(vm.ket)
            //vm.printBra(vm.env)
            fmt.Println()
        }
        vm.pop(&vm.bra,&e);
        //fmt.Println("e=",e)
       
        switch {
        case isNil(e):
           vm.ket = vm.cons(e,vm.ket)
        case isPrim(e):
           vm.evalPrim(e)
        case isSymb(e):
           vm.evalSymb(e)
        default:
           vm.ket = vm.cons(e,vm.ket)
        }

        if vm.needGc {
            vm.gc()
        }
        if isNil(vm.bra) {    // exit scope
            if vm.depth == startingDepth {
              break
            }
            vm.depth--
            _ = vm.popStack()  // for rec
            vm.bra = vm.popStack()
            vm.env = vm.popStack()
        }
    }
    vm.bra = vm.popStack()
}

func main() {
    fmt.Printf("rock'n roll\n")   
    vm := init_vm()
    
    //prog := "add 1 2"
    //prog := "whl [gt 0 dup add 1] 1 -50000000"  // 5e7, 3 sec on Mac
    //prog := "foo bar foo def foo' 2 def bar' 3"
 
    //prog := "eval foo' 1 2 def foo' [add] trace 1"   
    
    //prog := "eval [ rec gt 0 dup add 1 ] -5 trace 1"
    prog := "eval [rec gt 0 dup add 1 dup] -5 trace 0"
    //prog := "eval [ rec gt 0 dup add 1 ] -50000000"   // 5e7, 3.9 sec on MAc
    //prog := "eval [ rec gt 0 dup add 1 ] -500000000"   // 5e8, 24.9 sec on MAc

    vm.bra = vm.makeBra(prog)
    
    //vm.bra = vm.loadFile("test.clj")
   
    vm.printBra(vm.bra)
    //vm.printList(vm.bra) 
    //vm.printList(vm.reverse(vm.bra))
    
    vm.evalBra()
    //fmt.Println(vm.bra)

    //vm.gc()

    vm.printKet(vm.ket)

}


