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
  //test("cons 4 x'", "[x . 4]")  // dotted list
      

  // def, eval  
  test("def x' 2 10", "10")    // x bound to 2, def consumes value on ket
  test("x",           "[]")    // unbound variable evaluates to nill
  test("x def x' 2",  "2")     
  test("x y x def y' 3 def x' 2",  "2 3 2")     

  test("val x' def x' 2",       "2")
  test("x` def x' 2",           "2")  // escape value of symbol to ket
  test("val [1 2]",             "[1 2]")
  test("[1 2]`",                "[1 2]")
  test("x` def x' [1 2 3]",     "[1 2 3]")  
  test("val x' def x' [1 2 3]", "[1 2 3]")  
  test("x def x' [1 2 3]",      "1 2 3")  
  test("val x' def x' 3 x` def x' 2", "3 2")  

  test("f 2 3 def f' [add]",     "5")
  test("foo def foo' [add 1] 2", "3")
  test("eval [x def x' 2]",  "2")     


  test("x set x' 2",  "2")     
  test("x y x set y' 3 set x' 2",  "2 3 2")     
  test("f 2 3 set f' [add]",     "5")
  test("eval [x set x' 2]",  "2")     

  test("eval [x def x' 2] def x' 3",  "2")         // local scope
  test("eval [x] def x' 2]",  "2")    // inner scope can use value defined outside     
  test("eval [x] set x' 2]",  "2")    // inner scope can use value defined outside     
  test("x eval [x def x' 2] x def x' 3",  "3 2 3")  // def changes only within scope
  test("x eval [x set x' 2] x set x' 3",  "2 2 3")  // set changes also outside 

  test("eval [add 1] 2",    "3")
  test("eval [add] 1 2",    "3")
  test("eval [1 2 3] 4",    "1 2 3 4")
  test("eval [] 1",         "1")    // eval empty list
  test("eval add' 1 2",     "3")    // eval symbol
  test("eval foo' def foo' 5", "5")  
  test("eval foo' def foo' bar'", "bar")  
  //test("eval foo' 1 2 def foo' [add]", "3")  
  //test("eval foo' 1 2 def foo' add'", "3")  
  
  // lambda and lexical scoping
  test("eval eval [       [x def x' + 1 x`] def x' 10] def x' 1", "2")
  test("eval eval [lambda [x def x' + 1 x`] def x' 10] def x' 1", "11")
  test("eval swap eval dup eval [lambda [x def x' + 1 x`] def x' 10] def x' 1", "12 11")
  test("foo foo def foo' eval [lambda [x def x' + 1 x`] def x' 10] def x' 1","12 11")

  test("100 eval [f def x' 20] def f' lambda [x] def x' 10","100 10")
  test("100 eval [f def x' 20] def f'        [x] def x' 10","100 20")
  test("100 g def g' [f def x' 20] def f' lambda [x] def x' 10","100 10")
  test("  g g def g' [f def x' 20] def f' lambda [x] def x' 10","20 10")
  // the 2nd evaluation of g is the last command and due to tail elimination
  //  f is evaluated in the global scope

  // closure
  test("a 2 a 3 def a' [make_adder 4] a 2 a 3 def a' [make_adder 5]"+
       "def make_adder' [addx def x']"+
       "def addx' [+ x z def z']" , "6 7 7 8")


  // math and logic 
  test("add 2 3", "5")
  test("+ 2 3", "5")
  test("sub 3 2", "1")
  test("- 3 2", "1")
  test("mul 3 2", "6")
  test("* 3 2", "6")
  test("div 8 2", "4")
  test("/ 8 2", "4")
  test("/ 8 0", "0")
  test("- * 2 4 10", "-2")
  test("lt 4 10", "1")
  test("lt 10 4", "0")
  test("gt 10 4", "1")
  test("> 10 4", "1")
  test("> 10 10", "0")
  test("> 4 10", "0")

  test("add 2 x' def x' 3", "5")
  test("sub 3 x' def x' 2", "1")
  test("sub y' 2 def y' 3", "1")
  test("sub y' x' def y' 3 def x' 2", "1")

  // eq  
  test( "eq 2 2", "1")
  test("eq 2 3", "0")
  test("eq [] []", "1")
  test("eq [1 2 3] [1 2 3]", "1")
  test("eq [1 2 3] [1 2 4]", "0")
  test("eq 2 []", "0")
  test("eq 2 x'", "0")
  test("eq x' x'", "1")
  test("eq y' x'", "0")
  test("eq y' x' def y' x'", "0")

  // if  
  test( "if 20 30 1", "20")
  test("if 20 30 0", "30")
  test("if 20 30 [7]", "20")
  test( "if 20 30 []", "30")
  test( "if 20 30 foo'", "20")
  test( "if 20 30", "")
  test("if 20", "")
  test("if foo' bar' 1", "foo")

  // recur
  test("eval [ rec gt 0 dup add 1 ] -5", "0")   //simple loop

  test("__show results__", "")
}

