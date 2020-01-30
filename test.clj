; some test programs

; ****************************************************************
; Ackermann function
; def ack(m: BigInt, n: BigInt): BigInt = {
;  if (m==0) n+1
;  else if (n==0) ack(m-1, 1)
;  else ack(m-1, ack(m, n-1))
;}
;
;(define (A m n)
;    (cond
;        ((= m 0) (+ n 1))
;        ((= n 0) (A (- m 1) 1))
;        (else (A (- m 1) (A m (- n 1))))))

;Clojure
;(defn ack [m n]
;  (cond (zero? m) (inc n)
;        (zero? n) (ack (dec m) 1)
;        :else (ack (dec m) (ack m (dec n)))))

; Haskell
; ack 0 n = n + 1
; ack m 0 = ack (m-1) 1
; ack m n = ack (m-1) (ack m (n-1))

;Factor
;:: ackermann ( m n -- u )
;    {
;        { [ m 0 = ] [ n 1 + ] }
;        { [ n 0 = ] [ m 1 - 1 ackermann ] }
;        [ m 1 - m n 1 - ackermann ackermann ]
;    } cond ;

;V
;[ack
;       [ [pop zero?] [popd succ]
;         [zero?]     [pop pred 1 ack]
;         [true]      [[dup pred swap] dip pred ack ack ]
;       ] when].

;V
;[ack
;       [ [pop zero?] [ [m n : [n succ]] view i]
;         [zero?]     [ [m n : [m pred 1 ack]] view i]
;         [true]      [ [m n : [m pred m n pred ack ack]] view i]
;       ] when].

;Forth
;: acker ( m n -- u )
;	over 0= IF  nip 1+ EXIT  THEN
;	swap 1- swap ( m-1 n -- )
;	dup  0= IF  1+  recurse EXIT  THEN
;	1- over 1+ swap recurse recurse ;

ack 3 4 def ack' \[m n]
    [cond
      [ [ack - m 1 ack m - n 1]
        [ack - m 1 1]  [eq 0 n]
        [+ n 1]  [eq 0 m]
      ]
    ]

;(defn ack [m n]
;  (cond (zero? m) (inc n)
;        (zero? n) (ack (dec m) 1)
;        :else (ack (dec m) (ack m (dec n)))))

; ****************************************************************
; Factorial

; -- Clojure
;(def factorial (fn [n]
;    (loop [cnt n acc 1]
;       (if (zero? cnt)
;            acc
;          (recur (dec cnt) (* acc cnt))))))

;--Scheme
;(define (factorial n)
;  (do ((i 1 (+ i 1))
;       (accum 1 (* accum i)))
;      ((> i n) accum)))

;(define (factorial n)
;  (let loop ((i 1)
;             (accum 1))
;    (if (> i n)
;        accum
;        (loop (+ i 1) (* accum i)))))

;--- Postscript
;/fact {
;  dup 0 eq     % check for the argument being 0
;  {
;    pop 1      % if so, the result is 1
;  }
;  {
;    dup 1 sub fact       % call recursively with n - 1
;    mul        % multiply the result with n
;  } ifelse
;} def


;simple recursive
;fac 4 def fac'
;      [cond [[* fac - swap 1 dup] [1 drop] [eq 1 dup]]]
 
;fac 4 def fac' \[n] [cond [[* fac - n 1 n] 1 [eq 1 n]]]

;recursive (tail-recursive)

; recursive with rec & anonymous function
;fac 4 def fac' \[n]
;   [drop swap eval \[acc cnt]
;     [rec < cnt n * acc cnt + cnt 1] 1 1]

;fac 4 def fac' \[n]
;   [drop swap eval [rec < cnt n * acc cnt + cnt 1 def [acc cnt]] 1 1]

;fac 4 def fac' [eval if [1 drop] [* fac - swap 1 dup] eq 1 dup ]



; ****************************************************************
; Fibonacci
;(defn- fib-iter [max n i j]
;  (if (= n max)
;    j
;    (recur max
;           (inc n)
;           j
;           (+ i j))))

;(defn fib [max]
;  (if (< max 2)
;    max
;    (fib-iter max 1 0N 1N)))

; ---------------

;(defn fib [n]
;  (case n
;    0 0
;    1 1
;    (+ (fib (- n 1)) (fib (- n 2)))))

;---- scheme
;
;--iterative
;define (fib n)
;  (let loop ((cnt 0) (a 0) (b 1))
;    (if (= n cnt)
;        a
;        (loop (+ cnt 1) b (+ a b)))))
;
;-- recursive
;(define (fib-rec n)
;  (if (< n 2)
;      n
;      (+ (fib-rec (- n 1))
;         (fib-rec (- n 2)))))

; tail recursive scheme
;(define (fib n)
;  (let loop ((a 0) (b 1) (n n))
;    (if (= n 0) a
;        (loop b (+ a b) (- n 1)))))

;fib 6 def fib' \[n][
;   loop 1 1 n
;   def loop' \[a b n][
;       eval if 
;           [a]
;           [loop b + a b - n 1]
;           eq n 0
;   ]
;]

;fib 6 def fib' [
;   loop def loop' \[a b n][
;       eval if 
;           [a]
;           [loop b + a b - n 1]
;           eq n 0
;   ] 1 1
;]

;fib 6 def fib' [
;   loop def loop' [
;       eval if 
;           [a]
;           [loop b + a b - n 1]
;           eq n 0 def [a b n]
;    ] 1 1
;]

;fib 6 def fib' \[max] [
; drop drop fib-iter 1 0 1 
; def fib-iter' [rec < n max + n 1 j + i j def [n i j]]
;]

;fib 6 def fib' \[max] [
; drop drop fib-iter 1 0 1 
; def fib-iter' \[n i j] [rec < n max + n 1 j + i j]]
;]


;fib 6 def fib' [eval if [] [+ fib - swap 1 swap fib - swap 2 dup] > 2 dup]
;fib 6 def fib' [eval if [n] [+ fib - n 1 fib - n 2] < n 2 def n']

;def fib' \[n] 
;    [cond [ [+ fib - n 1 fib - n 2]
;             n [< n 2]]] 

; ************************************************


;whl [gt 0 dup add 1] 1 -50000000

;cond1 [] trace 1
;cond1 [[3] [2] [1]] trace 1
;cond1 [[3] [1] [eq 4] [2 drop] [lt 4 dup]] 3

;cond1 [[10 drop] [11 drop] [1]] trace 1
;cond1 [[2]] trace 1

;def cond1' \[[xs x b]] 
;  [eval if [b] [eval if  x` [cond1 xs`] b] eq [] x`]

;; here still something wrong
;def cond1' \[[xs b]] 
;  [eval if [b] [eval if cons swap [cond] swap car xs` b] eq [] xs`]

;def isNil' [eq []]
; dip with aux
;def dip' [Rto eval toR swap]
;def dip' [x` eval def x' swap]

; dip as a macro
;def dip' [cons eval' cons rot]

;cl.x'
;def [cl.x] 4
;\[] [def x' 5]

;; TRANSLATE SICP
; ****************************

;acc withdraw' 60
;acc deposit' 100
;acc withdraw' 60
;acc withdraw' 60
;acc deposit' 40
;def acc' make-acc 50 
;def make-acc' \[balance] [
;    dispatch` def dispatch' \[m][
;   ;\[m][
;       cond [
;           [unknown']
;           deposit [eq m deposit'] 
;           withdraw [eq m withdraw'] 
;       ]
;   ]
;   def withdraw' \[amount][
;       eval if 
;         [balance def [balance`] - balance amount]
;         [insuff']
;         gt balance amount
;   ] 
;   def deposit' \[amount][
;       balance def [balance`] + balance amount
;   ]
;]

; alternative shorter
;acc withdraw' 60
;acc deposit' 100
;acc withdraw' 60
;acc withdraw' 60
;acc deposit' 40
;def acc' make-acc 50 
;;def make-acc' \[balance] [
;def make-acc' [
;    \[m][ cond [
;           [unknown']
;           deposit [eq m deposit'] 
;           withdraw [eq m withdraw'] 
;          ]
;   ]
;   def withdraw' [
;       eval if 
;         [balance def [balance`] - balance]
;         [insuff' drop]
;         gt balance dup 
;   ] 
;   def deposit' [balance def [balance`] + balance]
;   def balance'
;]

;(define (make-account balance)
;  (define (withdraw amount)
;    (if (>= balance amount)
;        (begin (set! balance (- balance amount))
;               balance)
;        "Insufficient funds"))
;  (define (deposit amount)
;    (set! balance (+ balance amount))
;    balance)
;  (define (dispatch m)
;    (cond ((eq? m 'withdraw) withdraw)
;          ((eq? m 'deposit) deposit)
;          (else (error "Unknown request -- MAKE-ACCOUNT"
;                       m))))
;  dispatch)
; (define acc (make-account 50))
;((acc 'deposit) 40)
; => 90
;((acc 'withdraw) 60)
;=> 30





; ****************************
;(define (make-withdraw balance)
;  (lambda (amount)
;    (if (>= balance amount)
;        (begin (set! balance (- balance amount))
;               balance)
;        "Insufficient funds")))
;(define W1 (make-withdraw 100))
;(W1 50)

;  Bracket (in go)
;W1 61 W1 10 W1 30 
;def W1' account 100 
;def account' 
;    [lambda [eval if [balance set balance' sub balance amount]
;              [insuff']
;               gt balance amount def amount'] 
;               def balance']

; -------------------
; THIS WORKS NOW (bracket old julia)!!
;W1 61 W1 10 W1 30 
;def W1' \' account 100 trace 1
;def [account balance'] 
;   [ \ def [amount'] 
;           [eval if [balance def [balance`] sub balance amount]
;                    [insuff']
;           gt balance amount]]
; -------------------

