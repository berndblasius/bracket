package main

import (
       "fmt"
       "testing"
   )


func (vm *Vm) testCode(code, res string) int {
  vm.bra = vm.makeBra(code)
  vm.ket = nill
  result := vm.reverse(vm.makeBra(res))
  //printf("bra ");print_bra(vm->bra,vm); 
  vm.evalBra()
  if vm.isEqual(vm.ket, result) {
    //fmt.Println("test no error")
    //vm.printKet(vm.makeBra(code))
    //vm.printKet(vm.ket)
    //vm.printKet(result)
      return 1
  } else {
    //t.Error() 
    fmt.Println("test error")
    vm.printKet(vm.makeBra(code))
    vm.printKet(vm.ket)
    vm.printKet(result)
    return 0 
  }
}

func TestBracket(t *testing.T) {
  vm := init_vm()
  var c, r string
  ntests := 0
  success := 0

  ntests++; success += vm.testCode(
    "1 2 3",    
    "1 2 3")

  ntests++; success += vm.testCode(
    "1 2 3 ; this is a comment",    
    "1 2 3")

  ntests++; success += vm.testCode(
    "1 [ ]",    
    "1 [ ]")

  ntests++; success += vm.testCode(
    "[1 2 3]",    
    "[1 2 3]")


  c =  "1 [a ; comment \n [" + "2 3 ]]"    
  r =  "1 [a [2 3 ]]"    
  ntests++; success += vm.testCode(c,r)

  ntests++; success += vm.testCode(
    "dup 2",
    "2 2")
      
  ntests++; success += vm.testCode(
    "dup 2 3",
    "2 2 3")
      
  c = "dup [2 3] 10"
  r = "[2 3] [2 3] 10"
  ntests++; success += vm.testCode(c, r)

  ntests++; success += vm.testCode(
    "drop 2 3",
    "3")
      
  ntests++; success += vm.testCode(
    "drop 2",
    "")
      
  ntests++; success += vm.testCode(
    "swap 2 3",
    "3 2")
      

  /* *********************************************** */
  fmt.Println("Number of tests ", ntests)
  fmt.Println("Number of tests passed ", success)
  fmt.Println("Number of tests failed ", ntests - success)
  


}
