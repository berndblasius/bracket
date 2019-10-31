# Bracket
Simple concatenative language geared towards genetic programming


## Bracket is
- a simple, stack-based, functional, [concatenative language](https://en.wikipedia.org/wiki/Concatenative_programming_language)
- geared towards genetic programming, meaning i) that code is terse, having a minimal syntax, and ii) diversity of language features had been traded-off for the ability to efficiently store and compute a huge number of programs in series
- inspired by concatenative languages, such as [Joy](http://www.kevinalbrecht.com/code/joy-mirror/joy.html), [Factor](http://factorcode.org), [Consize](https://github.com/denkspuren/consize),
[Postscript](https://en.wikipedia.org/wiki/PostScript)
- includes (dynamically) lexically scoped variables and a normal (prefix) Polish notation
- here implemented in go (other implementation exist in Julia and C).

- still a work in progress...

## Motivation
Bracket is designed as a generic programming language for genetic programming.
It can be understood as a marriage between a functional language (aka Lisp) and a concatenative language (aka Forth).
The underlying assumption being that a functional language allows for terse,
very compact programs. A concatenative language ensures a minimal syntax, which
facilitates mutations of programs and helps to ensure that every syntacically correct code yields a running program. Thus, Bracket combines the expressive power of a
Lisp with the minimai syntax of Forth. 
Bracket includes lexically scoped enviroments and closures, but function arguments are passed on the global stack (the ket).

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
with normal Poish notation and lexical scoping.

## Details


#### Quotations, bras and kets
Bracket is a stack-based language. Similar to Forth,
data, functions, and code are hold in a stack (also called _quotation_ as in concatenative languages, and _list_ as in lisp).
A quotation can hold numbers, symbols (any text items), and other quotations. For example `[1 2 dup [2 +]]`.

In Bracket two stacks play a special role:
 - the _bra_, which holds the program code
 - the _ket_, a global stack which holds function arguments and the results of the computation

Making abuse of [Dirac's bra-ket notation from quantum mechanics](https://en.wikipedia.org/wiki/Bra%E2%80%93ket_notation),
bras and kets are written with angular brackets and vertical bars, 
e.g., `<foo|` for a bra, and `|bar>` for a ket.
All other stacks are denoted with square
brackets (e.g., `[baz]`).

The bra and all other quotations beside the ket use normal (reversed) Polish notation, where
the top of the stack is placed at the right position and the bottom of the stack on the left position.
Only the ket uses a reverse (infix) Polish notation, where the top of the stack is on the left and the bottom at the right.
For example, the stack `[1 2 [3 4] 5]` with topmost element 5 is written in bra-form as `<1 2 [3 4] 5|` and in  the (transposed) ket-form as `|5 [3 4] 2  1>`. 

A Bracket program essentially consists of the bra and the ket (and an environment) and 
is written as  `<bra | ket>`. 
In loose analogy to quantim-mechanis, bra and ket can be regarded as vectors,
where 
the ket corresponds to a vector of arguments, while the bra corresponds to a vector of operators. The result of a computation then is achieved by _applying_ the bra on the ket (aka scalar product in quantum mechanics) ` <bra | ket>`.



### Evaluation
During an evaluation, successively elements are taken from the bra. An element that is not a symbol is just pushed on top of the ket. A symbol is interpreted as an operator that takes arguments from the top of the ket, evaluates them and returns a modified bra and ket. The program terminates, when the bra is empty. The ket then holds the result of the computation.
Thereby, concatenation of quotations corresponds to the composition of functions.

### Examples (further examples in test directory)

- Numbers and quotation are sequentialy shifted from the bra to the ket. Since the ket is printed in reversed order 
the order of elements visually does not change
  - `<1 2 | >  `  first evaluates to `<1| 2>`
 and then to `<| 1 2>`
  - `<1 2 3 ; this is a comment | >  `  evaluates to `|1 2 3>`
  - `<1 [2 3]  | >  `  evaluates to `|1 [2 3]>`

- esc:
  Symbols can be escaped from bra to the ket without evaluation with the escape operator `esc` or the short notation `'` 
  - `<foo esc| >  `  evaluates to `|foo>`
  - `<foo'| >  `  evaluates to `|foo>`


- Stack shuffling operators
  - `<dup | 1 2>`   evaluates to `|1 1 2>`
  - `<swap | 1 2>`  evaluates to `|2 1>`
  - `<drop | 2 1>`  evaluates to `|1>`

- Math operators
  - `<+ | 3 2>`  evaluates to `|5>`
  - `<- | 3 2>`  evaluates to `|1>`
  - `<* | 3 2>`  evaluates to `|6>`

- Logical values (0 and empty list [] code for logical false, everything else is logical true)  
  - `<gt | 3 2>`  evaluates to `|1>`  ; greater than
  - `<gt | 2 3>`  evaluates to `|0>`  ; 
  - `<eq | 3 2>`  evaluates to `|0>`  ; equality
  - `<eq | 3 3>`  evaluates to `|1>`  
  - `<eq | foo foo>` evaluates to `|1>`  
  - `<eq | [1 x] [1 x]>` evaluates to `|1>`  

- `def` defines a new variable in current scope

  - `<def x'| 3>`  evaluates to `|>` ; x is bound to 3  
  - `<x x def x' 3 |>`  evaluates to `|3 3>`  
  - `<+ x def x' 2 3| >`  evaluates to `|5>`  
  - `<x [x def x' 5 ] x def x' 3 |>`  evaluates to `|3 5 3>`  


- `eval` stores the current brack, takes a stack from the ket, which is evaluated in a new environment, and finally restores the old bra
  - `<eval [+ 1] 3 | >`  evaluates to `|4>` 


### Built in primitives
- stack shuffling operator: `swap` `dup` `drop`
- math operators: `+`, `-`, `>` 
- list operations: `car`, `cdr`, `cons`
- logical and flow control: `if`, `recur`, `cond`, `whl`
- evaulation: `eval`
- variable definition: `def`

### Combinators
Bracket includes so-called combinators (similar to higher-order functions in other functional languages). Combinators are functions that act on, and usually unquote, quotations. 


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


