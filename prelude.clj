; Bracket prelude


; still missing:
; reverse, append (concat)

def Y' [
  eval if [drop nip] [swapd X]
  swap dip [eq dupd rot1 keep2 eval']
]
def X' [eval dup]

;def loop' [rec swap keep [eval]]
def loop' [rec swap dip eval' over eval']

;def whl' [eval if swap [] curry [whl] keep eval']

;whl f 
;when curry [whl] keep eval' f 
;when curry [whl] f eval f
;when [whl f']

;def whl' [eval if eval rot [whl Rto eval Ris] [drop Rto] toR swap]

def filter' [         ; reverse still missing
  each cons swap [eval if rot cons' drop' swap keep] rot1 [] ]

def size' [reduce [add 1 drop] 0]
def sum'  [reduce [add] 0]
def prod' [reduce [mul] 1]

; we still need to reverse !!!
def append' [reduce [cons] swap]   ; not yet working

;def spread' [eval if rot    ; concat not yet defined
;   [cons concat rot [dip]]
;   [drop]
;   dup swap ]

def cleave2' [drop drop each [keep2]]  ; eval a list of arguments on two arguments from ket
def cleave' [drop each [keep]]  ; eval a list of arguments on a single argument from ket

def reduce' [each swapd]
def unstack' [each []]

;def map1' [each rot1 [] append swap [cons]   ; for this we would need "append" first

def map' [     ; needs a reverse still
  drop drop eval if rot [
     rec dup swap dip2 [cons eval] over rot1 splt
  ]
  [drop drop drop] dup rot rot [] swap
]

def each' [
  drop drop eval if rot [
     rec dup swap dip2 eval' over rot1 splt
  ]
  [drop drop] dup swap
]

;def each' \[foo]
;  [drop eval [rec dup Rto foo toR swap car]]


;def each'
;   [drop Rto drop when
;     [rec dup Rto eval toR swap Ris spltr]
;   dup toR]


def repeat' \[n foo]     ; (n foo -- ) ; repeat foo n-times
  [eval [rec n def n' sub n 1 foo]]


;def rep' \[n foo]     ; (n foo -- ) ; repeat foo n-times
;  [rep1 n
;  def rep1 \[n1]
;    [eval if [rep1][drop] dup sub n1 1 foo]]

;def rep2' [    ; (n foo -- ) ; repeat foo n-times
;    rep2a
;    def rep2a'
;      [eval if [rep2a] [drop Rto drop]
;        dup sub swap 1 Rto eval toR swap Ris]
;    toR swap]

;def rep3' [    ; (n foo -- ) ; repeat foo n-times
;    rep3a
;    def rep3a'
;      [eval if [rep3a] [drop Rto drop]
;        dup sub swap 1 dip Ris]
;    toR swap]


;def do' [eval if [do Rto] [drop Rto] eval toR dup]


;def rep2' [    ; (n foo -- ) ; repeat foo n-times
;  drop Rto drop when
;    [rec dup sub swap 1 Rto eval toR swap Ris]
;  dup toR swap]


; cond via stack manipulations
;def cond2'
;  [eval if
;    [eval if
;         [eval if swap cons swap [cond2] spltr Rto eval toR]
;         [eval drop]
;       dup spltr]
;    [drop]
;  dup val]


;def [cond1 [ xs x b ]'] 
;  [eval if [b] [eval if  x` [cond1 xs`] b] isNil x`]

def cond' \[[xs b]]
  [eval if isNil xs` [b] [cond if b splt xs`]]


def curry' [cons esc' cons swap]  

def bi2star' [eval dip [dip2]]    ; eval two quotation on two arguments each from ket
def bistar' [eval dip [dip]]      ; eval two quotation on two arguments from ket
def tri2' [eval dip [keep2 dip [keep2]]] ; eval 3 functions on the same 2 ket arguments
def tri' [eval dip [keep dip [keep]]]    ; eval 3 functions on the same ket argument
def bi2' [eval dip [keep2]]   ; eval two functions on the same two ket arguments
def bi' [eval dip [keep]]     ; eval two functions on the same ket argument

;def keep' [Rto eval swap toR dup swap]
def keep2' [dip2 dip dup2'] ; eval function but retain two top arguments from ket
;def keep' [dip dip dup']
def keep' [dip eval' over]  ; eval function but retain top argument from ket

def do2' [dip2 dup]
def do' [dip dup]

; def dip2' [Rto Rto eval toR swap toR swap]
def dip3' [dip dip2' swap]  
def dip2' [dip dip' swap]  

; dip with aux
;def dip' [Rto eval toR swap]
;def dip' [x` eval def x' swap]

; dip as a macro
;def dip' [cons eval' cons rot]

def caar' [car car]
def cadr' [car cdr]
def cdar' [cdr car]
def cddr' [cdr cdr]
def caaar' [car car car]
def caadr' [car car cdr]
def cadar' [car cdr car]
def caddr' [car cdr cdr]
def cdaar' [cdr car car]
def cdadr' [cdr car cdr]
def cddar' [cdr cdr car]
def cdddr' [cdr cdr cdr]


; logical operators

def when1' [eval if rot1 [drop] rot1 dup]

;def when' [eval if Rto swap [] toR]
def when' [eval if rot1 []]

def unless' [eval if swap []]
def or' [if dup] 
;def and' [if Rto swap toR dup ]  ; [swapd dup]
def and' [if swap rot dup ]  ; [swapd dup]

;and a b  aab  baa
;aba 
;def nand' [if 1 swap]
;def nand' [if Rto Rto 1 toR toR ]
;def not' [if Rto 0 1 toR] 
def not' [if rot 0 1] 
; xor 0 1 ->1 ;  10-> 1;  00->0 11->0
;xor a b -  0 a b

;def consr
def spltr' [cdr swap car dup]   ; split right
def splt' [car swap cdr dup]

def rot4' [swap dip rot']  ; (abcd - dabc)
def rot14' [dip rot1' swap] ; (abcd - bcda)

;swap dip rot' 1 2 3 4
;swap 1 4 2 3


;def [dup2] [Rto swap Ris dup toR]  ; (ab - abab)
;def [dup2] [Rto swap toR dup Rto dup toR]  ; (ab - abab)
;def dup2' \[x] [x` swap x` dup]  ; (ab - abab)
def dup2' [over over]


;def [dupd] [a` b` b` def b' def a' ]
;def dupd' [Rto dup toR] ; (ab - abb)  ; deep dup
;def dupd' \[x] [x` dup] ; (ab - abb)  ; deep dup
;def dupd' [dip [dup]]
;def dupd' [dip dup']
def dupd' [rot dup swap]


;def swapd' [Rto swap toR] ; abc -- acb   ; deep swap
;def swapd' \[x] [x` swap] ; abc -- acb   ; deep swap
def swapd'  [swap rot]

def nip2' [drop swap drop swap]   ; (abc - a)
def nip'  [drop swap]  ; (ab - a)
;def [over [a b]'] [b` a` b`]   ; (ab - bab)
;def over' \[x] [swap x` dup]
;def over' [swap Rto dup toR]
def over' [swap rot dup swap]
;def over' [swap dip dup']

;def rot1' \[a b c]' [b` c` a`]  ; (abc - bca)
;def [rot1] [Rto swap toR swap]  ; (abc - bca)
def [rot1] [rot rot]
;def [rot1] [swap rot swap]
;def [rot] [swap Rto swap toR]  ; (abc - cab)

def isDef' [not isNil]
def isNil' [eq []]
def drop2' [drop drop]
def drop3' [drop drop drop]