package main

/* io functions for bracket */


import (
    "fmt"
   // "bufio"
   "bytes"
   "errors"
    "io/ioutil"
    //"os"
   // "strings"
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
   return box_symb(x)
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

func (vm *Vm) printElem(q value, invert bool) {
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
        vm.printList(q,invert)
   }
}

func (vm *Vm) printList(l value, invert bool) {
   if isNil(l) {
       print("[]")
   } else {
      if invert {
          //l=vm.reverseL(l)
      }
      fmt.Print("[")
      vm.printElem(vm.car(l),invert)
      l = vm.cdr(l)
      for isDef(l) {
          fmt.Print(" ")
          vm.printElem(vm.car(l), invert)
          l = vm.cdr(l)
      }
      fmt.Print("]")
  }
}

func (vm *Vm) printKet(l value) {
   if isNil(l) {
       print("[ >")
   } else {
      fmt.Print("[")
      vm.printElem(vm.car(l),false)
      l = vm.cdr(l)
      for isDef(l) {
          fmt.Print(" ")
          vm.printElem(vm.car(l), false)
          l = vm.cdr(l)
      }
      fmt.Print(">")
  }
  fmt.Println()
}


/*
void printList(any l, int invert, Vm *e) 
{
    if isNil(l) 
        printf("[]");
    else {
       if (invert) l=reverse(l, e);
       printf("[");
       printElem(car(l), invert, e);
       cdr_(&l);
       if (!isCons(l) && !isNil(l)){   // dotted list
            printf(" . ");
            printElem(l,invert,e);
       } else while(!isNil(l)) {
          printf(" ");
          printElem(car(l), invert, e);
          cdr_(&l);
       }
       printf("]");
   }
}
*/

func isSpace(c byte) bool {
   return c == ' ' || c == '\t'
}

//func nextchar(io *bufio.Reader) (byte, error) {
func nextchar(io *bytes.Reader) (byte, error) {
   c := byte('0')

   for {
      c, err := io.ReadByte()
      if err != nil { 
          return 0,err
      }
      if c == ';' {  // single line comment
         for {
            c, err = io.ReadByte()
            if err != nil { 
              return 0, err
            }
            if c == '\n' {
                break
            }
        }
     }
    if !isSpace(c){ 
        break
    }
   }

   return c,nil
}

//func read_token(r *bufio.Reader) byte {
func read_token(r *bytes.Reader, buf *bytes.Buffer) byte {
    token := byte(' ')
    ndots := 0  // number of dots in token
    for {
      c, err := r.ReadByte()
      if err != nil { 
          break
      }
      if c == '.' {
        ndots += 1
    }
      if ndots > 1 {// can have at most 1 dot in number
         r.UnreadByte()
         break
      }
      switch { 
      case c == '\n' || isSpace(c):
         break
      case !('0'<=c && c<'9' || c=='+' || c=='-')  :
         r.UnreadByte()
         break
      //case !(c=='.' || c == '*' || c == '/' || isletter(c)):
      //   r.UnreadByte()
      //   break
       case c=='`' || c == '\'' || c == '\\' || c == ';':
         r.UnreadByte()
         break
       default:
        //push!(buf,c)
        buf.WriteByte(c)
        //token = token * c
      }
   }
   return token
}

//func (vm *Vm) read_tokens(r *bufio.Reader) value {
func (vm *Vm) read_tokens(r *bytes.Reader) value {
    val := NIL
    var buf bytes.Buffer
    for {
        c,err := nextchar(r)
        if err != nil {
            return val  // end of stream
        }
        switch c {
        case ']':    // end of list
           return val
        case '[':    // begin of list
           newval := vm.read_tokens(r)
           val = vm.cons(newval, val)
         //elseif c == '\''   // escape
         //  val = cons!(ESC, val, vm)
         //elseif c == '`'   // Asc (escape value)
         //  val = cons!(ASC, val, vm)
         //elseif c == '\\'   // Backslash = lambda
         //  val = cons!(LAMBDA, val, vm)
         default:               // read new atom
           r.UnreadByte()
           read_token(r, &buf)
           //newval = atom(token)
              newval := NIL
           val = vm.cons(newval, val)
         }
    }
    return NIL
}

func parse(token []byte) (value, error) {
    if n, err := strconv.Atoi(string(token)); err == nil {
      return box_int(n), nil
    } 
    p,ok := str2prim[string(token)]
    if ok {
       return p, nil   // token is a primitive
    } else {
       return string2symbol(string(token)), nil  // token is a symbol
    }
   return NIL, errors.New("parse error, token not found")
}

func (vm *Vm) readFromTokens(tokens [][]byte, pos int) (value, int) {
  s := NIL
  s1 := NIL
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
    return newstr
}

func tokenize(str []byte) [][]byte {
   str = bytes.ReplaceAll(str, []byte("["), []byte(" [ "))
   str = bytes.ReplaceAll(str, []byte("]"), []byte(" ] "))
   str = bytes.ReplaceAll(str, []byte("'"), []byte(" esc "))
   str = removeComments(str)
   return bytes.Fields(str)
}


func (vm *Vm) makeBra(prog string) value {
    //tokens := tokenize([]byte("[ " + prog))
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
