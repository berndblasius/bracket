# Bracket
Simple concatenative functional language geared towards genetic programming


## Bracket is
- a simple, stack-based, functional, [concatenative language](https://en.wikipedia.org/wiki/Concatenative_programming_language)
- geared towards genetic programming, meaning that i) code is terse, having a minimal syntax, and ii) richness of language features is traded-off for the ability to efficiently store and run a large number of small programs in series
- inspired by concatenative languages, such as [Joy](http://www.kevinalbrecht.com/code/joy-mirror/joy.html), [Factor](http://factorcode.org), [Consize](https://github.com/denkspuren/consize),
[Postscript](https://en.wikipedia.org/wiki/PostScript),
but uses a normal (prefix) Polish notation

- a marriage of a concatenative language with Lisp, that is, it includes dynamically and lexically scoped variables, lambdas and closures

- here implemented in Go (other implementation exist in Julia and C)

- still a work in progress...


## Motivation
Bracket is designed as a generic programming language for genetic programming.
It can be understood as a marriage between a functional language (aka Lisp) and a concatenative language (aka Forth).
The underlying assumption being that a functional language allows for terse,
very compact programs. A concatenative language ensures a minimal syntax, which
facilitates mutations of programs and helps to ensure that every syntacically correct code yields a running program. Thus, Bracket combines the expressiveness of
Lisp with the minimai syntax of Forth. 
Bracket includes tail recursion, lexically scoped enviroments and closures, but function arguments are passed on a global stack (the ket).

To make the similarity to lisp-languages
more transparent, Bracket uses a normal (prefix) Polish notation, rather than reversed (postfix)
Polish notation as used by most other concatenative languages. While unsual at first,
this results in a remarkable formal similarity between Bracket and Lisp code.

Bracket does not aim to compete for a productive programming language. On purpose the language is lean. It has
neither the performance, nor does it support the many types, libraries, and the ecosystem of a mature programming language (e.g. Clojure, Common Lisp, or Factor, Postscript). In contrast, Bracket excels in its
design space: genetic programming. It allows expressive prgramming with a
minimal syntax and is optimized for the performant iteration of a huge number of runs of rather small
programs. 
In addition, Bracket can be regarded as an experiment how far we can go by combining a concatenative language
with prefix Poish notation and lexical scoping.

## Details

#### Quotations, bras and kets
Bracket is a stack-based language. Similar to Forth,
data, functions, and code are hold in stacks (also called _quotations_ as in concatenative languages, and _lists_ as in Lisp).
A quotation can hold any other literals (numbers, symbols) and other quotations. For example, the quotation `[1 2 dup [2 +]]` holds the numbers 1 and 2, the symbol `dup` and the quotation `[2 +]`.

Bracket currently supports as types only integers, symbols and quotations. But nothing precludes implementation of further types.

In Bracket two stacks play a special role:
 - the _bra_, which holds the current program code, and
 - the _ket_, a global stack which holds function arguments and the results of the computation

Making abuse of [Dirac's bra-ket notation from quantum mechanics](https://en.wikipedia.org/wiki/Bra%E2%80%93ket_notation),
bras and kets are written with angular brackets and vertical bars, 
that is, a left angular bracket `<foo|` for a bra , and a right angular bracket `|bar>` for a ket.
All other stacks are denoted with square
brackets (e.g., `[baz]`).

Apart from the ket, all other quotations (including the bra) use prefix (normal) Polish notation, where
the top of the stack is placed at the right position and the bottom of the stack 
is placed on the left.
Only the ket is written in postfix (reversed) Polish notation (similar to Forth), where the top of the stack is on the left and the bottom at the right.
For example, the stack `[1 2 [3 4] 5]` with topmost element 5 is written in bra-form as `<1 2 [3 4] 5|` and in 
ket-form is transposed to postfix notation `|5 [3 4] 2  1>`. Note that the small quotation `[3 4]` inside the ket uses prefix notation again.

A Bracket program essentially consists of a bra and a ket (as well as an environment) and 
is written as  `<bra|ket>`. 
In loose analogy to Dirac's notation in quantum-mechanics, bras and kets can be regarded as vectors,
where 
the ket corresponds to a vector of arguments, while the bra takes the role of a vector that acts as an operator. The result of the computation then is achieved by _applying_ the bra on the ket (aka scalar product in quantum mechanics) ` <bra|ket>`.



### Evaluation
During an evaluation, successively single elements are taken from the top of the bra. If the element is not a symbol or builtin-operator it is pushed on top of the ket. Otherwise it is interpreted as a function that takes arguments from the top of the ket (possibly also from the bra), evaluates them, and pushes the result on the ket (thereby possibly also modifying the bra).
The program terminates, when the bra is empty. The ket then holds the result of the computation.
Thus, computation in Bracket follows the usual program flow of a
[concatenative language](https://en.wikipedia.org/wiki/Concatenative_programming_language).
In particular, concatenation of quotations corresponds to the composition of functions.


### Documentation by examples
Further examples can be found in the test directory

#### Concatenative part

- Numbers and quotation are sequentialy shifted from the bra to the ket. Since the ket is printed in reverse notation 
the order of elements visually does not change
  - `<1 2| >  `  is first evaluates to `<1|2>`
 and then to `<|1 2>`
  - `<1 [2 3]| `  evaluates to `|1 [2 3]>`

  - single line comments separated by `;`
  - `<1 2 3 ; this is a comment |>  `  evaluates to `|1 2 3>`


- Stack shuffling operators
  - `<dup 1 2|`   evaluates to `|1 1 2>`  ; duplicate top element on ket
  - `<drop 1 2|`  evaluates to `|2>`  ; remove top element on ket
  - `<swap 1 2|`  evaluates to `|2 1>`
  - `<rot 1 2 3|`  evaluates to `|3 2 1>`


- Stack operators
  - `<car  [1 2 3]|`   evaluates to `|3>`  ; top of a list
  - `<cdr [1 2 3]|`   evaluates to `|[1 2]>` ; remainder of a list
  - `<cons 3 [1 2]|`   evaluates to `|[1 2 3]>`
  - `<cons car swap cdr dup [1 2 3]|`   evaluates to `|[1 2 3]>`

- Escaping:
  Symbols can be escaped from evaluation by framing into a quotation
  - `<car [dup] 1 2|`  evaluates to `|dup 1 2>`

  - Alternatively Bracket provides the 
  escape operator `esc` (or the short notation `'`), 
  thereby any element on top of the bra is shifted to the ket without evaluation 
  - `<dup esc 1 2|`  evaluates to `|dup 1 2>`
  - `<dup' 1 2|`  evaluates to `|dup 1 2>`

- Math operators
  - `<+ 3 2|`  evaluates to `|5>`
  - `<- 3 2|`  evaluates to `|1>`
  - `<* 3 2|`  evaluates to `|6>`
  - `</ 7 2|`  evaluates to `|3>` (integer division)
  - `<rnd 7|`  pushes a random number between 1 and 7 onto ket
  - `<rnd [1 2 3]|`  pushes a random element from the list onto ket

- Logical values (0 and empty list [] code for logical false, everything else is logical true)  
  - `<gt 3 2|`  evaluates to `|1>`  ; greater than
  - `<gt 2 3|`  evaluates to `|0>`  ; 
  - `<lt 2 3|`  evaluates to `|1>`  ; less than 
  - `<eq 3 2|`  evaluates to `|0>`  ; equality
  - `<eq 3 3|`  evaluates to `|1>`  
  - `<eq foo' foo'|` evaluates to `|1>`  
  - `<eq foo' foo'|` evaluates to `|1>`  
  - `<eq [1 x] [1 x]|` evaluates to `|1>`  

- `if` takes three elements from the ket, if the first element is true, the second element is pushed on the ket, else the third is pushed on the ket
  - `<if 1 foo' bar'|`  evaluates to `|foo>` 
  - `<if 0 foo' bar'|`  evaluates to `|bar>` 


- `eval` evaluates a quotation on top of the ket (which then acts as a function). 
During the evaluation, the current bra is first stored on the stack and then replaced by the quotation. Next, this new bra is evaluated. Finally, the 
old bra is restored from the stack. Note that all evaluations work on the same ket.  
In other concatenative languages this operator is called `i`. Bracket uses the notation `eval` because `i` is too valuable as a variable name.
  - `<eval [+ 2 3]|` evaluates to `|5>`
  - `<20 eval [+ 2 3] 10|` evaluates to `|20 5 10>`

- Conditional evaluation with `eval if`
  - `<eval if 1 [+ 1][- 1] 5|` evaluates to `|6>`
  - `<eval if 0 [+ 1][+ 2] 5|` evaluates to `|7>`  
  - combination of `eval if` statements is similar to `if-elseif` clauses in other languages  
  `<eval if 0 [+ 1][eval if 0 [+ 2] [+ 3]] 5|` evaluates to `|8>`

- `rec` anonymous recursion   
    At the begin of any evaluation of a quotation, a copy of the quotation is saved for possible recursion. 
    `rec` takes an argument from the ket, if the argument is true
    the bra is replaced by the saved quotation, i.e., a simple recursion
  - `<eval [rec gt 0 dup + 1 dup] -5|` evaluates to `|0 -1 -2 -3 -4 -5>` 
    (a simple loop from -5 to 0) 

  - `rec` does not need to be the last element  
  `<eval [100 rec gt 0 dup + 1 dup] -5|` evaluates to `|100 0 -1 -2 -3 -4 -5>`   

- Protection from evaluation.   
  `dip` evaluates a quotation; however before evaluation, the top element of the ket
  is stored on the bra (and thus protected). The quotation is then evaluated on the remainder of the ket 
  - `<dip [+ 2] 10 1|` first translates to `<10 eval|[+ 2] 1 >`    
  which then evaluates to `|10 3>`

-  Combinators  
  These operators can be combined to powerfull code-snippets and higher-order functions, called *combinators*, see functions defined in the prelude below.

- Turing completeness  
 This concatenative part already suffices to make Bracket Turing-complete. The difference to other concatenative languages (and the full power of Bracket) becomes apparent when combining this with the Lispy-part of Bracket.

#### Lispy part

Bracket supports variable bindings and nested environments.
An environments is a list of frames, and every frame is a list of bindings. 
New bindings can be generated with `def`.  
In contrast to many other concatenative languages, which allow to define new words with the `: ;` notation, Bracket's `def` follows the same semantics as any other concatenative operator.


- `def` binds a variable to a value in the current scope

  - `<def x' 3|`  evaluates to empty ket`|>`, but the symbol `x` is bound to the number 3 in the current environment (note that x needs to be escaped) 
  - `<x def x' 3|`  evaluates to `|3>` (evaluating a symbol that is bound to a number pushes the number on the ket) 
  - `<+ x 3 def x' 2| `  evaluates to `|5>`  (first, the number 2 is pushed on the ket, then x is bound to 2, then 3 and the binding of x are pushed on the stack, finally the + operator adds the two elements on the ket)
  - `<x def [x] 1|` evaluates to `|1>` (the symbol can also be provided in a quotation)
  - `<a b def [a b] 1 2|` evaluates to `|1 2>` (`def` applied to a list iteratively defines all elements from the list)

- Evaluation  
If the top of the bra is a variable, it is evaluated. For this, the current bra is saved and replaced by the binding of the variable. The binding is then evaluated, and finally the old bra is restored.
  - `<foo foo def foo' [add 1] 2|` evaluates to `|4>`
  - `<x|`evaluates to `|[]>` (an unbound variable evaluates to the empty list)

- `val` pushes the bound value of a symbol on ket without evaluation (short form is backtick `)
  - `<val l' def l' [1 2 3]|` evaluates to `|[1 2 3]>`
  - ``<l` def l' [1 2 3]|`` evaluates to `|[1 2 3]>`  ; short form


- `eval` and environments  
Evaluation also creates new environments. Thus, a more complete description of an evaluation is as following: During an evaluation,
the current bra is stored on the stack and replaced by the quotation.
Then, a new empty environment (or scope) is created on top of the current environment.
Then, the new bra is evaluated in this environment. Finally, the environment is destroyed and the
old bra restored from the stack.
Note that during the evaluation the bra was changed, a new environment was created, but the ket remained.
  - Dynamic scoping   
   New bindings are possible in inner scopes (these are dynamical scopes, as definitions are made during runtime). 
   Old scopes are preserved after leaving the inner scope.  
  `def` first searches for the variable in the current scope. If the variable is found, it is replaced by the new binding. If the variable does not exist in the current scope a new binding is created  
  `<x eval [x def x' 5 x] def x' 3|`  evaluates to `|3 5 3>`   

  - Redefining values  
  Values in deeper binding can be overwritten with the backtick operator(which acts similar to `set!` in Scheme)  
  ``<def [x`] 3|`` In this example, `def` searches in all nested scopes for a binding of `x`. If it is found, the binding is redefined. When no binding to the variable exists in the whole environment, a new binding is generated.
  
  -  `<x eval [x def x' 2] x def x' 3|` evaluates to  `3 2 3`   
  (def changes bindings only within current scope)  
  ``<x eval [x def [x`] 2] x def x' 3|`` evaluates to `2 2 3`   
  (set changes bindings also outside)


  - Tail recursion  
  `eval` and `rec` support tail recursion. That is, if the last operation during an evaluation is another evaluation, no new environment is created.

- `lambda`, closures and lexical scoping
  - `<lambda x' [+ x 1]|`  evaluates to `|[+ x 1 def x']>`  
  Thereby the following has happened:
    - a lexical closure, that contains the quotation and the current enviroment, is created and pushed on the ket. By printing the closure on the ket only the associated quotation is shown, so from inspection of the ket one cannot distinguish a closure from a pure quotation.
    - to handle the argumets of the lambda, the expression `def x'` is pushed on the quotation. Thus, when the closure is evaluated, at first the top value on the ket is bound to the symbol x. In this sense, 
    formally, we can regard this as a lambda expression, a function that takes the argument x.
  - short notation for lambda   
   `<\x' [+ x 1]|`  evaluates to `|[+ x 1 def x']>`  
  - usually the arguments are passed as a list  
   `<\[x] [+ x 1]|`  evaluates to `|[+ x 1 def x']>`   
   `<eval \[x] [+ x 1] 4|`  evaluates to `|5>`
  - closures can be bound to variables  
  `<f 4 def f' \[x] [+ x 1]|`  evaluates to `|5>`
  - lambdas can take multiple arguments  
   `<f 4 3 def f' \[x y] [- x y]|`  evaluates to `|1>` 



### Built in primitives
- stack shuffling operator: `swap`, `dup`, `drop`, `rot`
- math operators: `+`, `-`, `>` 
- list operations: `car`, `cdr`, `cons`
- logical and flow control: `if`, `rec`
- evaulation: `eval`, `dip`
- variable definition: `def`
- lambda: `lambda`
- escape and quotation: `esc`, `val`

##### Still missing
- Macros  
Macros are not yet implemented. The main reason being first, that in Bracket function arguments are not evaluated before function application. Thus, many algorithms that must be implemented as a macro in Lisp can be implemented as a function in Bracket. 
Second, at the moment there exists no compiler for Bracket, so there is no speed advantage to evaluate macros before compilation.

- Call-cc and continuations

- More types (strings, arrays, hashs, structs)


### Prelude examples
The file prelude.clj contains a number of predefined functions and combinators
and is loaded at the begin of a computation.  
Here some useful examples:

- `def swapd'  [swap rot]`; (abc -- acb) ; deep swap
- `def rot1' [rot rot]`   ; (abc -- bca)
- `def over' [swap rot dup swap]` ; (ab -- bab)
- `def rot4' [swap dip rot']`  ; (abcd -- dabc)
- `def rot14' [dip rot1' swap]` ; (abcd -- bcda)
- `def splt' [car swap cdr dup]`
- `def dip2' [dip dip' swap]`  
- `def dip3' [dip dip2' swap]`  
- `def keep' [dip eval' over`]`  ; eval function but retain top argument from ket
- `def bi' [eval dip [keep]]`     ; eval two functions on the same ket argument
- `def curry' [cons esc' cons swap]`  
- `def repeat' \[n foo] [eval [rec n def n' sub n 1 foo]]`;
 ; (n foo -- ) ; repeat foo n-times
- `def each' [drop drop eval if rot [rec dup swap dip2 eval' over rot1 splt][drop drop] dup swap]`
- `def reduce' [each swapd]`
- `def prod' [reduce [mul] 1]`
- `def sum'  [reduce [add] 0]`
- `def size' [reduce [add 1 drop] 0]`

### Code Examples
- Factorial
   - tail recursive
    ```clojure
    fac 4 ; this code evaluates to 24
    def fac' \[n]
       [drop swap eval \[acc cnt] 
         [rec lt cnt n * acc cnt + cnt 1] 1 1]
    ```


   - alternative code, using stack shuffling
   ```clojure
   fac 4 ; evaluates to 24
   def fac' 
      [eval if eq 1 rot 
            [1 drop] 
            [* fac - swap 1 dup] 
       dup]
   ```

- Ackermann function
  ```clojure
  ack 3 4 ; evaluates to 125
  def ack' \[m n]
     [eval if eq 0 m 
        [+ n 1] 
     [eval if eq 0 n 
         [ack - m 1 1] 
     [ack - m 1 ack m - n 1]]]
  ```



- simple closure for bank account
  ```clojure
  account withdraw' 20 account deposit' 40 
  def account' make-acc 50  
  def make-acc' [ 
      \[m][eval if eq m withdraw' 
              [withdraw] 
          [eval if eq m deposit' 
              [deposit] 
          [unknown']]] 
      def withdraw' [ 
          eval if gt balance rot  
          [balance def [balance`] - balance] 
          [insuff' drop] dup]  
      def deposit' [balance def [balance`] + balance] 
      def balance' ]
  ```


### Remarks on prefix notation
- Note that due to prefix notation functions must be defined in the source code **below** their usage!
- As stated above, Bracket is an experiment also in coding with prefix notation. First experiences seems to indicate that this is something one can get used to. In a certain sense this yields a natural code order, starting from the most general definitions, while specialize function appear later.

## Genetic programming
Concatenative_programming_languages lend themselves for genetic programming. Some of the reasons are: minimal and exceptional simple syntax, which facilitates manipulation of programs, genetic operations and mutations. Every quotation is a valid program. The terseness of the language, one of the drawbacks of such languages for human programmers, becomes a big plus in a genetic environment.

A simplified version of Bracket, intended for genetic progamming, is in the GeneBracket folder.

As a side effect this may make the Bracket a candidate for code golfing.

## Implementation

##### Immutable variables
Program code and variables in Bracket are stored as immutable linked lists. In
a genetic programming environment that means that genetical identical copies of
a programm can be shared by a single pointer. Linked lists are also a
persistent data structure. Thus, programs that differ only in a few mutations,
on average, can share much common code. Thus, the decision for linked lists as
basic data structure is a trade-off in programming efficiency vs storage
capability. Code of linked lists, is not localized in memory and thus looses in
performance. On the other hand, it allows a very compact storage of a huge
number of programs.


##### Memory, types,  and gc
The current implementation of Bracket stacks are stored as linked lists, using pre-allocated memory arenas (holding cons-cells). So
during run-time, no memory allocatations are necessary.
Bracket impements a garbage collector with Cheney copying algorithms (allowing
a non-recursive traversal of live-objects).

Variables types are stored with 4 tagbits, leaving the following data types: 60 bit integers, symbols
with max 10 characters, 32 bit floats, and linked lists.

##### Interpreter
Bracket is currently implemented as an intetreter. While nothing forbids the implementation as a compiled language, interpretation is more convenient for genetic programming (where the compact storage of code and the fast loading and start-up time are more important than efficiency of the programming itself). Being an interpreted language no macros are implemented (similar to PicoLisp and NewLisp).
