# Bracket
Simple concatenative language geared towards genetic programming


## Bracket is
- a simple, stack-base, functional, [concatenative language](https://en.wikipedia.org/wiki/Concatenative_programming_language)
- inspired by concatenative languages, such as [Joy](http://www.kevinalbrecht.com/code/joy-mirror/joy.html), [Factor](http://factorcode.org), Consize, [XV](http://www.nsl.com/k/xy/xy.htm), or postscript
- implemented in go.
- geared towards genetic programming. This means that code is terse, only a few types are realised. As a side effect this may make the language a candidate for code golfing

## Details


#### Quotations, bras and kets
Bracket is a stack-based language, using normal Polish notation, similar to Forth
Data, functions, and code in Bracket are hold in a stack (also called a quotation).

A quotation can hold Integers, Symbols (any text items) and other quotations. For example `[1 2 dup [2 +] ]`

All functions take the stack as argument and produce a new stack as result.

A Bracket program consists essentially of two stacks
 - the bra which holds the past of the computation
 - the ket, which hold the future of the computation

Making abuse of [Dirac's bra-ket notation] (https://en.wikipedia.org/wiki/Bra%E2%80%93ket_notation),
the program will be written as  `<bra | ket>`. The top of the stack is always at the `|`, the end of the stack is always at the "arrow" end. This makes it easy to see the order of a quote. For example, the stack `[1 2]` is written in bra-form as `<1 2|` and in  the (transposed) ket-form as `| 2  1 >`.

Thus, concatenation of quotations corresponds to the composition of functions.


### Evaluation
During an evaluation, successively elements are taken from the bra. An element that is not a symbol is just pushed on top of the ket. A symbol will take elements from the top of the ket, evaluated them and return a modified bra and ket. The program terminates, when the bra is empty.

### Built in primitives
`swap` `dup` `drop`

### Combinators
Bracket includes so-called combinators (similar to higher-order functions in other functional languages). Combinators are functions that act on, and usually unquote, quotations. 

### Examples
- `<dup | 1 2>  `  evaluates to `| 1 1 2 >`
- `<add | 1 2>  `  evaluates to `| 3 >`
- `<swap | 1 2>  `  evaluates to `| 2 1 >`
- `<drop | 2 1>  `  evaluates to `| 1 >`


### Genetic programming
Concatenative_programming_languages lend themselves for genetic programming. Some of the reasons are: minimal and exceptional simple syntax, which facilitates manipulation of programs, genetic operations and mutations. Every quotation is a valid program. The terseness of the language, one of the drawbacks of such languages for human programmers, becomes a big plus in a genetic environment.

