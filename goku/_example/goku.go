package main

import (
  "fmt"
  "io"
  "os"
)

const (
  PUSH byte = 1
  ADD byte = 2
  PRINT byte = 3
  HALT byte = 4
  JMPLT byte = 5
)

type op struct {
  name string
  nargs byte
}

var ops = map[byte]op{
  PUSH: op{"push", 1},
  ADD: op{"add", 0},
  PRINT: op{"print", 0},
  HALT: op{"halt", 0},
  JMPLT: op{"jmplt", 2},
}

type VM struct {
  code []byte
  pc byte
  stack []byte
  sp int
}

func (v *VM) trace(){
  addr := v.pc
  op :=  ops[v.code[v.pc]]
  args := v.code[v.pc+1:v.pc+op.nargs+1]
  stack := v.stack[0:v.sp+1]

  fmt.Printf("%04d: %s %v \t %v\n", addr,op.name, args, stack)
}

func (v *VM) run(){

  v.stack = make([]byte, 100)
  v.sp = -1
  v.pc = 0

  for{
    v.trace()
    //Fetch
    op := v.code[v.pc]
    v.pc++

    //Decode
    switch op {
    case PUSH:
      val := v.code[v.pc]
      v.pc++

      v.sp++
      v.stack[v.sp] = val
    case ADD:
      a := v.stack[v.sp]
      v.sp--
      b := v.stack[v.sp]
      v.sp--

      v.sp++
      v.stack[v.sp] = a + b
    case PRINT:
      val := v.stack[v.sp]
      v.sp--
      fmt.Println(val)
    case JMPLT:
      lt := v.code[v.pc]
      v.pc++
      addr := v.code[v.pc]
      v.pc++

      if v.stack[v.sp] < lt {
        v.pc = addr
      }
    case HALT:
      return

    }
  }
}
var filename string

func main() {
  if len(os.Args) == 1 { /* os.Args[0] is "compiler" or "compiler.exe" */
    fmt.Println("P-Put it i-in me, Senpai!")
    return
  }else if len(os.Args) == 2 {
    filename = os.Args[1]
    fmt.Println("You put it in me, Senpai! <3", filename )
  }else if len(os.Args) > 2 {
    fmt.Println("I cannot handle so many things, Senpai-sama! (-_-)" )
    return
  }

  v := &VM{}

  fi, err := os.Open(filename+".goku")
  if err != nil {
    panic(err)
  }
  // close fi on exit and check for its returned error
  defer func() {
    if err := fi.Close(); err != nil {
      panic(err)
    }
  }()
  bytebuf := make([]byte,1)
  for {
    // read a chunk
    n, err := fi.Read(bytebuf)
    if err != nil && err != io.EOF {
      panic(err)
    }
    if n == 0 { break }
    // write a chunk
    //code = append(code, bytebuf[0])
    v.code = append(v.code, bytebuf[0])
  }
  v.run()
}
