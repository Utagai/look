package breeze

import (
	"fmt"
	"strings"
	"text/scanner"
	"unicode"
)

// Token is a symbolic unit emitted by the tokenizer.
type Token int

// Various kinds of tokens in Liquid.
const (
	// Kind of ugly, but we know that scanner's tokens go in the negative range,
	// so we start with iota and go in the positive range and avoid conflicts
	// while avoiding the work of explicitly translating scanner tokens to breeze
	// tokens.
	// If we ever have issues with this, we can make everything off of iota and
	// write a conversion function.
	TokenStageSeparator Token = iota
	TokenFilter
	TokenSort
	TokenContains
	TokenEquals
	TokenGEQ
	TokenExists
	TokenExistsNot
	TokenFalse
	TokenTrue
	TokenNull
	TokenChar  Token = scanner.Char
	TokenFloat Token = scanner.Float

	TokenIdent  Token = scanner.Ident
	TokenInt    Token = scanner.Int
	TokenString Token = scanner.String
	TokenEOF    Token = scanner.EOF
)

// TODO: Should test this for exhaustiveness...
func (t Token) String() string {
	switch t {
	case TokenStageSeparator:
		return "StageSeparator"
	case TokenFilter:
		return "Filter"
	case TokenSort:
		return "Sort"
	case TokenContains:
		return "Contains"
	case TokenEquals:
		return "Equals"
	case TokenGEQ:
		return "GEQ"
	case TokenExists:
		return "Exists"
	case TokenExistsNot:
		return "Not Exists"
	case TokenFalse:
		return "False"
	case TokenTrue:
		return "True"
	case TokenNull:
		return "Null"
	case TokenChar:
		return "Char"
	case TokenFloat:
		return "Float"
	case TokenIdent:
		return "Ident"
	case TokenInt:
		return "Int"
	case TokenString:
		return "String"
	case TokenEOF:
		return "EOF"
	default:
		panic(fmt.Sprintf("unhandled token type: %d", t))
	}
}

const (
	// StageSeparatorString is the pipe character, which delimits the stages in a
	// breeze query.
	StageSeparatorString = "|"
)

type peeked struct {
	token    Token
	more     bool
	lastText string
}

// Tokenizer tokenizes an input string into breeze tokens.
// Under the hood, this is just a scanner.Scanner.
// In particular, it accomplishes two things we can't do with a scanner.Scanner:
//	* Give back tokens made for breeze in particular.
//  * Peek at token text.
type Tokenizer struct {
	s      scanner.Scanner
	peeked *peeked
}

// NewTokenizer returns a new Tokenizer.
func NewTokenizer(input string) Tokenizer {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))
	s.Mode = scanner.ScanChars |
		scanner.ScanFloats |
		scanner.ScanIdents |
		scanner.ScanInts |
		scanner.ScanStrings
	s.Error = func(s *scanner.Scanner, msg string) {}
	s.IsIdentRune = func(ch rune, i int) bool {
		if unicode.IsLetter(ch) {
			return true
		}

		switch ch {
		case '_': // Accept underscores in idents.
			return true
		case '!': // Accept ! in idents if they are at the beginning.
			return i == 0
		case '|': // Treat pipe as an identifier.
			return true
		}

		// Digits are OK, but only if they are after the first character.
		return i > 0 && unicode.IsDigit(ch)
	}

	return Tokenizer{
		s: s,
	}
}

// Next wraps scanner.Scanner#Next(). Returns the next token, as well as a
// boolean flag indicating if there are any more subsequent tokens. This
// function will never return a valid token (non-EOF) and simultaneously return
// false for the boolean flag. Therefore, iteration over this function can take
// the form:
//    for tok, ok := tokenizer.Next(); ok; tok, ok = tokenizer.Next() { ... }
func (t *Tokenizer) Next() (Token, bool) {
	if t.peeked != nil {
		token, more := t.peeked.token, t.peeked.more
		t.peeked = nil
		return token, more
	}

	tok, more := t.next()

	return tok, more
}

// Peek peeks at the next token to return, but does not advance the tokenizer
// past it.
func (t *Tokenizer) Peek() (Token, string, bool) {
	if t.peeked == nil {
		t.peeked = &peeked{}
		t.peeked.lastText = t.Text()
		t.peeked.token, t.peeked.more = t.next()
	}
	return t.peeked.token, t.s.TokenText(), t.peeked.more
}

func (t *Tokenizer) next() (Token, bool) {
	tok := t.s.Scan()
	if tok == scanner.EOF {
		// TODO: I think we can simplify this code by just returning TokenEOF, we
		// don't need the boolean.
		return TokenEOF, false
	}

	switch tok {
	case scanner.Ident:
		switch t.s.TokenText() {
		// Intercept keywords as special tokens.
		// TODO: This should be its own function.
		case "filter":
			return TokenFilter, true
		case "sort":
			return TokenSort, true
		case "contains":
			return TokenContains, true
		case "exists":
			return TokenExists, true
		case "!exists":
			return TokenExistsNot, true
		case StageSeparatorString:
			return TokenStageSeparator, true
		case "false":
			return TokenFalse, true
		case "true":
			return TokenTrue, true
		case "null":
			return TokenNull, true
		}
	default:
		switch t.s.TokenText() {
		// Intercept binary operators.
		case "=":
			return TokenEquals, true
		case ">":
			return TokenGEQ, true
		}
	}

	return Token(tok), true
}

// Text wraps scanner.Scanner#TokenText().
func (t *Tokenizer) Text() string {
	// Preserve this method's behavior of returning the current text, and not the
	// peeked text.
	if t.peeked != nil {
		return t.peeked.lastText
	}
	return t.s.TokenText()
}

// Position returns the current position in the string.
// This method technically returns the column position, which is equivalent
// because the string being tokenized is expected to be a single line.
func (t *Tokenizer) Position() int {
	return t.s.Position.Column
}
