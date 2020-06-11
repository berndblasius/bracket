// implementation of bracket in go
// bracket is a concatenative programming language geared towards genetic
// programming
package main

import (
    //"errors"
    "fmt"
    "math/rand"
    "os"
)

const cells = 24*1024*1024
//const cells = 1024*1024
const gcMargin  = cells - 24
const stackSize = 1024*1024

// Tagbits (from right  to left)
// three bits are used (from Bit 1 to Bit 3), Bit 4 is free and can be used for gc for tree traversals
// local pointer, global pointer, Int, Prim, Symbol, Float
// global pointers not used yet, will become meta gene pool of gene bracket
// Bit 1 = 0 ->  Cell
//    Bit 2 = 0 --> pointer to cell on local heap
//    Bit 2 = 1 --> pointer to cell on global heap  (not used in the moment)
//    Bit 3 = 0 --> cons (ie, list or quotation)
//    Bit 3 = 1 --> closure
// Bit 1 = 1 ->  Number or Symb
// Bit 2 = 0 --> Symb
//    Bit 3 = 0 --> assignable symbol
//    Bit 3 = 1 --> primitive
// Bit 2 = 1 ->  Number
//    Bit 3 = 0 --> Int
//    Bit 3 = 1 --> Float

const tagType    = 7  // mask with bits 111
const tagGlobal  = 2  // bits 010    cell on global heap
const tagCell    = 1  // bits 001
const tagCons    = 5  // bits 101
const tagClosure = 4  // bits 100
const tagPrim    = 5  // bits 101
const tagSymb    = 1  // bits 001
const tagNumb    = 2  // bits 010
const tagInt     = 3  // bits 011
const tagFloat   = 7  // bits 111

type value int

func boxCons(x int) value {return value(x<<4) }   // create a new local cons
func boxClosure(x int) value {return value(x<<4 | tagClosure)}   // create a new local closure
//func boxGlobal(x int) value {return value(x<<4 | tagGlobal)}
func boxPrim(x int) value {return value(x<<4 | tagPrim)}  // create a local primitive
func boxSymb(x int) value {return value(x<<4 | tagSymb)}
func boxInt(x int)  value {return value(x<<4 | tagInt)}

// floats not yet interpreted
//func box_float(x Int) Int { Int(reinterpret(Int32,x)) << 32 | tagFloat
//func unbox_float(x value) float = reinterpret(Float32, Int32(x>>32))

func unbox(x value) int  {return int(x)>>4}   // remove all tags
// in contrast to C, here the pointer is
// just the heap index, that is, a number
//func ptr(x value) int    {return int(x)>>4}   

func isInt(x value)    bool {return (x & tagType == tagInt)}
func isFloat(x value)  bool {return (x & tagType == tagFloat)}
func isPrim(x value)   bool {return (x & tagType == tagPrim)}
func isSymb(x value)   bool {return (x & tagType == tagSymb)}
func isLocal(x value)  bool {return (x & tagGlobal == 0)}
func isGlobal(x value) bool {return (x & tagGlobal == tagGlobal)}
func isCell(x value)   bool {return (x & tagCell == 0)}
func isAtom(x value)   bool {return (x & tagCell == tagCell)}
//func isAtom(x value)  bool {return !isCell(x)}
func isCons(x value)   bool {return (x & tagCons == 0)}
func isClosure(x value) bool {return (x & tagCons == tagClosure)}
func isNumb(x value)  bool {return (x & tagNumb) == tagNumb}
func isAbstractSymb(x value) bool {return (x & tagNumb) == 0}  // symbol or primitive

func isNil(x value) bool {return x == nill}
func isDef(x value) bool {return x != nill}

func (vm *Vm) isCell2(x value) bool {
    return isCell(x) && isCell(vm.cdr(x))
}
//isCons3(x,vm) = isCell(x) && isCell(cdr(x,vm)) && isCell(cddr(x,vm))


type cell struct {
    car  value
    cdr  value
}

const (  // bracket primitives
        nill value = iota<<4 | tagPrim   //  "nill" since "nil" already taken by golang
        dup
        drop
        swap
        rot
        cons
        car
        cdr
        eval
        dip
        rec
        def
        lambda
        add
        sub
        mul
        div
        gt
        lt
        rnd
        eq
        iff   // "iff" since "if" already taken by golang
        esc
        vesc
        val
        trace
        typ
        print
        unbound
)
        //rto
        //tor
        //cond
        //set
        //whl

var primStr = map[value] string {
    cons:"cons", car:"car", cdr:"cdr", def:"def", dip:"dip", dup:"dup", drop:"drop", 
    esc:"esc", eval:"eval", eq:"eq", iff:"if",  lambda:"\\",
    rec:"rec", swap:"swap", val:"val", vesc:"vesc", 
    add:"+", sub:"-", mul:"*", div:"/", gt:">", lt:"<",rnd:"rnd",
    rot:"rot", trace:"trace", typ:"typ", print:"print",
}
//cond:"cond",set:"set",dip:"dip",whl:"whl",
//rto:"toR", tor:"Rto", 

var str2prim = map[string] value {
    "cons":cons, "car":car, "cdr":cdr, "def":def, "dip":dip, "dup":dup, "drop":drop, 
    "esc":esc, "eval":eval, "eq":eq, "if":iff, "\\":lambda, "lambda":lambda,
    "rec":rec,  "swap":swap, "val":val, "vesc":vesc, 
    "add":add, "+":add, "sub":sub, "-":sub, "*":mul, "mul":mul, "/":div, "div":div,
    "gt":gt, ">":gt, "lt":lt, "<":lt, "rnd":rnd,
    "rot":rot,"trace":trace,"typ":typ,"print":print,
}
//"cond":cond,"set":set,"dip":dip,"whl":whl,
//"toR":tor, "Rto":rto,

type stats struct { // some statistics about the running program
    nInst   int   // number of Instructions
    nRecur  int   // recursion depth
    nSteps  int   // no of performed programming steps
    extent  int   // exent of genome at birth (vm size)
}

// virtual machine
type Vm struct {
    bra  value    // program, future of computation
    ket  value    // global data stack, past of computation
    env  value    // environment
    next int      // index to next entry on arena
    arena []cell  // memory arena to hold the cells
    brena []cell  // second arena, needed for copying gc
    stack []value  //
    stackIndex int
    needGc bool   // flag to indicate that heap space gets rare
    depth int     // current recursion depth
    stats stats   // some statistics about the running program
    trace int     //trace mode e: 0=no trace, 1=trace non-verbose, 3=verbose

}

func init_vm() Vm {
    a := make([]cell, cells)
    b := make([]cell, cells)
    stack := make([]value, stackSize)
    stats := stats{0,0,0,0}
    vm := Vm{nill,nill,nill,-1,a,b,stack,-1,false,0,stats,0}
    vm.env = vm.cons(nill,nill)
    return vm 
}

func (vm *Vm) reset() {
    vm.next = -1
    vm.stats = stats{0,0,0,0}
    vm.bra = nill
    vm.ket = nill
    vm.env = vm.cons(nill,nill)
    vm.stackIndex = -1
    vm.depth = 0
    vm.trace = 0
    vm.needGc = false
}

//  garbage collector  *********************************
//  implement Cheney copying algorithm
//    Cheney :  non-recursive traversal of live-objects
func (vm *Vm) relocate(c value) value {
   var c1 value
   if !isCell(c) {
       return c
   }
   indb := unbox(c)   // index into brena
   bcell := vm.brena[indb]
   if bcell.car == unbound {
       return bcell.cdr
   }
   inda := vm.next    // index into arena
   if isCons(c) {
     c1 = boxCons(inda)
   } else {
     c1 = boxClosure(inda)
   }
   //vm.arena[ind]   = vm.brena[indv]
   vm.arena[inda]  = bcell
   vm.brena[indb] = cell{unbound, c1}
   vm.next += 1
   return c1
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
   vm.env = vm.relocate(vm.env)
   for i:=0; i<=vm.stackIndex; i++ { 
     vm.stack[i] = vm.relocate(vm.stack[i])
   }

   // scan remaining objects in arena (including objects added by this loop)
   for finger < vm.next {
      c = vm.arena[finger]
      rcar := vm.relocate(c.car)
      rcdr := vm.relocate(c.cdr)
      vm.arena[finger] = cell{rcar,rcdr}
      //vm.arena[finger] = cell{vm.relocate(c.car), vm.relocate(c.cdr)}
      finger += 1
  }

   //fmt.Println("GC: live objects found: ", vm.next-1)
   //fmt.Println("stack ", vm.stackIndex, " ", vm.depth)
   if vm.next >= gcMargin {
       fmt.Println("Bracket GC, arena too small")
        panic("Vm stack overflow")
   }
   vm.needGc = false
   //fmt.Println("GC finished")
}

// **********************

func (vm *Vm) makeCons(pcar, pcdr value) int {
   vm.next += 1
   if vm.next > gcMargin {
     vm.needGc = true
   }
   vm.arena[vm.next] = cell{pcar,pcdr}
   return vm.next  // return index
}

func (vm *Vm) cons(pcar, pcdr value) value {
    return boxCons(vm.makeCons(pcar,pcdr))
}

func (vm *Vm) closure(pcar, pcdr value) value {
    return boxClosure(vm.makeCons(pcar,pcdr))
}

func (vm *Vm) stripClosure(cl *value) {
    if isClosure(*cl) {
        *cl = vm.car(*cl)
    } 
}

// ------- careful that we do not leak mutablilty
//         should be used only for environments
// modify car or cdr of a cell without allocating a new cell
// should only be used for bindings
func (vm *Vm) setcar(cl value, newcar value) {
    ind := unbox(cl)
    vm.arena[ind].car = newcar 
    //pcdr := vm.arena[ind].cdr
    //vm.arena[ind] = cell{newcar, pcdr} 
}

func (vm *Vm) setcdr(cl value, newcdr value) {
    ind := unbox(cl)
    vm.arena[ind].cdr = newcdr 
    //pcar := vm.arena[ind].car
    //vm.arena[ind] = cell{pcar, newcdr} 
}

// -------------------------------------------------

// unsafe, assumes p is a Cell
func (vm *Vm) car (p value) value {
    return vm.arena[p>>4].car
}

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


// pop top element from list (also from closure)
func (vm *Vm) pop(list, p *value) bool {
    if isCell(*list) {
        c := vm.arena[unbox(*list)]
        *p = c.car
        *list = c.cdr
        return true
    } else {  // pop from atom returns the atom and false
        *p = *list
        return false
    }
}

// pop top element from only from Cons (i.e. not from a closure)
func (vm *Vm) popCons(list, p *value) bool {
    if isCons(*list) {
        c := vm.arena[unbox(*list)]
        *p = c.car
        *list = c.cdr
        return true
    } else {  // pop from atom or closure returns the atom and false
        *p = *list
        return false
    }
}

// Pop first two items in list 
func (vm *Vm) pop2(list, p1, p2 *value) bool {
    return vm.pop(list, p1) && vm.pop(list, p2);  
}

// just count the number of conses, ie dotted pair has length 1
func (vm *Vm) length(list value) int {
   n := 0
   for isCell(list) {
       n += 1
       list = vm.cdr(list)
   }
   return n
}

// list length, without quoted values (but also including dotted pairs) 
func (vm *Vm) lengthNonQuoted(list value) int {
   n := 0
   for isCell(list) {
       elem := vm.car(list)
       //if elem == esc {    // a quoted element
       //   if isDef(vm.cdr(list)) {  
       //      list = vm.cdr(list)
       //   }
       //} else 
       if (elem != vesc && elem != lambda) {  // no quoted element
          n += 1
       }
       list = vm.cdr(list)
   }
   if isDef(list) {  // count last element in dotted pair
       n += 1
   }
   return n
}

// reverse a list
// if list contained a dotted pair, reverse returns normal list
//      but also a flag 
// reverse only to first occurence of a closure, because
//   closures can occur only at end of a list
//   to avoid infinite loop when printing environments
func (vm *Vm) reverse(list value) (value, bool) {
    var p value 
    l := nill
    for vm.popCons(&list,&p) { // take care not to pop from a closure
        l = vm.cons(p,l)
        if vm.needGc {
            vm.pushStack(l)
            vm.pushStack(list)
            vm.gc()
            list = vm.popStack()
            l = vm.popStack()
        }
    }
    if isClosure(list) { //take only quotation from closure, not the env
       l = vm.cons(vm.car(list), l)
       return l, true
    } else if isDef(list) {   // list contained a dotted pais
        l = vm.cons(list,l)
        return l, true
    } else {
        return l, false // list did not contain a dotted pair
    }
}

func (vm *Vm) isEqual(p1, p2 value) bool {
   vm.stripClosure(&p1)
   vm.stripClosure(&p2)
   if isCell(p1) && isCell(p2) { 
      return (vm.isEqual(vm.car(p1),vm.car(p2)) && 
              vm.isEqual(vm.cdr(p1),vm.cdr(p2)))
   } else { 
       return (p1 == p2)
   }
}

// stack functions  ---------------------------
func (vm *Vm) pushStack(x value) {
    if vm.stackIndex == stackSize {
        panic("Vm stack overflow")
    }
    vm.stackIndex++;
    vm.stack[vm.stackIndex] = x
}

// we don't check for underflow because it should never occur
func (vm *Vm) popStack() value {
    //if vm.stackIndex == 0 {
    //    return nill, errors.New("Vm stack underflow")
    //}
    x := vm.stack[vm.stackIndex]
    vm.stackIndex--
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

func (vm *Vm) printStack() {
    fmt.Println("stack: ")
    for i:=0; i<vm.stackIndex; i++ {
        vm.printElem(vm.stack[i]); fmt.Println()
    }
    fmt.Println()
}


// creates new empty environment
func (vm *Vm) newEnv(env value) value {
    return vm.cons(nill,env)
}

// ----------- bindings -----------------------------------

func (vm *Vm) findLocalKey(key, env value) value {
// search binding with key in current (= top of env) frame 
    if isNil(key) {
        return nill
    }
    for bnds := vm.car(env); isCell(bnds); bnds = vm.cdr(bnds) {
       bnd := vm.car(bnds)
       if vm.car(bnd) == key {
           return bnd
       }
    }
    return nill
}

func (vm *Vm) findKey(key value) value {
// search binding with key in whole environment 
    if isNil(key) {
        return nill
    }
    for env:=vm.env; isDef(env); env = vm.cdr(env) {
       bnd := vm.findLocalKey(key, env)
       if isDef(bnd) {
           return bnd
       }
    }
    return nill
}

func (vm *Vm) boundvalue(key value) value { // lookup symbol.. 
    bnd := vm.findKey(key)
    if bnd == nill {
        return nill
    } else {
        return vm.cdr(bnd)
    }
}      

func (vm *Vm) bindKey(key,val value) {
// search for key in top frame, if key found override
// otherwise make new binding in top frame
    env := vm.env
    bnd := vm.findLocalKey(key,env)
    if isNil(bnd) { // key does not yet exist
        bnd = vm.cons(key,val)
        vm.setcar(env, vm.cons(bnd, vm.car(env)) )
    } else { // key exists, just override val
        vm.setcdr(bnd, val)
    }
}

func (vm *Vm) setKey(key,val value) {
// search for key in full environment, if key found override
// otherwise make new binding in top frame
    env := vm.env
    bnd := vm.findKey(key)
    if isNil(bnd) { // key does not yet exist
        bnd = vm.cons(key,val)
        vm.setcar(env, vm.cons(bnd, vm.car(env)) )
    } else { // key exists, just override val
        vm.setcdr(bnd, val)
    }
}


/*
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
*/

func istrue(x value) bool {
    return x != nill && unbox(x) != 0
}


// *******************************************
func (vm *Vm) fDup() {
   if isCell(vm.ket) {
       vm.ket = vm.cons(vm.car(vm.ket), vm.ket)
   }
}

func (vm *Vm) fDrop() {
   if isCell(vm.ket) {
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

func (vm *Vm) fRot() {
   var a,b,c value 
   if vm.pop2(&vm.ket, &a, &b) && vm.pop(&vm.ket,&c){
       vm.ket = vm.cons(b,vm.ket)
       vm.ket = vm.cons(a,vm.ket)
       vm.ket = vm.cons(c,vm.ket)
   }
}


func (vm *Vm) fCons() {
    var p1,p2 value
    if vm.pop2(&vm.ket, &p1, &p2) {
        vm.stripClosure(&p2) // cons to a closure strips the closure
        // instead here we could cons to the quotation of the closure
        vm.ket = vm.cons(vm.cons(p1,p2),vm.ket)
    }
}

func (vm *Vm) fCar() {
    var head, p value
    if vm.pop(&vm.ket, &p) {
        vm.stripClosure(&p)
        if vm.pop(&p, &head) { // car a list
            //vm.ket = vm.cons(p,vm.ket) // leave the rest of the list on the ket 
            vm.ket = vm.cons(head,vm.ket)
        } else {    // car a symbol
          val := vm.boundvalue(p)   // lookup symbol.. 
          vm.stripClosure(&val)
          if isCons(val) {
              vm.ket = vm.cons(vm.car(val),vm.ket)
          }
        }
    }
}

func (vm *Vm) fCdr() {
    var head,p value
    if vm.pop(&vm.ket, &p) {
        vm.stripClosure(&p)
        if vm.pop(&p, &head) { // cdr a list
            vm.ket = vm.cons(p,vm.ket) 
        } else {
          val := vm.boundvalue(p)   // lookup symbol.. 
          vm.stripClosure(&val)
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
func myDiv(x,y int) int {
    if y==0 {
        return 0
    } else {
       return x/y
    }
}

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

// a bit spagetti, but doing the job
func (vm *Vm) fMath(op mathIntFunc) {
    var c1, c2, n1, n2 value
    if vm.pop2(&vm.ket, &n1, &n2) {
      if isSymb(n1) {
          n1 = vm.boundvalue(n1) 
      }
      if isSymb(n2) {
          n2 = vm.boundvalue(n2) 
      }
      if isNumb(n1) && isNumb(n2) {
          vm.ket = vm.cons(boxInt(op(unbox(n1), unbox(n2))),vm.ket)
      } else if isCell(n1) && isCell(n2) {
          vm.stripClosure(&n1)
          vm.stripClosure(&n2)
          c := nill
          for vm.pop(&n1,&c1) && vm.pop(&n2,&c2) {
             if isSymb(c1) {
                    c1 = vm.boundvalue(c1) 
             }
             if isSymb(c2) {
                    c2 = vm.boundvalue(c2) 
             }
             if isNumb(c1) && isNumb(c2) {
                 c = vm.cons(boxInt(op(unbox(c1), unbox(c2))) ,c)
             }
             if vm.needGc {
                  vm.pushStack(c)
                  vm.pushStack(n1)
                  vm.pushStack(n2)
                  vm.gc()
                  n2 = vm.popStack()
                  n1 = vm.popStack()
                  c = vm.popStack()
             }
          }
          c,_ = vm.reverse(c)
          vm.ket = vm.cons(c,vm.ket)
      } else if isCell(n1) {
          vm.stripClosure(&n1)
          c := nill
          for vm.pop(&n1,&c1) {
             if isSymb(c1) {
                    c1 = vm.boundvalue(c1) 
             }
             if isNumb(c1) && isNumb(n2) {
                 c = vm.cons(boxInt(op(unbox(c1), unbox(n2))) ,c)
             }
              if vm.needGc {
                  vm.pushStack(c)
                  vm.pushStack(n1)
                  vm.gc()
                  n1 = vm.popStack()
                  c = vm.popStack()
              }
          }
          c,_ = vm.reverse(c)
          vm.ket = vm.cons(c,vm.ket)
      } else if isCell(n2) {
          vm.stripClosure(&n2)
          c := nill
          for vm.pop(&n2,&c2) {
             if isSymb(c2) {
                    c2 = vm.boundvalue(c2) 
             }
             if isNumb(n1) && isNumb(c2) {
                 c = vm.cons(boxInt(op(unbox(n1), unbox(c2))) ,c)
             }
              if vm.needGc {
                  vm.pushStack(c)
                  vm.pushStack(n2)
                  vm.gc()
                  n2 = vm.popStack()
                  c = vm.popStack()
              }
          }
          c,_ = vm.reverse(c)
          vm.ket = vm.cons(c,vm.ket)
      }
  }
}
func (vm *Vm) fRnd() {
    var p value
    if vm.pop(&vm.ket, &p) {
        if isInt(p) {  // random number from 1 to p
            p1 := unbox(p)
            if p1 > 0 {
              p = boxInt(rand.Intn(p1)+1)
            } else {
              p = boxInt(0)
            }
        } else if isCell(p) {
            vm.stripClosure(&p)
            n := vm.length(p)
            n1 := rand.Intn(n)
            for i:=0; i<n1; i++ {
                p = vm.cdr(p)
            }
            p = vm.car(p)
        }
        vm.ket = vm.cons(p, vm.ket)
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
   //if vm.pop2(&vm.ket, &p1, &p2) && vm.pop(&vm.ket,&p) {
   if vm.pop(&vm.ket,&p)&& vm.pop2(&vm.ket, &p1, &p2) {
       if istrue(p) {
          vm.ket = vm.cons(p1, vm.ket)
      } else {
          vm.ket = vm.cons(p2, vm.ket)
      }
   }
}

func (vm *Vm) fDip() {
   var q1,q2 value
   if vm.pop2(&vm.ket,&q1, &q2) {
        vm.bra = vm.cons(q2, vm.bra) 
        vm.bra = vm.cons(eval, vm.bra)
        vm.ket = vm.cons(q1, vm.ket)
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
        //vm.stripClosure(&key)
        if isCell(key) {
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

func (vm *Vm) fTyp() { // type of an element
    var p value
    var t int
    if vm.pop(&vm.ket,&p) {
        switch {
        case isInt(p):
            t=1
        case isPrim(p):
            t=2
        case isSymb(p):
            t=3
        case isCons(p):
            t=4
        case isClosure(p):
            t=5
        default:
            t=0
        }
        vm.ket = vm.cons(boxInt(t),vm.ket)
    }
}

func (vm *Vm) fPrint() {
    var p value
    if vm.pop(&vm.ket,&p){
        vm.printElem(p)
        fmt.Print(" ")
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

func (vm *Vm) fLambda() {
   var quote,keys value
   if vm.pop2(&vm.ket, &keys, &quote) {
       if isAtom(quote) {
          quote = vm.boundvalue(quote)
       }
       if isAtom(quote) { // we need a quotation to do lambda
           return
       } 
       if isDef(keys) {                 // if arguments are not nill ..
           quote = vm.cons(def, quote)  // .. push a definition on quote
           quote = vm.cons(keys, quote)
           if isAtom(keys) {
               quote = vm.cons(esc, quote)
           }
       }
       if isCons(quote) { // make a closure (only of not yet)
          env := vm.newEnv(vm.env)
          vm.ket = vm.cons(vm.closure(quote,env), vm.ket)
       }
   }
}

/*
func (vm *Vm) pushN(list value, n int) (value, int) {
// take  at most n values from ket and push to list
// return list and number of pushed values
// keep list save in case of gc 
    l := nill           
    n1 = 0
    for i:=0; i<n; i++ {
        if vm.pop(&vm.ket, &val) {
            n1 += 1   
            l = vm.cons(val,l)
            if vm.needGc {
                vm.pushStack(list)
                vm.gc()
                list = vm.popstack()
            }
        } else {
            break       // ket to small 
        }
    } 
    return l, n1
}
*/

func (vm *Vm) deepBind(keys, val value) {
// recursively bind all values of list keys to atom val
// keys must be a list, val an atom
    //vm.printElem(keys); fmt.Println()
    //vm.printElem(val); fmt.Println()
    var key value
    for vm.pop(&keys,&key) {
        if isAtom(key) {
            vm.bindKey(key,val)
            if vm.needGc {
                vm.pushStack(keys)
                vm.gc()
                keys = vm.popStack()
            }
        } else {  // key itself is a list
            vm.pushStack(keys)
            vm.deepBind(key, val)
            keys = vm.popStack()
        }
    }
} 

func (vm *Vm) match(keys, vals value) {
// bind elements from keys to elements from vals with pattern matching
// keys must be a list
    //fmt.Println("match")
    //vm.printElem(keys); fmt.Println()
    //vm.printElem(vals); fmt.Println()
   var key, val value
   if isAtom(vals) {
       vm.deepBind(keys, vals)
       return
   }                             // Q: do we need an else here??
   for vm.pop(&keys, &key) {
       if isNil(keys) {
           vm.bindKey(key,vals)
       } else {
           vm.pop(&vals, &val)   // Q: do we need to check that vals is a list??
           vm.bindKey(key,val)
       }
       if isAtom(vals) {
         vm.deepBind(keys, vals)
         return
       }
   }
}

func (vm *Vm) fDef() {
   var key, k, val value
   var n1 int
   if vm.pop(&vm.ket, &key) {
       if isAtom(key) { 
           if vm.pop(&vm.ket, &val) {
              vm.bindKey(key,val)  // bind key to val in top env-frame
           }
       } else if isDef(key) {      // binding a list of keys
           n := vm.lengthNonQuoted(key)
           n1 = 0    // push max n values from ket to stack
           for i:=0; i<n && vm.pop(&vm.ket, &val); i++ {
                   n1 += 1   
                   vm.pushStack(val)
           } 
           for i:=0; i<n1; i++ {  // make the bindings
               vm.pop(&key, &k) 
               // Q: do need to check that key is still a list in the line below??
               if k == vesc && vm.pop(&key, &k){  // this is interpreted as set
                  vm.setKey(k,vm.popStack())
               } else if isAtom(k) {
                  vm.bindKey(k,vm.popStack())
               } else {
                  elem := vm.popStack()
                  vm.pushStack(key)  // safe key in case of gc
                  vm.match(k,elem)
                  key = vm.popStack()
               }
           }
       }
   }
}


/* we remove set (replaced by backtick)
func (vm *Vm) fSet() {
   var key, val value
   if vm.pop2(&vm.ket, &key, &val) {
       if isAtom(key) {
         vm.setKey(key,val)  // bind key to val in top env-frame
       }
       //vm.printList(vm.env); fmt.Println()
       //vm.printList(vm.car(vm.env)); fmt.Println()
   }
}
*/

func (vm *Vm) fEval() {
    var op value
    if vm.pop(&vm.ket,&op){
        switch {
        case isCons(op):
            vm.evalCons(op)
        case isClosure(op):
            vm.evalClosure(op)
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

func (vm *Vm) evalCons(op value) {
    if isCell(vm.bra) {
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

func (vm *Vm) evalClosure(clos value) {
    op := vm.car(clos)
    env := vm.newEnv(vm.cdr(clos))
    if isCell(vm.bra) {  // no tail position
       vm.depth++
       vm.pushStack(vm.env)
       vm.pushStack(vm.bra)
       vm.pushStack(op)
    } else {   // tail position
       vm.replaceStack(op)
    }
    vm.env = env
    vm.bra = op
}


func (vm *Vm) evalNumb(n value) {
    vm.ket = vm.cons(n,vm.ket)
}

func (vm *Vm) evalSymb(sym value) {
    val := vm.boundvalue(sym)
    if isCons(val) {
        vm.evalCons(val)
    } else if isClosure(val) {
        vm.evalClosure(val)
    } else {
        vm.ket = vm.cons(val,vm.ket)
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
    case rot:
        vm.fRot()
    case cons:
        vm.fCons()
    case car:
        vm.fCar()
    case cdr:
        vm.fCdr()
    case eval:
        vm.fEval()
    case dip:
        vm.fDip()
    case rec:
        vm.fRec()
    case def:
        vm.fDef()
    //case set:
    //    vm.fSet()
    case lambda:
        vm.fLambda()
    //case whl:
    //    vm.fWhl()
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
    case rnd:
        vm.fRnd()
    case eq:
        vm.fEq()
    case iff:
        vm.fIf()
    //case cond:
    //    vm.fCond()
    case esc:
        vm.fEsc()
    case vesc:
        vm.fVesc()
    case val:
        vm.fVal()
    case trace:
        vm.fTrace()
    case typ:
        vm.fTyp()
    case print:
        vm.fPrint()
    default:
        fmt.Println("Error: unknown primitive")
        vm.printElem(p); fmt.Println()
        
    }
}


func (vm *Vm) evalBra() {
      //fmt.Println("Start eval ")
      // Q: should we check for isAtom(vm.bra) ??
    startingDepth := vm.depth
    vm.pushStack(vm.bra)
    var e value
    for {
        if vm.trace > 0 { 
            //fmt.Println("trace")
            vm.printBra(vm.bra)
            vm.printKet(vm.ket)
            vm.printBra(vm.env)
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
        if isAtom(vm.bra) {    // exit scope
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
    
    // load prelude
    vm.bra = vm.loadFile("prelude.clj")
    vm.evalBra()


    //prog := "whl [gt 0 dup add 1] 1 -50000000"  // 5e7, 3 sec on Mac
 
    
    //prog := "eval [ rec gt 0 dup add 1 ] -5 trace 1"
    //prog := "eval [rec gt 0 dup add 1 dup] -5 trace 0"
    //prog := "1 \\' 10  "
    //prog := "eval [ rec gt 0 dup add 1 ] -50000000"   // 5e7, 3.9 sec on MAc
    //prog := "eval [ rec gt 0 dup add 1 ] -500000000"   // 5e8, 24.9 sec on MAc
    
    /*prog := "ack 3 10 def ack' \\[m n]"+
    "[cond "+
    "  [ [ack - m 1 ack m - n 1]"+
    "    [ack - m 1 1]  [eq 0 n]"+
    "    [+ n 1]  [eq 0 m]] ]"
*/

   /* prog := "ack 3 8 def ack' "+
    "[cond "+
    "  [ [ack - m 1 ack m - n 1]"+
    "    [ack - m 1 1]  [eq 0 n]"+
    "    [+ n 1]  [eq 0 m]] def [m n]]"

*/

    //prog := "drop2 1 2 3"
    //prog := "if 0 20"

    //prog := "fac 4 def fac' [eval if rot [1 drop] [* fac - swap 1 dup] eq 1 dup]"
    //prog := "+ [2 1 7] [4 8]"
    //prog := "rot 1"
    prog := "rot 1"

    if len(os.Args) == 2 {
        vm.bra = vm.loadFile(os.Args[1])
    } else {
       vm.bra = vm.makeBra(prog)
    }

    vm.printBra(vm.bra)
    
    vm.evalBra()
    //fmt.Println(vm.bra)


    vm.printKet(vm.ket)

}


/* todos

 - defining the symbol ket, creates a new local ket in 
    the current environment (could be useful in combination 
    with closures)

 - dot operator to access symbol bindings in closure
 
 - reverse if (to emulate cond)
   and direct evaluate (so we avoid ubiquitous eval if)
  
- "if"+
   "gt"  is parsed wrong
*/

// go build  -gcflags=-B