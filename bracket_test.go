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
      

  // def, val  
  test("def x' 2 10", "10")    // x bound to 2, def consumes value on ket
  test("x",           "[]")    // unbound variable evaluates to nill
  test("x def x' 2",  "2")     
  test("x y x def y' 3 def x' 2",  "2 3 2")     
  
  test("foo-bar def foo-bar' 2", "2")         // dash in symbol name
  test("foo_bar def foo_bar' 2", "2")         // underscore in symbol name
  test("foo1 def foo1' 2", "2")               // digit in symbol name
  test("Fo-13_g def Fo-13_g' 2", "2")         // digit in symbol name
  test("f123456789  def f123456789' 2", "2")  // symbol name with 10 characters

  test("val x' def x' 2",       "2")
  test("x vesc def x' 2",       "2")  // vesc, escape value of symbol to ket
  test("x` def x' 2",           "2")  // backtick = short for vesc
  test("val [1 2]",             "[1 2]")
  test("[1 2]`",                "[1 2]")
  test("x` def x' [1 2 3]",     "[1 2 3]")  
  test("val x' def x' [1 2 3]", "[1 2 3]")  
  test("x def x' [1 2 3]",      "1 2 3")  
  test("val x' def x' 3 x` def x' 2", "3 2")  

  test("f 2 3 def f' [add]",     "5")
  test("foo def foo' [add 1] 2", "3")
  test("foo 1 2 def foo' [add]", "3")    // def a symbol that is not builtin
  test("foo 1 2 def foo' add'", "+ 1 2") // this not 

  test("a def [a] 2",             " 2")    // def: bind a list of keys
  test("a b def [a b] 1 2",       "1 2")    
  test("a b c def [a b c] 1 2 3", "1 2 3")  
  test("a b c def [a b c] 1 2",   "[] 1 2") 
  test("a b` def [a b] 1 [2 3]",   "1 [2 3]")    
  test("def [] 1 2",              "2")    

  test("a` b` def [[a b]] 2", "2 2")         // pattern matching
  test("a` b` c` d` def [[a b [c d]]] 2", "2 2 2 2")
  test("b` def [[b]] [1 2 3]", "[1 2 3]") 
  test("a` b` def [[a b]] [1 2 3]", "[1 2] 3") 
  test("a` b` c` def [[a b c]] [1 2 3]", "[1] 2 3") 

  test("a def [a`] 2",             " 2")    // set
  test("a b def [a` b] 2 3",       " 2 3")    
  test("b a b def [a` b`] 2 3",    "3 2 3")    
  test("f 2 3 def [f`] [add]",     "5")
  test("eval [x def [x`] 2]",      "2")     

  test("x set x' 2",  "2")                  // set as a primary operator, obsolete now
  test("x y x set y' 3 set x' 2",  "2 3 2")     
  test("f 2 3 set f' [add]",     "5")
  test("eval [x set x' 2]",  "2")     

  test("eval [x def x' 2]",  "2")     
  test("eval [x def x' 2] def x' 3",  "2")         // local scope
  test("eval [x] def x' 2]",  "2")    // inner scope can use value defined outside     
  test("eval [x] set x' 2]",  "2")    // inner scope can use value defined outside     
  test("eval [x] def [x`] 2]", "2")    // inner scope can use value defined outside     
  test("x eval [x def x' 2] x def x' 3",  "3 2 3")  // def changes only within scope
  test("x eval [x set x' 2] x set x' 3",  "2 2 3")  // set changes also outside 
  test("x eval [x def [x`] 2] x def [x`] 3",  "2 2 3")  // set changes also outside 
  test("x eval [x def [x`] 2] x def x' 3",  "2 2 3")  // set changes also outside 

  test("eval [add 1] 2",    "3")
  test("eval [add] 1 2",    "3")
  test("eval [1 2 3] 4",    "1 2 3 4")
  test("eval [] 1",         "1")    // eval empty list
  test("eval add' 1 2",     "3")    // eval symbol
  test("eval foo' def foo' 5", "5")  
  test("eval foo' def foo' bar'", "bar")  


  // lambda and lexical scoping
  test("lambda x' [+ x 1]","[+ x 1 def x']")
  test("f def f' lambda x' [+ x 1] 10","11")
  test("eval lambda x' [+ x 1] 10","11")
  // shorter syntax, but we need golang raw strings and cannot use backqoute for val
  test(`eval \x' [+ x 1] 10`,"11")      
  test(`f def f' \x' [+ x 1] 10`,"11")
  // alternative use interpreted strings and "\\""
  test("eval \\x' [+ x 1] 10","11")      

  test("eval \\[x] [+ x 1] 10","11")      // lambda with list of arguments
  test("eval \\[x y] [- x y] 10 2","8")      

  test("eval foo` def foo' \\[x] [+ x 1] 10","11")      // backquote or val .. 
  test("eval val foo' def foo' \\[x] [+ x 1] 10","11")  // .. puts closure on ket 
  test("eval val val foo' def foo' \\[x] [+ x 1] 10","11")  // val on closure puts quotation on stack 
  
  test("eval eval [       [x def x' + 1 x`] def x' 10] def x' 1", "2")
  test("eval eval [\\ [] [x def x' + 1 x`] def x' 10] def x' 1", "11")
  test("eval swap eval dup eval [\\[] [x def [x`] + 1 x`]] def x' 10",         "12 11")
  test("x foo foo def foo' eval [\\[] [x def [x`] + 1 x`]] def x' 10",         "12 12 11")
  test("x foo foo def foo' eval [\\[] [x set x' + 1 x`]]  def x' 10",          "12 12 11")
  test("x foo foo def foo' eval [\\[] [x def x' + 1 x`]]  def x' 10",          "10 11 11")
  test("x foo foo def foo' eval [\\[] [x def [x`] + 1 x`] def x' 10] def x' 1", "1 12 11")
  test("x foo foo def foo' eval [\\[] [x def x' + 1 x`]   def x' 10] def x' 1", "1 11 11")
  test("x foo foo def foo' [x def x' + 1 x`]   def x' 10",                     "10 11 11")
  test("x foo foo def foo' [x def [x`] + 1 x`] def x' 10",                     "12 12 11")

  test(`100 eval [f def x' 20] def f' \[] [x] def x' 10`,"100 10")
  test(`100 eval [f def x' 20] def f'     [x] def x' 10`,"100 20")
  test(`100 g def g' [f def x' 20] def f' \[] [x] def x' 10`,"100 10")
  test(`  g g def g' [f def x' 20] def f' \[] [x] def x' 10`,"20 10")
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
  test("/ 9 2", "4")
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

  // dip
  test("dip [+ 1] 5 2", "5 3")
  test("dip [+ 1] [+ 10] 2", "[+ 10] 3")
  test("dip [1 2 3] 4", "4 1 2 3")


  // cond
  test("cond [[2]]",       "2")                 
  test("cond [[10] [11] [0]]",   "10")
  test("cond [[10] [11] [1]]",   "11")
  test("cond [[3] [1] [eq 4] [2 drop] [lt 4 dup]] 3", "3")                 
  test("cond [[3] [1] [eq 4] [2 drop] [lt 4 dup]] 5", "2")                 
  test("cond foo' 5 def foo' [[3] [1] [eq 4] [2 drop] [lt 4 dup]]", "2")    

  // recur
  test("eval [ rec gt 0 dup add 1 dup] -5", "0 -1 -2 -3 -4 -5")   //simple loop
  test("foo def foo' [ rec gt 0 dup add 1 dup] -5", "0 -1 -2 -3 -4 -5")   //simple loop


  // small examples, to give a feeling for the language
  
  // simple closure for bank account
  test("acc withdraw' 60 "+
  "acc deposit' 100 acc withdraw' 60 acc withdraw' 60 acc deposit' 40 "+
  "def acc' make-acc 50 "+ 
  "def make-acc' [ "+
    "\\[m][ cond [ "+
          "[unknown']"+
          "deposit [eq m deposit']"+ 
          "withdraw [eq m withdraw'] ]"+
   "]"+
   "def withdraw' ["+
       "eval if"+ 
         "[balance def [balance`] - balance]"+
         "[insuff' drop]"+
         "gt balance dup "+ 
   "]"+ 
   "def deposit' [balance def [balance`] + balance]"+
   "def balance' ]", "70 130 insuff 30 90")

  // Factorial
   //simple recursive
   test("fac 4 def fac' [cond [[* fac - swap 1 dup] [1 drop] [eq 1 dup]]]", "24")
   test("fac 4 def fac' \\[n] [cond [[* fac - n 1 n] 1 [eq 1 n]]]", "24")
   // tail recursive
   test("fac 4 def fac' \\[n]"+
      "[drop swap eval \\[acc cnt] [rec < cnt n * acc cnt + cnt 1] 1 1]", "24")

   test("fac 4 def fac' \\[n]"+
      "[drop swap eval [rec < cnt n * acc cnt + cnt 1 def [acc cnt]] 1 1]", "24")

   test("fac 4 def fac' [eval if [1 drop] [* fac - swap 1 dup] eq 1 dup]", "24")


  // Fibonacci numbers with simple recursion
  // simple recursion
  test("fib 6 "+
   "def fib' \\[n] [cond [ [+ fib - n 1 fib - n 2] n [< n 2]]]", "8") 
  test("fib 6 def fib' \\[n] [eval if [n] [+ fib - n 1 fib - n 2]  < n 2]", "8")
  test("fib 6 def fib' [eval if [] [+ fib - swap 1 swap fib - swap 2 dup] > 2 dup]", "8")

  test("fib 6 def fib' \\[n]["+    // tail recursive
        "loop 1 1 n def loop' \\[a b n][ "+
             "eval if [a] [loop b + a b - n 1]"+
           "eq n 0 ] ]", "13")
   
  test("fib 6 def fib' \\[max] ["+
      "drop drop fib-iter 1 0 1 "+
      "def fib-iter' [rec < n max + n 1 j + i j def [n i j]] ]", "13")

   test("fib 6 def fib' \\[max] ["+
      "drop drop fib-iter 1 0 1 "+
      "def fib-iter' \\[n i j] [rec < n max + n 1 j + i j]] ]", "13")

  // Ackermann function
  test("ack 3 4 def ack' \\[m n]"+
    "[cond "+
    "  [ [ack - m 1 ack m - n 1]"+
    "    [ack - m 1 1]  [eq 0 n]"+
    "    [+ n 1]  [eq 0 m]] ]",  "125")


  test("__show results__", "")
}

