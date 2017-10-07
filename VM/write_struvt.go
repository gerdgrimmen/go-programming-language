package main

import
(
  "os"
  "fmt"
  "bytes"
)

type DefFunc struct {
  functionName int
  numArgs int
}

func main() {
  bytebuf := make([]byte,1)
  intbuf := make([]int,0)

  byts :=  []byte{2,3,4,5}
  reed := bytes.NewReader(byts)
  for{
    n, _ := reed.Read(bytebuf)
    if n == 0 { break }
    // write a chunk
    //code = append(code, bytebuf[0])
    intbuf = append(intbuf, int(bytebuf[0]))
  }
  runess := bytes.Runes(byts)
  fmt.Println("Runes: ",runess)
  fmt.Println(reed)
  fmt.Println(byts,intbuf)
  ints := byts
  fmt.Println(ints)
  var asd []uint8
  asdd := &byts
  asd = byts
  fmt.Println(asdd)
  fmt.Println(asd)


  thisbuf := bytes.NewBuffer([]byte{})
  n,err := thisbuf.WriteRune(321)
  if err != nil {
    panic(err)
  }
  fmt.Println(n)

  n,err = thisbuf.WriteString("ABC")
  if err != nil {
    panic(err)
  }
  r,s,err := thisbuf.ReadRune()
  if err != nil {
    panic(err)
  }
  fmt.Println(r,s)
  r,s,err = thisbuf.ReadRune()
  if err != nil {
    panic(err)
  }
  fmt.Println(r,s)
  r,s,err = thisbuf.ReadRune()
  if err != nil {
    panic(err)
  }
  fmt.Println(r,s)
  r,s,err = thisbuf.ReadRune()
  if err != nil {
    panic(err)
  }
  fmt.Println(r,s)


  // open output file
  fo, err := os.Create("written.runes")
  if err != nil { panic(err) }
  // close fo on exit and check for its returned error
    defer func() { if err := fo.Close(); err != nil { panic(err) } }()
  //if _, err := fo.Write(32); err != nil { panic(err) }
  /*
  fmt.Println(n)
  r,s,err := thisbuf.ReadRune()
  if err != nil {
    panic(err)
  }
  fmt.Println(r,s)
  r,s,err = thisbuf.ReadRune()
  if err != nil {
    panic(err)
  }
  fmt.Println(r,s)
  r,s,err = thisbuf.ReadRune()
  if err != nil {
    panic(err)
  }
  fmt.Println(r,s)
  //*/
 /*/
 funct := &DefFunc{4,3}
  fo, err := os.Create("written.struct")
  if err != nil {
    panic(err)
  }
  // close fo on exit and check for its returned error
  defer func() {
    if err := fo.Close(); err != nil {
        panic(err)
    }
  }()
  if _, err := fo.Write(byte(funct.functionName)); err != nil {
    panic(err)
  }
  */
}
