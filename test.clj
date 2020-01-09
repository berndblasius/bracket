; some test programs
;whl [gt 0 dup add 1] 1 -50000000


cond1 [[10 drop] [11 drop] [1]] trace 1

;cond1 [2] trace 1
def cond1' [
   eval if [b] [cond1 if car xs` b] isNil xs` 
   def xs' def b' car
]

def isNil' [eq []]

;; TRANSLATE SICP
; ****************************

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
;((acc 'deposit) 40)
; => 90
;((acc 'withdraw) 60)
;=> 30

; (define acc (make-account 50))

;acc withdraw' 60   ;m30 
;acc deposit' 40    ; 90
;def acc' make_acc 50
;def make-acc' [
;  lambda [
;
;  ]
;  def balance']

; THIS WORKS NOW !!
;acc withdraw' 60  ;30
;acc deposit' 40   ; 90
;def acc' \' make-acc 50
;
;def [make-acc balance'] 
;  [\ dispatch'
;   def [withdraw amount'] 
;     [eval if [balance def [balance`] sub balance amount]
;              [insuff']
;     gt balance amount]
;   def [deposit amount'] 
;     [balance def [balance`] add balance amount]
;   def [dispatch m'] 
;     [cond [[unknown']
;             [deposit] [eq m` deposit']
;             [withdraw] [eq m` withdraw']]]]





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

