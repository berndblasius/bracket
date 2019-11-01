package main

import (
       "fmt"
       "testing"
   )

// generate a closure to safe test statistics
func (vm *Vm) makeTest() (f func(string, string)){
  var ntests, success int
  f = func(code, res string){
      if code == "__show results__" {
         fmt.Println("Number of tests ", ntests)
         fmt.Println("Number of tests passed ", success)
         fmt.Println("Number of tests failed ", ntests - success)
      } else {
          ntests++
          vm.reset()
          vm.bra = vm.makeBra(code)
          vm.ket = nill
          result,_ := vm.reverse(vm.makeBra(res))
          vm.evalBra()
          if vm.isEqual(vm.ket, result) {
            //fmt.Println("test no error")
            //vm.printKet(vm.makeBra(code))
            //vm.printKet(vm.ket)
            //vm.printKet(result)
              success++
              return 
          } else {
            //t.Error() 
            fmt.Println("test error")
            vm.printKet(vm.makeBra(code))
            vm.printKet(vm.ket)
            vm.printKet(result)
            return 
          }
     }
  }
  return
}

func TestBracket(t *testing.T) {
  vm := init_vm()
  //var c, r string
  test := vm.makeTest() 

  test("1 2 3",     "1 2 3")   // values on bra are shifted on ket
  test("1 2 3; this is a comment",    "1 2 3")
  test( "1 [ ]",    "1 [ ]")
  test("[1 2 3]",   "[1 2 3]")
  test("1 [a ; comment \n [" + "2 3 ]]",  "1 [a [2 3 ]]")

  // stack shuffling
  test( "dup 2", "2 2")
  test("dup 2 3", "2 2 3")
  test("dup [2 3] 10", "[2 3] [2 3] 10")
  test("drop 2 3", "3")
  test("drop 2", "")
  test("swap 2 3", "3 2")

  test("x esc 4",  "x 4")       // escape symbol
  test("x' 4",     "x 4")
  test("[1 2 3]'", "[1 2 3]")   // escape a list or symbol..
  test("3'", "3")               // .. just move to ket

  // car, cdr, cons
  test("car [1 2 3] 4", "3 [1 2] 4")
  test("car [x]", "x []")  // bring a symbol on ket
  test("car []", "")       // car empty list
  test("car 1", "")        // car atom
  test("car x' def x' [1 2]",  "2")           // car a symbol ...
  test("car x' car x' def x' [1 2]",  "2 2")  // .. leaves the symbol intact
  test("car x' def x' 1",  "")           

  test("cdr [1 2 3] 4", "[1 2] 4")
  test("cdr []", "[]")    // cdr empty list
  test("cdr 1", "[]")     // cdr atom
  test("cdr x' def x' [1 2 3]",  "[1 2]")            // cdr a symbol ...
  test("x` cdr x' def x' [1 2 3]",  "[1 2 3] [1 2]") // ... leaves symbol intact
  test("cdr foo' def foo' bar'",  "[]")           

  test("cons 1 []", "[1]")
  test("cons 4 [1 2 3]", "[1 2 3 4]")
  test("cons [4 5] [1 2 3]", "[1 2 3 [4 5]]")
  test("cons car [1 2 3 4]", "[1 2 3 4]")
  test("cons car x' cdr x' def x' [1 2 3 4]", "[1 2 3 4]")
      

  // def, eval  
  test("def x' 2 10", "10")    // x bound to 2, def consumes value on ket
  test("x",           "[]")    // unbound variable evaluates to nill
  test("x def x' 2",  "2")     
  test("x y x def y' 3 def x' 2",  "2 3 2")     

  test("val x' def x' 2",   "2")
  test("x` def x' 2",       "2")  // escape value of symbol to ket
  test("val [1 2]",         "[1 2]")
  test("[1 2]`",            "[1 2]")
  test("x` def x' [1 2 3]", "[1 2 3]")  


  test("eval [x def x' 2]",  "2")     // inner scope
  test("eval [x def x' 2] def x' 3",  "2")     // inner scope
  test("x eval [x def x' 2] x def x' 3",  "3 2 3")     // inner scope

  test("eval [add 1] 2",    "3")
  test("eval [add] 1 2",    "3")
  test("eval [1 2 3] 4",    "1 2 3 4")
  test("eval [] 1",         "1")    // eval empty list
  test("eval add' 1 2",     "3")    // eval symbol
  test("eval foo' def foo' 5", "5")  
  test("eval foo' def foo' bar'", "bar")  
  //test("eval foo' 1 2 def foo' [add]", "3")  
  //test("eval foo' 1 2 def foo' add'", "3")  

  // math  
  test("add 2 3", "5")
  test("+ 2 3", "5")
  test("gt 10 4", "1")
  test("> 10 4", "1")
  test("> 10 10", "0")
  test("> 4 10", "0")

  // eq  
  test( "eq 2 2", "1")
  test("eq 2 3", "0")
  test("eq [] []", "1")
  test("eq [1 2 3] [1 2 3]", "1")
  test("eq 2 []", "0")
  test("eq 2 x'", "0")
  test("eq x' x'", "1")
  test("eq y' x'", "0")

  // if  
  test( "if 20 30 1", "20")
  test("if 20 30 0", "30")
  test("if 20 30 [7]", "20")
  test( "if 20 30 []", "30")
  test( "if 20 30 foo'", "20")
  test( "if 20 30", "")
  test("if 20", "")
  test("if foo' bar' 1", "foo")


  test("__show results__", "")
}

