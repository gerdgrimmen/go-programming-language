package main

import (
  "bufio"
  "bytes"
  "io"
  "strings"
  "fmt"
)

type Token int

const (
  // Special Tokens
  ILLEGAL Token = iota
  EOF
  WS

  // Literals
  IDENTIFIER           // fields, table_name

  // Misc characters
  ASTERISK        // *
  COMMA           // ,
  BRACKETOPEN
  BRACKETCLOSE
  CODEBRACKETOPEN
  CODEBRACKETCLOSE
  EQUAL
  PLUS
  MINUS
  SLASH
  DOT
  AMPERSAND
  OR
  LESSTHAN
  GREATERTHAN
  COLON
  SEMICOLON


  // Keywords
  DEF
  IF
  ELSE
  PIPE
  RETURN
  MAIN
  IMPORT
  PACKAGE
  BREAK
  PRINT
  TYPE
  VAR
)

func isWhitespace(ch rune) bool{
  return ch ==' ' || ch =='\t' || ch =='\n'
}

func isLetter(ch rune) bool{
  return (ch >= 'a' && ch <='z') || (ch >= 'A' && ch <='Z')
}

func isDigit(ch rune) bool{
  return (ch >= 0 && ch <= 9)
}

var eof = rune(0)

type SelectStatement struct {
  Fields []string
  TableName string
}

type DefFunc struct {
  functionName string
  args []string
  numArgs int
  functionBody []string // TODO: Back to int
}

func (def *DefFunc) clean(){
  def.functionName = ""
  def.args = make([]string, 0)
  def.numArgs = 0
  def.functionBody = make([]string,0) // TODO: back to int
}


//Scanner represents a lexical scanner.
type Scanner struct {
  r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
  return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the buffered reader.
// returns the rune(0) if an error occurs( or io.EOF is returned)
func (s *Scanner) read() rune {
  ch, _, err := s.r.ReadRune()
    if err != nil { return eof }
  return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread(){ _ = s.r.UnreadRune()}

// Scn returns the next Token and literal Value
func (s *Scanner) Scan() (tok Token, lit string) {
  // Read the Next Rune
  ch := s.read()

  // if we see whitespace then consume all contiguous whitespace
  // if we see a letter then consume as an ident or reserved word
  if isWhitespace(ch) {
    s.unread()
    return s.scanWhitespace()
  }else if isLetter(ch) {
    s.unread()
    return s.scanIdent()
  }

  // otherwise read the individual character.
  switch ch {
    case eof:
      return EOF, ""
    case '*':
      return ASTERISK, string(ch)
    case ',':
      return COMMA, string(ch)
    case '{':
      return CODEBRACKETOPEN, string(ch)
    case '}':
      return CODEBRACKETCLOSE, string(ch)
    case '(':
      return BRACKETOPEN, string(ch)
    case ')':
      return BRACKETCLOSE, string(ch)
    }
  return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
  // Create a buffer and read the current character into it
  var buf bytes.Buffer
  buf.WriteRune(s.read())

  // Read every subsequent whitespace character into the buffer
  // Non-Whitespace characters and EOF will cause the loop to exit.
  for {
    if ch := s.read(); ch == eof {
      break
    } else if !isWhitespace(ch) {
      s.unread()
      break
    } else {
      buf.WriteRune(ch)
    }
  }

  return WS, buf.String()
}

// scanIdent consumes the urrent rune and all contiguous ident runes
func (s *Scanner) scanIdent() (tok Token, lit string){
  // Create a buffer and read the current character into it
  var buf bytes.Buffer
  buf.WriteRune(s.read())

  // Read every susequent ident character into the buffer
  // Non-ident characters and EOF will cause the loop to exit
  for {
    if ch := s.read(); ch == eof {
      break
    } else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
      s.unread()
      break
    } else {
      _, _ = buf.WriteRune(ch)
    }
  }

  // If the string matches a keyword then return
  switch strings.ToUpper(buf.String()) {
  case "DEF":
      return DEF, buf.String()
  }
  // Otherwise return as a regular identifier
  return IDENTIFIER, buf.String()
}

type Parser struct{
  s *Scanner
  buf struct{
    tok Token     // last read token
    lit string    // last read literal
    n int         // buffer size (max=1)
  }
}

// NewParser returns a new instance of Parser
func NewParser(r io.Reader) *Parser {
  return &Parser{s:NewScanner(r)}
}

// scan returns the next token from the underlying scanner
// if a token has been unscanned then read that instead
func (p *Parser) scan() (tok Token, lit string){
  // IF we have a token on the buffer, then return it.
  if p.buf.n != 0 {
    p.buf.n = 0
    return p.buf.tok, p.buf.lit
  }

  // Otherwise read the next token from the scan
  tok, lit = p.s.Scan()

  // Save it to the buffer in case we unscan later
  p.buf.tok, p.buf.lit = tok, lit

  return
}

// unscan pushes the previously read token back onto the buffer
func (p *Parser) unscan() {p.buf.n = 1}

// scanIgnoreWhitespace scans the next non-whitespace token
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
  tok, lit = p.scan()
  if tok == WS {
    tok, lit = p.scan()
  }
  return
}

func (p *Parser) Parse() ([]DefFunc, error) {
  functslice := []DefFunc{}
  funct := DefFunc{}

  for{
    if tok, lit := p.scanIgnoreWhitespace(); tok != DEF {
      return nil, fmt.Errorf("found %q, expected DEF", lit)
    }

    tok, lit := p.scanIgnoreWhitespace()
    if tok != IDENTIFIER {
      return nil, fmt.Errorf("found %q, expected functionName", lit)
    }
    funct.functionName = lit//= append(funct.functionName, lit)

    if tok, lit := p.scanIgnoreWhitespace(); tok != BRACKETOPEN {
      return nil, fmt.Errorf("found %q, expected BRACKETOPEN", lit)
    }

    i := 0
    for {
      // Read a field
      tok, lit := p.scanIgnoreWhitespace()
      if tok != IDENTIFIER && tok != ASTERISK {
        p.unscan()
        /*
        tok, _ := p.scanIgnoreWhitespace()
        if tok != BRACKETCLOSE {
          return nil, fmt.Errorf("found %q, expected BRACKETCLOSE", lit)
        }
        p.unscan()
        //*/
        break//return nil, fmt.Errorf("found %q, expected argument", lit)
      }
      funct.args = append(funct.args, lit)
      i++

      // If the next token is not a comma then break the loop.
      if tok,_ := p.scanIgnoreWhitespace();  tok != COMMA {
        p.unscan()
        break
      }
    }
    funct.numArgs = i
      if tok, lit := p.scanIgnoreWhitespace(); tok != BRACKETCLOSE {
        return nil, fmt.Errorf("found %q, expected BRACKETCLOSE", lit)
      }
    if tok, lit := p.scanIgnoreWhitespace(); tok != CODEBRACKETOPEN {
      return nil, fmt.Errorf("found %q, expected CODEBRACKETOPEN", lit)
    }
    /*
    functionBody: []int
    */

    //*/
    if tok, lit := p.s.scanIdent(); tok == CODEBRACKETCLOSE {
      return nil, fmt.Errorf("found %q, expected CODEBRACKETCLOSE", lit)
    }

    functslice = append(functslice, funct)
    funct.clean()
    if _, lit := p.scanIgnoreWhitespace(); lit != "" {
      p.unscan()
    }else{
      break
    }
  }
  return functslice, nil
}

func main() {
  pars := NewParser(strings.NewReader("DEF main(asd,dsa){} def asdasd(){} def aaaaaa(dddsss,sdsd){}"))
  functslice, err  := pars.Parse()
    if err != nil {fmt.Println(err);return}
    fmt.Println(functslice)
}
