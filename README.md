# Bracket
Simple concatenative functional language geared towards genetic programming


## Bracket is
- a simple, stack-based, functional, [concatenative language](https://en.wikipedia.org/wiki/Concatenative_programming_language)
- geared towards genetic programming, meaning i) that code is terse, having a minimal syntax, and ii) diversity of language features had been traded-off for the ability to efficiently store and compute a huge number of programs in series
- inspired by concatenative languages, such as [Joy](http://www.kevinalbrecht.com/code/joy-mirror/joy.html), [Factor](http://factorcode.org), [Consize](https://github.com/denkspuren/consize),
[Postscript](https://en.wikipedia.org/wiki/PostScript)
- includes (dynamically) lexically scoped variables and a normal (prefix) Polish notation
- here implemented in go (other implementation exist in Julia and C)

- still a work in progress...

## Motivation
Bracket is designed as a generic programming language for genetic programming.
It can be understood as a marriage between a functional language (aka Lisp) and a concatenative language (aka Forth).
The underlying assumption being that a functional language allows for terse,
very compact programs. A concatenative language ensures a minimal syntax, which
facilitates mutations of programs and helps to ensure that every syntacically correct code yields a running program. Thus, Bracket combines the expressive power of a
Lisp with the minimai syntax of Forth. 
Bracket includes tail recursion, lexically scoped enviroments and closures, but function arguments are passed on the global stack (the ket).

To make the similarity to lisp-languages
even more transparent, Bracket uses a normal Polish notation (rather than reversed
Polish notation, as used by most other concatenative languages). While unsual at first,
this leads to a remarkable formal similarity between Bracket and lisp code.

Bracket does not aim to compete for a productive programming language. On purpose the language is lean. It has
neither the performance, nor does it support the many types, libraries, and ecosystem of a mature programming language (e.g. clojure, common lisp, or factor, postscript). In contrast, Bracket excels in its
design space: genetic programming. It allows expressive prgramming with a
minimal syntax and is optimized for the iterative performant calculation of a huge number of small
programs. 
In addition, Bracket can be regarded as an experiment how far we can go by combining a concatenative language
with prefix Poish notation and lexical scoping.

## Details


#### Quotations, bras and kets
Bracket is a stack-based language. Similar to Forth,
data, functions, and code are hold in a stack (also called _quotation_ as in concatenative languages, and _list_ as in lisp).
A quotation can hold numbers, symbols, and other quotations. For example `[1 2 dup [2 +]]`.

In Bracket two stacks play a special role:
 - the _bra_, which holds the program code
 - the _ket_, a global stack which holds function arguments and the results of the computation

Making abuse of [Dirac's bra-ket notation from quantum mechanics](https://en.wikipedia.org/wiki/Bra%E2%80%93ket_notation),
bras and kets are written with angular brackets and vertical bars, 
that is, a left angular bracket `<foo|` for a bra , and a right angular bracket `|bar>` for a ket.
All other stacks are denoted with square
brackets (e.g., `[baz]`).

Apart from the ket, all other quotations (including the bra) use prefix (normal) Polish notation, where
the top of the stack is placed at the right position and the bottom of the stack on the left.
Only the ket is written in prefix (polish) Polish notation (similar to Forth), where the top of the stack is on the left and the bottom at the right.
For example, the stack `[1 2 [3 4] 5]` with topmost element 5 is written in bra-form as `<1 2 [3 4] 5|` and in 
ket-form is transposed to postfix notation `|5 [3 4] 2  1>` (note that the small quotation `[3 4]` inside the ket uses prefix again).

A Bracket program essentially consists of a bra and a ket (as well as an environment) and 
is written as  `<bra|ket>`. 
In loose analogy to Dirac's notation in quantum-mechanics, bra and ket can be regarded as vectors,
where 
the ket corresponds to a vector of arguments, while the bra takes the role of a vector that acts as an operator. The result of the computation then is achieved by _applying_ the bra on the ket (aka scalar product in quantum mechanics) ` <bra|ket>`.



### Evaluation
During an evaluation, successively single elements are taken from the top of the bra. If the element is not a symbol it is pushed on top of the ket. If the element is a symbol it is interpreted as an operator that takes arguments from the top of the ket, evaluates them, and returns a modified bra and ket. The program terminates, when the bra is empty. The ket then holds the result of the computation.
Thus, computation in Bracket follows the usual logistic of a
[concatenative language](https://en.wikipedia.org/wiki/Concatenative_programming_language).
In particular, concatenation of quotations corresponds to the composition of functions.

Bracket supports nested environments.
An environments is a linked list of frames, and every frame is a linked list of bindings.  
When evaluating a symbol on the top of the bra (or with the builtin `eval` entering a new scope) the current bra is stored on the stack, a new environment is created, and the associated quotation is evaluated in this environment, and finally the old bra is restored from the stack.

### Examples (further examples can be found in the test directory)

- Numbers and quotation are sequentialy shifted from the bra to the ket. Since the ket is printed in reverse notation 
the order of elements visually does not change
  - `<1 2 | >  `  first evaluates to `<1| 2>`
 and then to `<| 1 2>`
  - `<1 2 3 ; this is a comment |>  `  evaluates to `|1 2 3>`
  - `<1 [2 3]| `  evaluates to `|1 [2 3]>`

- esc:
  Symbols can be escaped from evaluation
  with the escape operator `esc` (or the short notation `'`) 
  and thereby asrejust shifted from the bra to the ket without evaluation 
  - `<foo esc|`  evaluates to `|foo>`
  - `<foo'|`  evaluates to `|foo>`

- Stack shuffling operators
  - `<dup 1 2|`   evaluates to `|1 1 2>`
  - `<swap 1 2|`  evaluates to `|2 1>`
  - `<drop 1 2|`  evaluates to `|2>`

- Stack operators
  - `<car  [1 2 3]|`   evaluates to `|3 [1 2]>`
  - `<cdr [1 2 3]|`   evaluates to `|[1 2]>`
  - `<cons 3 [1 2]|`   evaluates to `|[1 2 3]>`
  - `<cons car [1 2 3]|`   evaluates to `|[1 2 3]>`

- Math operators
  - `<+ 3 2|`  evaluates to `|5>`
  - `<- 3 2|`  evaluates to `|1>`
  - `<* 3 2|`  evaluates to `|6>`
  - `</ 7 2|`  evaluates to `|3>` (integer division)

- Logical values (0 and empty list [] code for logical false, everything else is logical true)  
  - `<gt 3 2|`  evaluates to `|1>`  ; greater than
  - `<gt 2 3|`  evaluates to `|0>`  ; 
  - `<eq 3 2|`  evaluates to `|0>`  ; equality
  - `<eq 3 3|`  evaluates to `|1>`  
  - `<eq foo foo|` evaluates to `|1>`  
  - `<eq [1 x] [1 x]|` evaluates to `|1>`  

- `if` takes three elements from ket, if third element is true, the first element is pushed on the ket, else the second is pushed on the ket
  - `<if foo' bar' 1|`  evaluates to `|foo>` 
  - `<if foo' bar' 0|`  evaluates to `|bar>` 

- `rec` anonymous recursion (supports tail recursion)
  - `<eval [rec gt 0 dup add 1] -5|` evaluates to `|0>` 
    at the begin of any evaluation the bra is saved, 
    rec takes an argument from the ket, if the argument is true
    the bra is replaced by the saved bra, i.e., a simple recursion
    (this example codes a simple loop from -5 to 0) 

- `def` defines a new variable in current scope

  - `<def x' 3|`  evaluates to `|>`, but the symbol is bound to the number 3 in the current environment (note that x needs to be escaped) 
  - `<x def x' 3|`  evaluates to `|3>` (evaluating a symbol that is bound to a number pushes the number on the ket) 
  - `<+ x 3 def x' 2| `  evaluates to `|5>`  (first, the number 2 is pushed on the ket, then x is bound to 2, then 3 and the binding of x are pushed on the stack, finally the + operator adds the two elements on the ket)
  - `<foo def foo' [add 1] 2|` evaluates to `|3>`

- `val` pushes the bound value of a symbol on ket without evaluation (short form is backtick `)
  - `<val l' def l' [1 2 3]|` evaluates to `|[1 2 3]>`
  - ``<l` def l' [1 2 3]|`` evaluates to `|[1 2 3]>`  ; short form


- `eval` and environments  
  - `<eval [+ 1] 3|`  evaluates to `|4>` ; the quotation is evaluated in a new environment
  - `<x eval [x def x' 5 x] def x' 3|`  evaluates to `|3 5 3>` ; new bindings are possible in inner scopes (these are dynamical scopes as definitions are made during runtime); old scopes are preserved after leaving the inner scope

- `lambda`, closures and lexical scoping
  - `<lambda x' [+ x 1]|`  evaluates to `|[+ x 1 def x']>`  
  Thereby the following has happened:
    - a lexical closure is created that contains the quotation and the current enviroment. By printing the closure on the ket only the associated quotation is shown, so from inspection of the ket one cannot distinguish a closure from a pure quotation.
    - to handle the argumets of the lambda, the expression `def x'` is pushed on the quotation. Thus, formally we can regard this as a lambda expression, a function that takes the argument x.
  - `<\x' [+ x 1]|`  evaluates to `|[+ x 1 def x']>`  (short notation for lambda)
  - `<eval \x' [+ x 1] 4|`  evaluates to `|5>`
  - `<f 4 def f' \x' [+ x 1]|`  evaluates to `|5>`
  - `<f 4 def f' \[x] [+ x 1]|`  evaluates to `|5>` (when the argument is in a list, we don't need to escape)
  - `<f 4 3 def f' \[x y] [- x y]|`  evaluates to `|1>` (lambda with multiple arguments)



### Built in primitives
- stack shuffling operator: `swap`, `dup`, `drop`
- math operators: `+`, `-`, `>` 
- list operations: `car`, `cdr`, `cons`
- logical and flow control: `if`, `recur`, `cond`, `whl`
- evaulation: `eval`
- variable definition: `def`
- lambda: `lambda`
- escape and quotation: `esc`, `val`

### Combinators
Bracket includes so-called combinators (similar to higher-order functions in other functional languages). Combinators are functions that act on, and usually unquote, quotations. 


### Examples
- Factorial
   - tail recursive
    ```clojure
    fac 4 ; this code evaluates to 24
    def fac' \[n]
       [drop swap eval \[acc cnt] 
         [rec < cnt n * acc cnt + cnt 1] 1 1]
    ```

   - alternative code, using stack shuffling
   ```clojure
   fac 4 ; evaluates to 24
   def fac' 
      [eval if [1 drop] 
               [* fac - swap 1 dup] 
       eq 1 dup]
   ```

- Ackermann function
  ```clojure
  ack 3 4 ; evaluates to 125
  def ack' \[m n]
    [cond [ 
        [ack - m 1 ack m - n 1]
        [ack - m 1 1]  [eq 0 n]
        [+ n 1]  [eq 0 m]]] 
  ```

### Genetic programming
Concatenative_programming_languages lend themselves for genetic programming. Some of the reasons are: minimal and exceptional simple syntax, which facilitates manipulation of programs, genetic operations and mutations. Every quotation is a valid program. The terseness of the language, one of the drawbacks of such languages for human programmers, becomes a big plus in a genetic environment.

As a side effect this may make the language a candidate for code golfing

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


