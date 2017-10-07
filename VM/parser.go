package main

import
(
  "fmt"
  "strings"
  "bytes"
  "os"
  "io"
  "strconv"
)
const (
  PUSH byte = 1
  ADD byte = 2
  PRINT byte = 3
  HALT byte = 4
  JMPLT byte = 5
)


var string_code string

func main() {
  fi, err := os.Open("to_parse.goku")
  if err != nil {
    panic(err)
  }
  // close fi on exit and check for its returned error
  defer func() {
    if err := fi.Close(); err != nil {
      panic(err)
    }
  }()
  //*
  //Decode
  switchy := func(token string)byte{
    switch token {
      case "PUSH":
        return PUSH
      case "ADD":
        return ADD
      case "PRINT":
        return PRINT
      case "JMPLT":
        return JMPLT
      case "HALT":
        return HALT
      default:
        if s, err := strconv.Atoi(token); err == nil {
            //fmt.Printf("%T, %v", s, s)
            return byte(s)
        }
    }
    return 0
  }
 //*/
  bytebuf := make([]byte,1)
  code := make([]byte,0)
  for {
    // read a chunk
    n, err := fi.Read(bytebuf)
    if err != nil && err != io.EOF {
      panic(err)
    }
    if n == 0 { break }
    // write a chunk
    //code = append(code, bytebuf[0])
    code = append(code, bytebuf[0])
  }
  newbuf := bytes.NewBuffer(code)
  r := strings.NewReplacer(" ", "")
  string_code = r.Replace(newbuf.String())
  //string_code = newbuf.String()
  newcode := make([]byte,0)
  i := 0
  for _, row := range strings.Split(string_code, "\n") {
    if row != ""{
      for _, token := range strings.Split(row, ",") {
          if token != "" && switchy(token) != 0{
            // problem with strings.split() to give me an empty("") string but i somehow can't for it
            // changed my constants to start at 1 - using 0 as skip this "non-instruction" part
            fmt.Print(token)
            fmt.Println(switchy(token))
            newcode = append(newcode, switchy(token))
          }
      }
      //fmt.Printf("%04d: %s\n",i,row)
      i += 1
    }
  }
  fmt.Println(newcode)
  //*
  //*/

  // open output file
  fo, err := os.Create("parsed.goku")
  if err != nil { panic(err) }
  // close fo on exit and check for its returned error
    defer func() { if err := fo.Close(); err != nil { panic(err) } }()
  if _, err := fo.Write(newcode); err != nil { panic(err) }
}
