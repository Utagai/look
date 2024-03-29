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

	// HACK: OK, so this is technically wrong. In particular, we should not be
	// making tokens for every 'keyword' in the language here. We should only be
	// identifying idents here, and the parser level that has more context can
	// decide how to interpret the ident, namely for the sake of determining if a
	// given ident is a field name or an actual keyword.

	// Stages:
	TokenStageSeparator Token = iota
	TokenFilter
	TokenSort
	TokenGroup
	TokenMap

	// Punctuators
	TokenLParen
	TokenRParen
	TokenLSqBracket
	TokenRSqBracket
	TokenComma

	// Binary comparison operators:
	TokenEquals
	TokenGEQ
	TokenContains

	// Binary expression operations:
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide

	// Keyword consts:
	TokenFalse
	TokenTrue
	TokenNull

	// Scanner package token types:
	TokenChar   Token = scanner.Char
	TokenFloat  Token = scanner.Float
	TokenIdent  Token = scanner.Ident
	TokenInt    Token = scanner.Int
	TokenString Token = scanner.String
	TokenEOF    Token = scanner.EOF
)

func (t Token) String() string {
	switch t {
	case TokenStageSeparator:
		return "StageSeparator"
	case TokenFilter:
		return "Filter"
	case TokenSort:
		return "Sort"
	case TokenGroup:
		return "Group"
	case TokenMap:
		return "Map"
	case TokenLParen:
		return "LParen"
	case TokenRParen:
		return "RParen"
	case TokenComma:
		return "Comma"
	case TokenLSqBracket:
		return "LSqBracket"
	case TokenRSqBracket:
		return "RSqBracket"
	case TokenContains:
		return "Contains"
	case TokenEquals:
		return "Equals"
	case TokenGEQ:
		return "GEQ"
	case TokenPlus:
		return "Plus"
	case TokenMinus:
		return "Minus"
	case TokenMultiply:
		return "Multiply"
	case TokenDivide:
		return "Divide"
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
		panic(fmt.Sprintf("unhandled token type: %d (%v)", t, scanner.TokenString(rune(t))))
	}
}

const (
	// StageSeparatorString is the pipe character, which delimits the stages in a
	// breeze query.
	StageSeparatorString = "|"
)

type peeked struct {
	token    Token
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
		case '.': // If we see a '.' and it is the first position of the token, treat it as ident.
			return i == 0
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
// function will return TokenEOF when it has exhausted all tokens in the input:
//    for tok := tokenizer.Next(); tok != breeze.TokenEOF; tok = tokenizer.Next() { ... }
func (t *Tokenizer) Next() Token {
	if t.peeked != nil {
		token := t.peeked.token
		t.peeked = nil
		return token
	}

	tok := t.next()

	return tok
}

// Peek peeks at the next token to return, but does not advance the tokenizer
// past it.
func (t *Tokenizer) Peek() (Token, string) {
	if t.peeked == nil {
		lastText := t.Text()
		t.peeked = &peeked{}
		t.peeked.lastText = lastText
		t.peeked.token = t.next()
	}
	return t.peeked.token, t.s.TokenText()
}

func (t *Tokenizer) next() Token {
	tok := t.s.Scan()
	if tok == scanner.EOF {
		return TokenEOF
	}

	switch tok {
	case scanner.Ident:
		return t.convertIdentToken(tok)
	default:
		switch t.s.TokenText() {
		// Intercept binary operators.
		case "=":
			return TokenEquals
		case ">":
			return TokenGEQ
		case "(":
			return TokenLParen
		case ")":
			return TokenRParen
		case "[":
			return TokenLSqBracket
		case "]":
			return TokenRSqBracket
		case ",":
			return TokenComma
		case "+":
			return TokenPlus
		case "-":
			return TokenMinus
		case "*":
			return TokenMultiply
		case "/":
			return TokenDivide
		}
	}

	return Token(tok)
}

// Converts a token from the scanner into a breeze-specific Token type, if
// possible.
func (t *Tokenizer) convertIdentToken(tok rune) Token {
	switch t.s.TokenText() {
	case StageSeparatorString:
		return TokenStageSeparator
	case "filter":
		return TokenFilter
	case "sort":
		return TokenSort
	case "group":
		return TokenGroup
	case "map":
		return TokenMap
	case "contains":
		return TokenContains
	case "false":
		return TokenFalse
	case "true":
		return TokenTrue
	case "null":
		return TokenNull
	default:
		return Token(tok)
	}
}

// Text wraps scanner.Scanner#TokenText().
// In other words, returns the current/last-returned token.
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
