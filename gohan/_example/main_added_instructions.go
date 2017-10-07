package main

import "fmt"
/*
Still needed:
Globaldata
Localdata
Datatype - Stack/List, Strings, "Goroutines"/Jobs/Workers, atoms/symbols

"keys" for variable access

Packages/Modules extend specific/current Instructionset / Functionset
Export to Go-Package
Make it Buildable with Go -> 1 .exe File

Different Ideas
Interactive Shell
//*/

/*
Possible:
DECLARE Variable (local)
ACCESS Variable (local)
LOOP
IF SKIP Function --- JUMP !?

//*/
const (
  PUSH = iota
  ADD
  PRINT
  HALT
  JMPLT
  JMPGT
  SUB
  DIV
  MULT
  MODU
  STORE
  POP
  CALL // CALL another []byte - put parameters/args on stack, put numargs on stackset pointer on right position and return to last address after execution
  JUMP
  GSTORE
  GLOAD
  RET
  PIPE
  LOAD
  SET
  REC
)

var FUNCS [][]int=[][]int{
                      []int{        // after call v.run(FUNCS[val]) set v.code again to code
                                    // and set the pointer back to the state before the call
                          PUSH, 600,
                          DIV,2,
                          MODU,170,
                          MULT,2,
                          PRINT,
                          STORE,

                          //*
                          CALL, 2, 2,5,9,
                          PUSH, 4,
                          ADD,
                          PRINT,
                          STORE,
                          PIPE, 3, 2,0,7,
                          PUSH, 4,
                          ADD,
                          JMPLT,25,25,
                          HALT,
                          //*/
                        },
                        []int{
                            PUSH,4,
                            LOAD, -3,
                            LOAD, -4,
                            CALL, 1,0,
                            ADD,
                            JMPLT, 7,2,
                            PRINT,
                            STORE,
                            RET,
                          },
                         []int{
                            PUSH,6,
                            PUSH,3,
                            SUB,
                            PRINT,
                            STORE,
                            RET,
                          },
                          []int{
                            REC,
                            LOAD, -3,
                            JMPLT, 1, 21,
                            JMPGT, 0, 9,
                            PUSH, 1,
                            SUB,
                            SET, -3,
                            LOAD, -3,
                            JMPGT, 0,0,
                            PRINT,
                            STORE,
                            RET,
                          },
                          []int{
                            PUSH, 2,
                            CALL, 3, 2,5,3,
                            PRINT,
                            HALT,
                          },
                        }
/*
const  PUSH = 0
const  ADD = 1
const  PRINT = 2
const  HALT = 3
const JMPLT = 4
//*/
type op struct {
  name string
  nargs int
}

var ops = map[int]op{
  PUSH: op{"push", 1},
  ADD: op{"add", 0},
  PRINT: op{"print", 0},
  HALT: op{"halt", 0},
  JMPLT: op{"jl", 2},
  JMPGT: op{"jg", 2},
  SUB: op{"subs",0},
  DIV: op{"divi",1},
  MULT: op{"mult",1},
  MODU: op{"modulo",1},
  STORE: op{"store", 0},
  POP: op{"pop",0},
  CALL: op{"call",2}, // function name and numargs
  JUMP: op{"jump", 1}, // jump to instruction
  GSTORE: op{"gstore",1}, // takes current "HEAD" and puts it in global address
  GLOAD: op{"gload", 1}, // puts on stack the data on given global address
  RET: op{"ret", 0},
  PIPE: op{"pipe", 2},
  LOAD: op{"load",1},
  SET: op{"set",1},
  REC: op{"set",0},
}


type VM struct {
  code *[]int
  scope int
  pc int
  stack []int
  sp int
  fp int
}

func (v *VM) trace(){
  addr := v.pc
  op :=  ops[(*v.code)[v.pc]]
  args := (*v.code)[v.pc+1:v.pc+op.nargs+1]
  stack := v.stack[0:v.sp+1]

  fmt.Printf("%02d: %s %v \t\t %v\n", addr,op.name, args, stack)
}

func (v *VM) run(c []int){

  v.stack = make([]int, 100)
  v.sp = -1
  v.fp = 0
  v.code = &c
  v.pc = 0
  v.scope = 0
  for{
    v.trace()
    //Fetch
    op := (*v.code)[v.pc]
    v.pc++

    //Decode
    switch op {
    case PUSH:
      val := (*v.code)[v.pc]
      v.pc++
      v.sp++
      v.stack[v.sp] = val
    case ADD:               // TODO: Shorten ADD - SUB - MULT - DIVI - MODU
      v.stack[v.sp-1] += v.stack[v.sp]
      v.sp--
      /*
      a := v.stack[v.sp]
      v.sp--
      b := v.stack[v.sp]
      v.sp--
      v.sp++
      v.stack[v.sp] = a + b
      //*/
    case PRINT:
      val := v.stack[v.sp]
      v.sp--
      fmt.Println(val)
    case JMPLT:
      lt := (*v.code)[v.pc]
      v.pc++
      addr := (*v.code)[v.pc]
      v.pc++

      if v.stack[v.sp] < lt {
        v.pc = addr
      }
    case JMPGT:
      gt := (*v.code)[v.pc]
      v.pc++
      addr := (*v.code)[v.pc]
      v.pc++

      if v.stack[v.sp] > gt {
        v.pc = addr
      }
    case HALT:
      return
    case SUB:
      v.stack[v.sp-1] -= v.stack[v.sp]
      v.sp--
    case DIV:
      v.stack[v.sp] /= (*v.code)[v.pc]
      v.pc++
    case MULT:
      v.stack[v.sp] *= (*v.code)[v.pc]
      v.pc++
    case MODU:
      v.stack[v.sp] %= (*v.code)[v.pc]
      v.pc++
    case STORE:
      v.sp++
    case POP:
      v.sp--
    case CALL:
      funct := (*v.code)[v.pc]
      v.pc++
      numArgs := (*v.code)[v.pc]
      v.pc++
      v.sp++
      v.stack[v.sp] = v.scope
      v.scope = funct
      for i := numArgs;i >= 1;i--{
        v.sp++
        v.stack[v.sp] = (*v.code)[v.pc]
        v.pc++
      }
      v.sp++
      v.stack[v.sp] = numArgs
      v.sp++
      v.stack[v.sp] = v.pc
      v.sp++
      v.stack[v.sp] = v.fp
      v.pc = 0
      v.fp = v.sp
      v.code = &FUNCS[funct]
    case RET:
      retVal := v.stack[v.sp]
      v.sp = v.fp
      v.fp = v.stack[v.sp]
      v.sp--
      v.pc = v.stack[v.sp]
      v.sp--
      numArgs := v.stack[v.sp]
      for i := numArgs;i >= 1;i--{
        v.sp--
      }
      v.sp--
      scope := v.stack[v.sp]
      v.scope = scope
      v.stack[v.sp] = retVal
      v.code = &FUNCS[scope]
    case PIPE:
      piped := v.stack[v.sp]
      funct := (*v.code)[v.pc]
      v.pc++
      numArgs := (*v.code)[v.pc]
      v.pc++
      v.stack[v.sp] = v.scope
      v.sp++
      v.stack[v.sp] = piped
      v.pc++
      v.scope = funct
      for i := numArgs;i >= 2;i--{
        v.sp++
        v.stack[v.sp] = (*v.code)[v.pc]
        v.pc++
      }
      v.sp++
      v.stack[v.sp] = numArgs
      v.sp++
      v.stack[v.sp] = v.pc
      v.sp++
      v.stack[v.sp] = v.fp
      v.pc = 0
      v.fp = v.sp
      v.code = &FUNCS[funct]

// CALL another []byte - put parameters/args on stack, put numargs on stack
// set pointer on right position and return to last address after execution
    case JUMP:
      val := (*v.code)[v.pc] // TODO: shorten this(JUMP) Instruction
      v.pc = val
    case GSTORE: // TODO: implement GSTORE Instruction
      return
    case GLOAD: // TODO: implement GLOAD Instruction
      return
    case LOAD:
      adr := (*v.code)[v.pc]
      v.pc++
      val := v.stack[v.fp+adr]
      v.sp++
      v.stack[v.sp] = val
    case  SET:
      adr := (*v.code)[v.pc]
      v.pc++
      val := v.stack[v.sp]
      v.stack[v.fp+adr] = val
    case  REC:
      v.sp = v.fp
    }
  }
}

func main() {
  v := &VM{}
  v.run(FUNCS[v.scope])
}
