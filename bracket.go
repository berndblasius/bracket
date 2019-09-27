package main

/* go version of bracket with type tags

*/


import "fmt"

const CELLS = 100*1024*1024
const GC_MARGIN = CELLS - 24


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


const TAG_TYPE   = 15 // mask with bits 1111
const TAG_LOCAL  = 0  // bits 0000    cell on local heap
const TAG_GLOBAL = 8  // bits 1000    cell on global heap
const TAG_PRIM   = 1  // bits 0001
const TAG_SYMB   = 9  // bits 1001
const TAG_INT    = 3  // bits 0011
const TAG_FLOAT  = 11 // bits 1011

type value int

func box_cell(x int) value {return value(x<<4) }   // since TAG_LOCAL=0
func box_global(x int) value {return value(x<<4 | TAG_GLOBAL)}
func box_prim(x int) value {return value(x<<4 | TAG_PRIM)}
func box_symb(x int) value {return value(x<<4 | TAG_SYMB)}
func box_int(x int)  value {return value(x<<4 | TAG_INT)}
//func box_float(x Int) Int { Int(reinterpret(Int32,x)) << 32 | TAG_FLOAT}

func unbox(x value) int  {return int(x)>>4}   // remove all tags
func ptr(x value) int    {return int(x)>>4}   // in contrast to C, here the pointer is
                   // just the heap index, that is, a number
//func unbox_float(x value) float = reinterpret(Float32, Int32(x>>32))

func isInt(x value)    bool {return (x & TAG_TYPE) == TAG_INT}
func isFloat(x value)  bool {return (x & TAG_TYPE) == TAG_FLOAT}
func isPrim(x value)   bool {return (x & TAG_TYPE) == TAG_PRIM}
func isSymb(x value)   bool {return (x & TAG_TYPE) == TAG_SYMB}
func isLocal(x value)  bool {return (x & TAG_TYPE) == TAG_LOCAL}
func isGlobal(x value) bool {return (x & TAG_TYPE) == TAG_GLOBAL}

func isCons(x value) bool {return (x & 1) == 0}
func isNumb(x value) bool {return (x & 3) == 3}
func isAbstractSymb(x value) bool {return (x & 3) == 1}

func isNil(x value) bool {return x == NIL}
func isDef(x value) bool {return x != NIL}
 

type Cell struct {
    car  value
    cdr  value
}

const (
        NIL value = iota<<4 | TAG_PRIM
        DUP
        DROP
        SWAP
        WHL
        ADD
        GT
        UNBOUND
)


type Vm struct {
    bra  value
    ket  value
    root value
    next int   // index to next entry on heap
    heap []Cell
    altheap []Cell
    need_gc bool
}

func (vm *Vm) reset() {
    vm.next = 0
    //vm.heap = [cells]Cell
    //var heap [cells]Cell
}

func init_vm() Vm {
    h := make([]Cell, CELLS)
    ah := make([]Cell, CELLS)
    vm := Vm{NIL,NIL,NIL,0,h,ah,false}
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
   ah := vm.altheap[indv]
   if ah.car == UNBOUND {
       return ah.cdr
   }
   ind := vm.next
   vm.heap[ind]   = vm.altheap[indv]
   vm.altheap[indv] = Cell{UNBOUND, box_cell(ind)}
   vm.next += 1
   return box_cell(ind)
}

func (vm *Vm) gc() {
   fmt.Println("starting gc ************************************************")
   var c Cell
   vm.altheap, vm.heap = vm.heap, vm.altheap
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

   // scan remaining objects in heap (including objects added by this loop)
   for finger < vm.next {
      c = vm.heap[finger]
      vm.heap[finger] = Cell{vm.relocate(c.car), vm.relocate(c.cdr)}
      finger += 1
  }

   //println("GC: live objects found: ", vm.next-1)

   if vm.next >= GC_MARGIN {
       fmt.Println("Bracket GC, heap too small")
   }
   vm.need_gc = false
   fmt.Println("GC finished")
}

// **********************


func (vm *Vm) cons(pcar value, pcdr value) value {
   vm.next += 1
   if vm.next > GC_MARGIN {
     vm.need_gc = true
   }
   vm.heap[vm.next] = Cell{pcar,pcdr}
   return box_cell(vm.next)  // return a boxed index
}

func (vm *Vm) cons_it(pcar value, pcdr *value) {
   *pcdr = vm.cons(pcar, *pcdr)
}

/*
func (vm *Vm) pop_a(list value) (value, value) { //unsafe, we assume that list is a cons
    c := vm.heap[list.(Ptr)]
    return c.car, c.cdr
}
func (vm *Vm) pop2_a(list value) (value, value, value) { //unsafe, we assume that list is a cons
   var p,p1,l value
   p,l = vm.pop_a(list)
   p1,l = vm.pop_a(l)
   return p, p1, l
}
*/

/* Pop first item in list 
   we need care to differentiate between empty list
   and list which has empty list as top element
*/
func (vm *Vm) pop(list *value, p *value) bool {
    //ptr, isCons := (*list).(Ptr)
    if isNil(*list) {
       *p = NIL
       return false
    }
    if isCons(*list) {
        ind := unbox(*list)
        c := vm.heap[ind]
        *p = c.car
        *list = c.cdr
    } else {  // pop from atom returns the atom
        *p = *list
        *list = NIL
    }
    return true
}

// Pop first two items in list 
func (vm *Vm) pop2(list, p1, p2 *value) bool {
    return vm.pop(list, p1) && vm.pop(list, p2);  
}

func (vm *Vm) car (p value) value {
    return vm.heap[p>>4].car
}

//func (vm *Vm) car_ptr (p Ptr) value {
//    return vm.heap[p].car
//}

func (vm *Vm) cdr (p value) value {
    return vm.heap[p>>4].cdr
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



//isCons2(x,vm) = isCons(x) && isCons(cdr(x,vm))
//isCons3(x,vm) = isCons(x) && isCons(cdr(x,vm)) && isCons(cddr(x,vm))


//func pop2(list,vm) {
//   p,l = pop(list,vm)
//   if isCons(l)
//       p1,l = pop(l,vm)
//       return p, p1, l
//   end
//   return p, l , NIL
//}

func (vm *Vm) f_dup() {
   //if isCons(vm.ket)
   //ket_ptr, ok := vm.ket.(Ptr)
   //if ok {
       vm.ket = vm.cons(vm.car(vm.ket), vm.ket)
      // vm.ket = vm.cons(vm.car_ptr(ket_ptr), vm.ket)
   //end
    //}
       // fmt.Println("dup")

}


func (vm *Vm) f_plus() {
    var n1,n2 value
    if vm.pop2(&vm.ket, &n1, &n2) {
      n3 := box_int(unbox(n1) + unbox(n2))
      vm.ket = vm.cons(n3,vm.ket)
      //fmt.Println("add: ",n3)
  }
}

func gt(x,y int) int {
   if x>y {
       return 1
   } else  {
       return 0
   }
}

func (vm *Vm) f_gt() {
    var n1,n2 value
    if vm.pop2(&vm.ket, &n1, &n2) {
      //n1,n2, vm.ket = vm.pop2(vm.ket)
      n3 := box_int(gt(unbox(n1), unbox(n2)))
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
          return x != NIL
    default:
          return true
    }
}


func (vm *Vm) f_whl() {
   fmt.Println("Start whl ")
   var q,b value
   if vm.pop2(&vm.ket,&q, &b) {
      //vm->root = cons(vm->bra,vm->root,vm);
      //vm->root = cons(q,vm->root,vm);
       vm.root = vm.cons(vm.bra,vm.root)
       vm.root = vm.cons(q,vm.root)
       for istrue(b) {
         vm.bra = vm.car(vm.root)
         vm.eval_bra()
         if !vm.pop(&vm.ket,&b){
             break
         }
       }
       vm.root = vm.cdr(vm.root)
       vm.bra = vm.car(vm.root)
       vm.root = vm.cdr(vm.root)
   }
}

func (vm *Vm) eval_bra() {
      //fmt.Println("Start eval ")
    var e value
    for {
        vm.pop(&vm.bra,&e);
        //e, vm.bra = vm.pop(vm.bra)
        //fmt.Println("e=",e)
       
        switch e { 
        case DUP:
            vm.f_dup()
        case WHL:
            vm.f_whl()
        case ADD:
            vm.f_plus()
        case GT:
            vm.f_gt()
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
	fmt.Printf("rock'n roll\n")     // c
    //var e value
    //var bra, e value 
    //bra = vm.cons(DUP,NIL)
    //e, bra = vm.pop(bra)
    vm := init_vm()
    //vm.bra = vm.cons(ADD,vm.bra)
    //vm.bra = vm.cons(LT,vm.bra)
    //vm.bra = vm.cons(DUP,vm.bra)
    //vm.bra = vm.cons(I(30),vm.bra)
    //vm.bra = vm.cons(I(40),vm.bra)
    //vm.ket = vm.cons(I(1),vm.ket)

    // char str[] = "whl [ gt 0 dup add 1 ] 1 -2";
    // While loop for benchmark
    len := -10
    len = -5e7   // 3 sec on MAc
    len = -5e8   // 18 sec on MAc
    //len = -50000000
    q := vm.cons(GT,NIL)
    q = vm.cons(box_int(0),q)
    q = vm.cons(DUP,q)
    q = vm.cons(ADD,q)
    q = vm.cons(box_int(1),q)
    vm.bra = vm.cons(WHL,vm.bra)
    vm.bra = vm.cons(q,vm.bra)
    vm.bra = vm.cons(box_int(1),vm.bra)
    vm.bra = vm.cons(box_int(len),vm.bra)
    
    
    vm.eval_bra()
    fmt.Println(vm.bra)

    //vm.gc()

    l := Cell{box_int(1),box_cell(2)}
    fmt.Println(l)
}


// 2-15, 14-20, 59-60
