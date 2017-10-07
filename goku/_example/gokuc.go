package main

import (
  "fmt"
  "os"
)

const (
  PUSH byte = 1
  ADD byte = 2
  PRINT byte = 3
  HALT byte = 4
  JMPLT byte = 5
)

var filename string

func main() {
  if len(os.Args) == 1 { /* os.Args[0] is "compiler" or "compiler.exe" */
    fmt.Println("P-Put it i-in me, Senpai!")
    return
  }else if len(os.Args) == 2 {
    filename = os.Args[1]
    fmt.Println("You put it in me, Senpai! <3", filename )
  }else if len(os.Args) > 2 { /* os.Args[0] is "main" or "main.exe" */
    fmt.Println("But Senpai, I can't Handle so many!")
    return
  }

  code := []byte{
    PUSH, 2,
    PUSH, 3,
    ADD,
    JMPLT, 10, 2,
    PRINT,
    HALT,
  }
  fmt.Println(code)
  // open output file
  fo, err := os.Create(filename+".goku")
  if err != nil {
    panic(err)
  }
  // close fo on exit and check for its returned error
  defer func() {
    if err := fo.Close(); err != nil {
        panic(err)
    }
  }()
  if _, err := fo.Write(code); err != nil {
    panic(err)
  }



}
