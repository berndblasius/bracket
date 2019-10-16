package main

/* go version of bracket with type tags

*/


import (
    "fmt"
)


const cells = 100*1024*1024
const gc_margin = cells - 24

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
func ptr(x value) int    {return int(x)>>4}   // in contrast to C, here the pointer is
                   // just the heap index, that is, a number
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
        nill value = iota<<4 | tagPrim
        dup
        drop
        swap
        whl
        add
        gt
        unbound
)

var primStr = map[value] string {
    dup: "dup", drop: "drop", swap: "swap", whl: "whl",
    add: "add", gt: "gt",
}

var str2prim = map[string] value {
    "dup": dup, "drop": drop, "swap": swap, "whl": whl,
    "add": add, "gt": gt,
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
    need_gc bool   // flag to indicate that heap space gets rare
    stats stats   // some statistics about the running program
    trace int   //trace mode e: 0=no trace, 1=trace non-verbose, 3=verbose

}

func (vm *Vm) reset() {
    vm.next = 0
    vm.stats = stats{0,0,0,0}
    vm.bra = nill
    vm.ket = nill
    vm.aux = nill
    vm.env = nill
    vm.root = nill
    vm.trace = 0
    //vm.arena = [cells]cell
    //var arena [cells]cell
}

func init_vm() Vm {
    h := make([]cell, cells)
    ah := make([]cell, cells)
    stats := stats{0,0,0,0}
    vm := Vm{nill,nill,nill,nill,nill,0,h,ah,false,stats,0}
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
   //vm.env = vm.relocate(vm.env)
   //for i = 1: vm.stackindex
   //  vm.stack[i] = relocate!(vm.stack[i],vm)
   //end

   // scan remaining objects in arena (including objects added by this loop)
   for finger < vm.next {
      c = vm.arena[finger]
      vm.arena[finger] = cell{vm.relocate(c.car), vm.relocate(c.cdr)}
      finger += 1
  }

   //println("GC: live objects found: ", vm.next-1)

   if vm.next >= gc_margin {
       fmt.Println("Bracket GC, arena too small")
   }
   vm.need_gc = false
   fmt.Println("GC finished")
}

// **********************


func (vm *Vm) cons(pcar, pcdr value) value {
   vm.next += 1
   if vm.next > gc_margin {
     vm.need_gc = true
   }
   vm.arena[vm.next] = cell{pcar,pcdr}
   return boxCell(vm.next)  // return a boxed index
}

func (vm *Vm) consVal(pcar value, pcdr *value) {
   *pcdr = vm.cons(pcar, *pcdr)
}

/* Pop first item in list 
   we need care to differentiate between empty list
   and list which has empty list as top element
*/
func (vm *Vm) pop(list, p *value) bool {
    //ptr, isCons := (*list).(Ptr)
    if isNil(*list) {
       *p = nill
       return false
    }
    if isCons(*list) {
        c := vm.arena[unbox(*list)]
        *p = c.car
        *list = c.cdr
    //} else {  // pop from atom returns the atom
    //    *p = *list
    //    *list = nill
    } else {  // pop from atom returns the atom and false
        //*p = *list
        *p = nill
        return false
    }
    return true
}

// Pop first two items in list 
func (vm *Vm) pop2(list, p1, p2 *value) bool {
    return vm.pop(list, p1) && vm.pop(list, p2);  
}

func (vm *Vm) car (p value) value {
    return vm.arena[p>>4].car
}

func (vm *Vm) cdr (p value) value {
    return vm.arena[p>>4].cdr
}

// unsafe
func (vm *Vm) caar(p value) value { 
    return vm.car(vm.car(p))
}
//func cadr(p Ptr, vm Vm) value { 
//    return car(cdr(p,vm).(Ptr),vm)
//}
//func cddr(p Ptr, vm Vm) value { 
//    return cdr(cdr(p,vm).(Ptr),vm)
//}


func (vm *Vm) reverse(list value) value {
    var p value 
    l := nill
    if vm.pop(&list,&p) {
       if isAtom(list) && isDef(list) { // dotted pair
         return vm.cons(list,p)
       } else {
          vm.consVal(p,&l)
       }
    }
    for vm.pop(&list,&p) {
       vm.consVal(p,&l)
    }
    return l
}

// [ 1 ptr] [2 nil]  --> [2 ptr] [1 nil]
// [1 ptr] [2 3]

func (vm *Vm) reverse1(list value) value {
    var p value 
    l := nill
    for vm.pop(&list,&p) {
       vm.consVal(p,&l)
    }
    return l
}

func (vm *Vm) length(list value) int {
// just count the number of conses, ie dotted list has length 1
   n := 0
   for isDef(list) {
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

func (vm *Vm) fPlus() {
    var n1,n2 value
    if vm.pop2(&vm.ket, &n1, &n2) {
      n3 := boxInt(unbox(n1) + unbox(n2))
      vm.ket = vm.cons(n3,vm.ket)
      //fmt.Println("add: ",n3)
  }
}

func myGt(x,y int) int {
   if x>y {
       return 1
   } else  {
       return 0
   }
}

func (vm *Vm) fGt() {
    var n1,n2 value
    if vm.pop2(&vm.ket, &n1, &n2) {
      //n1,n2, vm.ket = vm.pop2(vm.ket)
      n3 := boxInt(myGt(unbox(n1), unbox(n2)))
      vm.ket = vm.cons(n3,vm.ket)
      //fmt.Println("lt: ",n3)
    }
}

/*
  
func (vm *Vm) f_lt() {
    n1, vm.bra = split(vm.bra)
    n2, vm.bra = split(vm.bra)
    vm.bra = con(lt(n2,n1),vm.bra)
}

  
func (vm *Vm) f_whl() {
    q, vm.bra = split(vm.bra)
    b, vm.bra = split(vm.bra)
    ket_safe = vm.ket
    while is_true(b)
      vm.ket = q
      run!(vm)
      if isempty(vm.bra); break; end
      b, vm.bra = split(vm.bra)
    end
    vm.ket = ket_safe
  }
void f_whl(Vm *vm) 
{
  //    printf("WHL  \n");
  //   printElem(bra); printf("\n");
  //    printElem(ket); printf("\n");
     // printElem(root); printf("\n \n");
   any q,b;
   if (pop2(&vm->ket,&q, &b)) {
      vm->root = cons(vm->bra,vm->root,vm);
      vm->root = cons(q,vm->root,vm);
     // printf("WHL  \n");
     // printElem(bra); printf("\n");
     // printElem(ket); printf("\n");
     // printElem(root); printf("\n \n");
       while (unbox(b)) {
         //bra = q;
        //printf("WHL 2  \n");
         vm->bra = car(vm->root);
       // printElem(bra); printf("\n");
       // printElem(root); printf("\n \n");
         eval_bra(vm);
       // printf("WHL 3  \n");
       // printElem(bra); printf("\n");
       // printElem(ket); printf("\n");
         if (!pop(&vm->ket,&b)) break;
       }
       vm->root = cdr(vm->root); 
       pop(&vm->root,&vm->bra);
   }
}
*/

//istrue(l) = isDef(l) ?  (unbox(l) != 0) : false
//istrue(l) = isDef(l) && unbox(l) != 0
//isfalse(l) = isNil(l) || unbox(l) == 0
//istrue(x) ((isNumb(x) && (unbox(x) != 0)) || (isSymb(x) && isDef(x)) || (isCons(x)))

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

func (vm *Vm) evalBra() {
      //fmt.Println("Start eval ")
    var e value
    for {
        vm.pop(&vm.bra,&e);
        //e, vm.bra = vm.pop(vm.bra)
        //fmt.Println("e=",e)
       
        switch e { 
        case dup:
            vm.fDup()
        case drop:
            vm.fDrop()
        case swap:
            vm.fSwap()
        case whl:
            vm.fWhl()
        case add:
            vm.fPlus()
        case gt:
            vm.fGt()
        default:
            vm.ket = vm.cons(e,vm.ket)
            //fmt.Println("default")
            
        }

        if vm.need_gc {
            vm.gc()
        }
        if isNil(vm.bra) {break}
    }
}

func main() {
    fmt.Printf("rock'n roll\n")   
    vm := init_vm()
    
    //prog := "add 1 2"
    //prog := "whl [gt 0 dup add 1] 1 -50000000"  // 5e7, 3 sec on Mac
    prog := "1 2 3 ; this is a comment 4"    
    vm.bra = vm.makeBra(prog)
    
    //vm.bra = vm.loadFile("test.clj")
   
    vm.printBra(vm.bra)
    //vm.printList(vm.bra) 
    //vm.printList(vm.reverse(vm.bra))
    
    vm.evalBra()
    fmt.Println(vm.bra)

    //vm.gc()

    vm.printList(vm.bra)
    vm.printKet(vm.ket)
    //vm.printElem(vm.bra, true)

    
    fmt.Println("test dotted")
    //ll := vm.cons(1,2)
    //ll := vm.cons(dup,drop)
    //ll := vm.cons(boxInt(10),boxInt(20))
    ll := vm.cons(boxInt(10),nill)
    ll = vm.cons(dup,ll)
    ll = vm.makeBra("dup drop [1 2 3]")
    ll = vm.cons(ll, swap)
    ll = vm.cons(swap, boxInt(10))

    vm.printList(ll); fmt.Println()

    lll := vm.reverse(ll)

    
    
    vm.printList(lll)
    vm.printList(vm.reverse(ll))
}

