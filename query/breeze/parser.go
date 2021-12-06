package breeze

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// Parser parses breeze queries.
// TODO: It is weird that this parser can only run once. I think we should be
// make this create actual parser instances or something that do expire, or
// something.
type Parser struct {
	input     string
	tokenizer Tokenizer
}

// NewParser creates a Parser.
func NewParser(input string) *Parser {
	return &Parser{
		input:     input,
		tokenizer: NewTokenizer(input),
	}
}

// ParseError is an error encountered during parsing.
type ParseError struct {
	error
	query    string
	position int
}

func maybeWrapInParseError(t Tokenizer, query string, err error) error {
	if err == nil {
		return nil
	}

	return &ParseError{
		error:    err,
		query:    query,
		position: t.Position(),
	}
}

// LocalError returns the deepest, original source error that caused this
// ParseError.
func (pe *ParseError) LocalError() error {
	err := pe.error
	for {
		nextErr := errors.Unwrap(err)
		if nextErr != nil {
			err = nextErr
			continue
		} else {
			// Otherwise, err was the innermost err.
			return err
		}
	}
}

// ErrorDescription returns a human-readable, more helpful message for a parse
// error.
func (pe *ParseError) ErrorDescription() string {
	localErrMsg := pe.LocalError().Error()
	var sb strings.Builder
	sb.WriteString(pe.query)
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat(" ", pe.position-1))
	sb.WriteString("^")
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat(" ", pe.position-1))
	sb.WriteString("| ")
	sb.WriteString(localErrMsg)
	sb.WriteString(" |")
	sb.WriteString("\n")

	return sb.String()
}

// Parse parses the query and returns stages for it.
func (p *Parser) Parse() ([]Stage, error) {
	stages, err := p.parse()
	return stages, maybeWrapInParseError(p.tokenizer, p.input, err)
}

func (p *Parser) parse() ([]Stage, error) {
	// A breeze query is a series of stages delimited by '|'. We keep parsing
	// these stages until we get nothing back.
	stages := []Stage{}
	for {
		stage, err := p.parseStage()
		if err == io.EOF {
			break // No more stages to parse.
		} else if err != nil {
			return nil, fmt.Errorf("failed to parse: %w", err)
		}
		stages = append(stages, stage)
		// Now chomp the stage separator and repeat.
		token := p.tokenizer.Next()
		if token == TokenEOF {
			break // No more stages to parse.
		}
		if token != TokenStageSeparator {
			return nil, fmt.Errorf("expected a stage separator after a stage, but got: %q", p.tokenizer.Text())
		}
	}
	return stages, nil
}

func (p *Parser) parseStage() (Stage, error) {
	stageToken := p.tokenizer.Next()
	if stageToken == TokenEOF {
		return nil, io.EOF
	}
	// Determine the stage type. If not recognized, this is a parse error.
	switch stageToken {
	case TokenFilter:
		return p.parseFilter()
	case TokenSort:
		return p.parseSort()
	case TokenGroup:
		return p.parseGroup()
	case TokenMap:
		return p.parseMap()
	default:
		return nil, fmt.Errorf("unrecognized stage: %q", p.tokenizer.Text())
	}
}

func (p *Parser) parseFilter() (*Filter, error) {
	uChecks := []*UnaryCheck{}
	bChecks := []*BinaryCheck{}
	for {
		if token, _ := p.tokenizer.Peek(); token == TokenStageSeparator || token == TokenEOF {
			break // No more checks to parse.
		}
		// TODO: I think our code might actually be a lot cleaner if we went with a single
		// Check type that covers both the unary & binary case.
		uCheck, bCheck, err := p.parseCheck()
		if err == io.EOF {
			break // No more checks to parse.
		} else if err != nil {
			return nil, fmt.Errorf("failed to parse check: %w", err)
		}

		switch {
		case uCheck != nil:
			uChecks = append(uChecks, uCheck)
		case bCheck != nil:
			bChecks = append(bChecks, bCheck)
		default:
			panic("unreachable, if err == nil => uCheck | bCheck cannot be nil")
		}
	}
	return &Filter{UnaryChecks: uChecks, BinaryChecks: bChecks}, nil
}

func (p *Parser) parseSort() (*Sort, error) {
	field, err := p.parseField()
	if err != nil {
		return nil, fmt.Errorf("failed to parse field: %w", err)
	}

	descending := p.parseSortOrder()

	return &Sort{
		Field:      field,
		Descending: descending,
	}, nil
}

func (p *Parser) parseGroup() (*Group, error) {
	var groupByFieldPtr *string
	if p.parseBy() {
		groupByField, err := p.parseField()
		if err != nil {
			return nil, fmt.Errorf("failed to parse the group-by field: %w", err)
		}

		groupByFieldPtr = &groupByField
	}

	aggFunc, err := p.parseAggFunc()
	if err != nil {
		return nil, fmt.Errorf("failed to parse aggregate function: %w", err)
	}

	aggregateField, err := p.parseField()
	if err != nil {
		return nil, fmt.Errorf("failed to parse field: %w", err)
	}

	return &Group{
		AggFunc:        *aggFunc,
		GroupByField:   groupByFieldPtr,
		AggregateField: aggregateField,
	}, nil
}

func (p *Parser) parseMap() (*Map, error) {
	assignments := []FieldAssignment{}
	for {
		if token, _ := p.tokenizer.Peek(); token == TokenStageSeparator || token == TokenEOF {
			break // No more checks to parse.
		}

		// TODO: I think our code might actually be a lot cleaner if we went with a single
		// Check type that covers both the unary & binary case.
		assignment, err := p.parseAssignment()
		if err == io.EOF {
			break // No more checks to parse.
		} else if err != nil {
			return nil, fmt.Errorf("failed to parse check: %w", err)
		}

		assignments = append(assignments, *assignment)
	}
	return &Map{Assignments: assignments}, nil
}

func (p *Parser) parseBy() bool {
	_, tokStr := p.tokenizer.Peek()

	if tokStr == "by" {
		_ = p.tokenizer.Next()
		return true
	}

	return false
}

func (p *Parser) parseCheck() (*UnaryCheck, *BinaryCheck, error) {
	field, err := p.parseField()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse field: %w", err)
	}

	uOp, bOp, err := p.parseOp()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse op: %w", err)
	}

	if uOp != "" {
		return &UnaryCheck{
			Field: field,
			Op:    uOp,
		}, nil, nil
	}

	value, err := p.parseConstValue()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse constant value: %w", err)
	}

	return nil, &BinaryCheck{
		Field: field,
		Value: value.(*Const), // TODO: We probably can and should get rid of this cast.
		Op:    bOp,
	}, nil
}

func (p *Parser) parseAssignment() (*FieldAssignment, error) {
	field, err := p.parseField()
	if err != nil {
		return nil, fmt.Errorf("failed to parse field: %w", err)
	}

	token := p.tokenizer.Next()
	if token == TokenEOF {
		return nil, errors.New("expected an equals (=), but reached end of query")
	}
	err = p.parseEquals(token)
	if err != nil {
		return nil, fmt.Errorf("failed to parse op: %w", err)
	}

	value, err := p.parseConstValue()
	if err != nil {
		return nil, fmt.Errorf("failed to parse constant value: %w", err)
	}

	return &FieldAssignment{
		Field: field,
		Assignment: ValueOrExpr{
			Value: value,
		},
	}, nil
}

func (p *Parser) parseField() (string, error) {
	token := p.tokenizer.Next()
	if token == TokenEOF {
		return "", errors.New("expected a field, but reached end of query")
	}

	// TODO: When we fix the fact that we probably shouldn't have token types for
	// each keyword, we can revert this back to TokenIdent, I think.
	if token == TokenIdent || token == TokenString {
		return p.tokenizer.Text(), nil
	}

	// TODO: We should make this error more obvious in its meaning. What we are
	// really trying to say is that we expected a field, but we got some other
	// kind of identifier here.
	// In general, I think many of our errors of of this kind ("expected X, but
	// got %q"). I think what we really wanna do is something along the lines of
	// "expected X, but got %q (%T)" (but without exposing internal types).
	return "", fmt.Errorf("expected a field identifier, but got %q (%s)", p.tokenizer.Text(), token.String())
}

func (p *Parser) parseOp() (UnaryOp, BinaryOp, error) {
	token := p.tokenizer.Next()
	if token == TokenEOF {
		return "", "", errors.New("expected an operator, but reached end of query")
	}

	if uOp, err := p.parseUnaryOp(token); err == nil {
		return uOp, "", nil
	} else if bOp, err := p.parseBinaryOp(token); err == nil {
		return "", bOp, nil
	}

	return "", "", fmt.Errorf("unrecognized operator: %q (%v)", p.tokenizer.Text(), token)
}

func (p *Parser) parseUnaryOp(token Token) (UnaryOp, error) {
	switch token {
	case TokenExists, TokenExistsNot:
		return UnaryOp(p.tokenizer.Text()), nil
	default:
		return "", fmt.Errorf("unrecognized unary operator: %q (%v)", p.tokenizer.Text(), token)
	}

}

func (p *Parser) parseBinaryOp(token Token) (BinaryOp, error) {
	switch token {
	case TokenGEQ, TokenEquals, TokenContains:
		return BinaryOp(p.tokenizer.Text()), nil
	default:
		return "", fmt.Errorf("unrecognized binary operator: %q (%v)", p.tokenizer.Text(), token)
	}
}

func (p *Parser) parseEquals(token Token) error {
	if token == TokenEquals {
		return nil
	}

	return fmt.Errorf("expected to find an equals (=), but found %v", token)
}

func isQuotedString(text string) bool {
	if len(text) < 2 {
		return false
	}

	return text[0] == '"' && text[len(text)-1] == '"'
}

func (p *Parser) parseConstValue() (Value, error) {
	token := p.tokenizer.Next()
	if token == TokenEOF {
		return nil, errors.New("expected a constant value, but reached end of query")
	}

	switch token {
	case TokenFloat, TokenInt:
		return &Const{
			Kind:        ConstKindNumber,
			Stringified: p.tokenizer.Text(),
		}, nil
	case TokenChar, TokenString:
		text := p.tokenizer.Text()
		if !isQuotedString(text) {
			return nil, fmt.Errorf("expected a properly quoted string, but got %q", text)
		}
		textWithoutQuotes := text[1 : len(text)-1]
		return &Const{
			Kind:        ConstKindString,
			Stringified: textWithoutQuotes,
		}, nil
	case TokenFalse, TokenTrue:
		return &Const{
			Kind:        ConstKindBool,
			Stringified: p.tokenizer.Text(),
		}, nil
	case TokenNull:
		return &Const{
			Kind:        ConstKindNull,
			Stringified: p.tokenizer.Text(),
		}, nil
	case TokenIdent:
		// Treat this as a string.
		return &Const{
			Kind:        ConstKindString,
			Stringified: p.tokenizer.Text(),
		}, nil
	default:
		return nil, fmt.Errorf("expected a constant value, but got: %q", p.tokenizer.Text())
	}
}

func (p *Parser) parseSortOrder() bool {
	maybeSortOrder, sortOrderText := p.tokenizer.Peek()
	if maybeSortOrder == TokenIdent {
		switch sortOrderText {
		case "asc":
			// If we get either asc/desc, then we want to _consume_ the token.
			p.tokenizer.Next()
			return false
		case "desc":
			p.tokenizer.Next()
			return true
		}
	}

	// If it isn't a sort keyword, then take the default behavior of ascending.
	return false
}

func (p *Parser) parseAggFunc() (*AggregateFunc, error) {
	tok := p.tokenizer.Next()
	aggFuncText := p.tokenizer.Text()
	if tok == TokenIdent {
		var aggFunc AggregateFunc
		switch aggFuncText {
		case "sum":
			aggFunc = AggFuncSum
		case "avg":
			aggFunc = AggFuncAvg
		case "count":
			aggFunc = AggFuncCount
		case "min":
			aggFunc = AggFuncMin
		case "max":
			aggFunc = AggFuncMax
		case "mode":
			aggFunc = AggFuncMode
		case "stddev":
			aggFunc = AggFuncStdDev
		default:
			return nil, fmt.Errorf("unrecognized aggregate function: %q", aggFuncText)
		}

		return &aggFunc, nil
	}

	return nil, fmt.Errorf("expected an aggregate func, but found %q", aggFuncText)
}
