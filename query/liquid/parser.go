package liquid

import (
	"errors"
	"fmt"
	"io"
)

// Parses liquid queries.
type Parser struct {
	input     string
	tokenizer Tokenizer
}

func NewParser(input string) *Parser {
	return &Parser{
		input:     input,
		tokenizer: NewTokenizer(input),
	}
}

func (p *Parser) Parse() ([]Stage, error) {
	// A liquid query is a series of stages delimited by '|'. We keep parsing
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
		token, more := p.tokenizer.Next()
		if !more {
			break // No more stages to parse.
		}
		if token != TokenStageSeparator {
			return nil, fmt.Errorf("expected a stage separator after a stage, but got: %q", p.tokenizer.Text())
		}
	}
	return stages, nil
}

func (p *Parser) parseStage() (Stage, error) {
	stageToken, more := p.tokenizer.Next()
	if !more {
		return nil, io.EOF
	}
	// Determine the stage type. If not recognized, this is a parse error.
	switch stageToken {
	case TokenFind:
		return p.parseFind()
	case TokenSort:
		return p.parseSort()
	default:
		return nil, fmt.Errorf("unrecognized stage: %q", p.tokenizer.Text())
	}
}

func (p *Parser) parseFind() (*Find, error) {
	checks := []*Check{}
	for {
		if token, _, _ := p.tokenizer.Peek(); token == TokenStageSeparator {
			break // No more checks to parse.
		}
		check, err := p.parseCheck()
		if err == io.EOF {
			break // No more checks to parse.
		} else if err != nil {
			return nil, fmt.Errorf("failed to parse check: %w", err)
		}
		checks = append(checks, check)
	}
	return &Find{Checks: checks}, nil
}

func (p *Parser) parseCheck() (*Check, error) {
	field, err := p.parseField()
	if err != nil {
		return nil, fmt.Errorf("failed to parse field: %w", err)
	}

	op, err := p.parseOp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse op: %w", err)
	}

	value, err := p.parseConstValue()
	if err != nil {
		return nil, fmt.Errorf("failed to parse constant value: %w", err)
	}

	return &Check{
		Field: field,
		Value: value,
		Op:    op,
	}, nil
}

func (p *Parser) parseField() (string, error) {
	token, more := p.tokenizer.Next()
	if !more {
		return "", errors.New("expected a field, but reached end of query")
	}

	if token == TokenIdent {
		return p.tokenizer.Text(), nil
	}
	return "", fmt.Errorf("expected a field identifier, but got %q", p.tokenizer.Text())
}

func (p *Parser) parseOp() (BinaryOp, error) {
	token, more := p.tokenizer.Next()
	if !more {
		return "", errors.New("expected a binary operator, but reached end of query")
	}

	switch token {
	case TokenGEQ, TokenEquals, TokenContains:
		return BinaryOp(p.tokenizer.Text()), nil
	default:
		return "", fmt.Errorf("unrecognized binary operator: %q (%v)", p.tokenizer.Text(), token)
	}
}

func isQuotedString(text string) bool {
	if len(text) < 2 {
		return false
	}

	return text[0] == '"' && text[len(text)-1] == '"'
}

func (p *Parser) parseConstValue() (*Const, error) {
	token, more := p.tokenizer.Next()
	if !more {
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

func (p *Parser) parseSortOrder() bool {
	maybeSortOrder, sortOrderText, more := p.tokenizer.Peek()
	if more && maybeSortOrder == TokenIdent {
		switch sortOrderText {
		case "asc":
			return false
		case "desc":
			return true
		}
	}

	// If it isn't a sort keyword, then take the default behavior of ascending.
	return false
}
