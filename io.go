package main

/* io functions for bracket */


import (
    "fmt"
    "bytes"
    "errors"
    "io/ioutil"
    "strconv"
)

/* compared to Base64 we place the digits at the beginning 
and use 'minus' and 'underscore' as additional chars  
0..9 have position 0..9, 'A' .. 'Z' have position 10..35
'a' .. 'z' have position 36..61, '-' has position 62 and '_' has 63*/
const base64_enc_table = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"

// position in map is ASCII index, value is the index in Base64
var base64_dec_map = map[byte]int{
    '0':0, '1':1, '2':2, '3':3, '4':4, '5':5, '6':6, '7':7, '8':8, '9':9,
    'A':10,'B':11,'C':12,'D':13,'E':14,'F':15,'G':16,'H':17,'I':18,'J':19,
    'K':20,'L':21,'M':22,'N':23,'O':24,'P':25,'Q':26,'R':27,'S':28,'T':29,
    'U':30,'V':31,'W':32,'X':33,'Y':34,'Z':35,'a':36,'b':37,'c':38,'d':39,
    'e':40,'f':41,'g':42,'h':43,'i':44,'j':45,'k':46,'l':47,'m':48,'n':49,
    'o':50,'p':51,'q':52,'r':53,'s':54,'t':55,'u':56,'v':57,'w':58,'x':59,
    'y':60,'z':61,'-':62,'_':63, 
}

func string2symbol(str string) value {
// encode each character into 6 bits Base64-value
// only 10*6=60 bits are used,
// the remaining 4 bits can be used either as flags bits (at the right side)
   x := 0
   nok := 0
   for i, _ := range str {
       x1, ok := base64_dec_map[str[i]]
       if ok {
            x  = ( x << 6 ) | x1
            nok += 1
            if nok == 10 {break}
        }
   }
   return boxSymb(x)
}

func symbol2string(symb value) string {
// decode symbol back to string for output
   x := unbox(symb)  // remove the flag bits
   s := make([]byte, 10)
   for i:=9; i>=0; i-- {
       s[i] = base64_enc_table[x & 63]
       x = x>>6
   }
   for i:=0;i<10;i++{ // remove trailing 0's
      if s[i] != '0' {
        return string(s[i:10])
      }
   }
   return ""
}

func (vm *Vm) printElem(q value) {
   switch {
   case isInt(q):
       fmt.Print(unbox(q))
   //case isFloat(q):
   //    fmt.Print(unbox_float(q))
   case isNil(q):
        fmt.Print("[]")
   case isPrim(q):
        fmt.Print(primStr[q])
   case isSymb(q):
        fmt.Print(symbol2string(q))
   default:
        vm.printList(q)
   }
}

func (vm *Vm) printInnerList(l value, invert bool) {
   if isDef(l) {
      if invert {
          l=vm.reverse(l)
      }
      vm.printElem(vm.car(l))
      l = vm.cdr(l)
      if (!isCons(l) && isDef(l)){   // dotted list
            fmt.Print(" . ")
            vm.printElem(l)
       } else { 
          for isDef(l) {
             fmt.Print(" ")
             vm.printElem(vm.car(l))
             l = vm.cdr(l)
          }
      }
  }
}

func (vm *Vm) printList(l value) {
      fmt.Print("[")
      vm.printInnerList(l,true)
      fmt.Print("]")
}

func (vm *Vm) printKet(l value) {
      fmt.Print("[")
      vm.printInnerList(l,false)
      fmt.Println(">")
}

func (vm *Vm) printBra(l value) {
      fmt.Print("<")
      vm.printInnerList(l,true)
      fmt.Println("]")
}

func parse(token []byte) (value, error) {
    if n, err := strconv.Atoi(string(token)); err == nil {
      return boxInt(n), nil
    } 
    p,ok := str2prim[string(token)]
    if ok {
       return p, nil   // token is a primitive
    } else {
       return string2symbol(string(token)), nil  // token is a symbol
    }
   return nill, errors.New("parse error, token not found")
}

func (vm *Vm) readFromTokens(tokens [][]byte, pos int) (value, int) {
  s := nill
  s1 := nill
  for pos < len(tokens){
    token := tokens[pos]
    pos++
    switch string(token) { 
    case "]" :
      return s, pos
    case "[":
      s1, pos = vm.readFromTokens(tokens, pos)
      s = vm.cons(s1,s)
    default:
      p,err := parse(token)
      if err == nil {
          s = vm.cons(p,s)
       } else {
       fmt.Println("parse error")
      }
    }
  }
  return s, pos
}

func removeComments (str []byte ) []byte {
    l := len(str)
    newstr := make([]byte, l)
    j := 0
    for i:=0; i<l; i++ {
       switch str[i] {
       case ';':
           for str[i] != '\n' && i != l-1 {
              i += 1
           }
       case '\n':
          newstr[j] = ' '
          j += 1
       default:
          newstr[j] = str[i]
          j += 1
      }
    }
    return newstr[0:j]
}

func tokenize(str []byte) [][]byte {
   str = bytes.ReplaceAll(str, []byte("["), []byte(" [ "))
   str = bytes.ReplaceAll(str, []byte("]"), []byte(" ] "))
   str = bytes.ReplaceAll(str, []byte("'"), []byte(" esc "))
   str = removeComments(str)
   return bytes.Fields(str)
}

func (vm *Vm) makeBra(prog string) value {
    tokens := tokenize([]byte(prog))
    val,_ := vm.readFromTokens(tokens, 0)
    return val
}


func (vm *Vm) loadFile(fname string) value{
    //file,_ := os.Open(prog)
    //r := bufio.NewReader(file)
    b, _ := ioutil.ReadFile(fname)
    tokens := tokenize(b)
    val,_ := vm.readFromTokens(tokens, 0)
    return val
}


// so far, printList of dotted lists works only for
// dotted pairs (because not clear how to define reverse..
// --> maybe remove print of dotted pairs totally??
