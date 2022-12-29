package breeze

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

// Parser parses breeze queries.
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
	exprs := []Expr{}
	for {
		token, _ := p.tokenizer.Peek()
		if token == TokenStageSeparator || token == TokenEOF {
			break // No more checks to parse.
		} else if token == TokenComma {
			p.tokenizer.Next()
			continue
		}

		expr, err := p.parseExpr(p.tokenizer.Next())
		if err == io.EOF {
			break // No more checks to parse.
		} else if err != nil {
			return nil, fmt.Errorf("failed to parse filter expression: %w", err)
		}

		exprs = append(exprs, expr)
	}

	if len(exprs) > 0 {
		log.Printf("Returning exprs: %+v", exprs[0])
	}
	return &Filter{Exprs: exprs}, nil
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
		token, _ := p.tokenizer.Peek()
		if token == TokenStageSeparator || token == TokenEOF {
			break // No more checks to parse.
		}

		assignment, err := p.parseAssignment()
		if err == io.EOF {
			break // No more checks to parse.
		} else if err != nil {
			return nil, fmt.Errorf("failed to parse assignment: %w", err)
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
		return nil, err
	}

	expr, err := p.parseExpr(p.tokenizer.Next())
	if err != nil {
		return nil, err
	}

	return &FieldAssignment{
		Field:      field,
		Assignment: expr,
	}, nil
}

func (p *Parser) parseField() (string, error) {
	token := p.tokenizer.Next()
	if token == TokenEOF {
		return "", errors.New("expected a field, but reached end of query")
	}

	if token == TokenIdent || token == TokenString {
		return p.tokenizer.Text(), nil
	}

	return "", fmt.Errorf("expected a field identifier, but got %q (%s)", p.tokenizer.Text(), token.String())
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

func binaryOpToToken(bOp BinaryOp) Token {
	switch bOp {
	case BinaryOpPlus:
		return TokenPlus
	case BinaryOpMinus:
		return TokenMinus
	case BinaryOpMultiply:
		return TokenMultiply
	case BinaryOpDivide:
		return TokenDivide
	default:
		panic(fmt.Sprintf("unrecognized binary op: %v", bOp))
	}
}

func getBinaryOpPrecedence(bOp BinaryOp) int {
	// A mapping of expression operator tokens linked to their precedence level.
	// Higher precedence values means the token has a higher precedence.
	binaryOpPrecedence := map[Token]int{
		TokenMultiply: 1,
		TokenDivide:   1,
		TokenPlus:     0,
		TokenMinus:    0,
	}

	bOpToken := binaryOpToToken(bOp)

	prec, ok := binaryOpPrecedence[bOpToken]
	if !ok {
		panic(fmt.Sprintf("unrecognized binary op: %v", bOp))
	}

	return prec
}

func (p *Parser) parseExpr(token Token) (Expr, error) {
	var leftExpr Expr
	var err error

	if token == TokenLParen {
		leftExpr, err = p.parseExpr(p.tokenizer.Next())
		if err != nil {
			return nil, fmt.Errorf("failed to parse expression: %w", err)
		}
		if p.tokenizer.Next() != TokenRParen {
			return nil, fmt.Errorf("expected a closing paranthesis, but got %q", p.tokenizer.Text())
		}
	} else {
		// We should always expect _at least_ a single value, aka, a single-term
		// expression. If we don't find this at least, that means the expression
		// doesn't exist in the query even though it should.
		leftExpr, err = p.parseValue(token)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value in expr: %w", err)
		}
	}

	token, _ = p.tokenizer.Peek()
	bOp, err := p.parseBinaryOp(token)
	if err != nil {
		// No more tokens for this expression.
		return leftExpr, nil
	}
	_ = p.tokenizer.Next()

	// Grab the next token, so that we can check if the next suffix of the
	// expression input is parenthesized for precedence determination.
	token = p.tokenizer.Next()
	rightExpr, err := p.parseExpr(token)

	return p.applyPrecedence(leftExpr, bOp, rightExpr, token == TokenLParen), err
}

func (p *Parser) applyPrecedence(leftExpr Expr, leftOp BinaryOp, rightExpr Expr, rightIsParenthesized bool) *BinaryExpr {
	// HACK(?): Fixes the precedence before committing to the binary expr.
	// Classically, this is achieved implicitly through the recursion
	// tree, the textbook example being parseTerm and parseFactor
	// functions. I decided to try something different to see if the
	// code gets more readable, since the classical approach involves
	// building the shape of the AST _implicitly_, which seems a bad
	// decision for a trait so important (and it took me personally a
	// while to see how precedence is encoded into the call stack). In
	// other words, I think it is more readable to explicitly handle
	// precedence in the parser logic rather than have it done
	// implicitly through construction.  All that said, this if
	// statement is not easy to look at and moving around trees aren't
	// easy to visualize either... so maybe the readability improvement
	// is minor. Furthermore, maybe my thinking is buggy here! I'll let
	// this sit for now and revisit it later (or, if this is buggy, I'll
	// probably just give up and do it 'the right way').
	if rightBinaryExpr, ok := rightExpr.(*BinaryExpr); ok &&
		!rightIsParenthesized &&
		getBinaryOpPrecedence(leftOp) > getBinaryOpPrecedence(rightBinaryExpr.Op) {
		// We have the higher precedence operator, call it *. We have a value X,
		// and a right-side expression with lower precedence operator, +: Y + Z.
		// Without this, we would do X * Y + Z -> X * (Y + Z).
		// With this, we do X * Y + Z -> (X * Y) + Z.
		leftExpr = &BinaryExpr{
			Left:  leftExpr,
			Op:    leftOp,
			Right: rightBinaryExpr.Left,
		}
		leftOp = rightBinaryExpr.Op
		rightExpr = rightBinaryExpr.Right
	}

	return &BinaryExpr{
		Left:  leftExpr,
		Op:    leftOp,
		Right: rightExpr,
	}
}

func (p *Parser) parseValue(token Token) (Value, error) {
	if token == TokenEOF {
		return nil, errors.New("expected a value, but reached end of query")
	}

	// The following attempts at parsing the possibilities is actually _incorrect_.
	// If any of these fail midway through their parse, then the tokenizer will
	// have advanced past where token in this context is, meaning the later parse
	// functions will fail.
	// That said, this is basically the whole ambiguous grammar issue. Our parser
	// is recursive descent, so it is limited to LL(k) grammars for it to run
	// without backtracking. Right now, I think all these parse functions
	// represent _unambiguous_ non-terminals, so there's no actual issue here
	// (they'll fail with k = 1 tokens), but if we do encounter an issue, we'll
	// have to either bump up k and pass in multiple tokens, or implement some
	// rewind feature to the tokenizer.
	constValue, constErr := p.parseConstValue(token)
	if constErr == nil {
		return constValue, nil
	}

	fieldRefValue, fieldRefErr := p.parseFieldRef(token)
	if fieldRefErr == nil {
		return fieldRefValue, nil
	}

	functionValue, funcErr := p.parseFunction(token)
	if funcErr == nil {
		return functionValue, nil
	}

	arrayValue, arrayErr := p.parseArray(token)
	if arrayErr == nil {
		return arrayValue, nil
	}

	return nil, fmt.Errorf(
		"failed to parse a value; expected a constant value (%s), field reference (%s), function (%s), or array (%s)",
		constErr,
		fieldRefErr,
		funcErr,
		arrayErr,
	)
}

func (p *Parser) parseBinaryOp(token Token) (BinaryOp, error) {
	switch token {
	case TokenPlus:
		return BinaryOpPlus, nil
	case TokenMinus:
		return BinaryOpMinus, nil
	case TokenMultiply:
		return BinaryOpMultiply, nil
	case TokenDivide:
		return BinaryOpDivide, nil
	case TokenEquals:
		return BinaryOpEquals, nil
	case TokenGEQ:
		return BinaryOpGeq, nil
	case TokenContains:
		return BinaryOpContains, nil
	}

	return "", fmt.Errorf("unrecognized binary operator: %q (%v)", p.tokenizer.Text(), token)
}

func (p *Parser) parseConstValue(token Token) (*Scalar, error) {
	if token == TokenEOF {
		return nil, errors.New("expected a constant value, but reached end of query")
	}

	switch token {
	case TokenFloat, TokenInt:
		return &Scalar{
			Kind:        ScalarKindNumber,
			Stringified: p.tokenizer.Text(),
		}, nil
	case TokenChar, TokenString:
		text := p.tokenizer.Text()
		if !isQuotedString(text) {
			return nil, fmt.Errorf("expected a properly quoted string, but got %q", text)
		}
		textWithoutQuotes := text[1 : len(text)-1]
		return &Scalar{
			Kind:        ScalarKindString,
			Stringified: textWithoutQuotes,
		}, nil
	case TokenFalse, TokenTrue:
		return &Scalar{
			Kind:        ScalarKindBool,
			Stringified: p.tokenizer.Text(),
		}, nil
	case TokenNull:
		return &Scalar{
			Kind:        ScalarKindNull,
			Stringified: p.tokenizer.Text(),
		}, nil
	default:
		return nil, fmt.Errorf("expected a constant value, but got: %q", p.tokenizer.Text())
	}
}

func (p *Parser) parseArray(token Token) (Array, error) {
	if token != TokenLSqBracket {
		return nil, fmt.Errorf("expected array to start with '[', but found %q", p.tokenizer.Text())
	}
	token = p.tokenizer.Next() // Advance past the bracket we just parsed.

	if token == TokenRSqBracket { // Empty array.
		p.tokenizer.Next() // Advance past the closing ']'.
		return []Expr{}, nil
	}

	exprs := []Expr{}
	for ; token != TokenRSqBracket && token != TokenEOF; token = p.tokenizer.Next() {
		expr, err := p.parseExpr(token)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %dth array member: %w", len(exprs)+1, err)
		}
		exprs = append(exprs, expr)
		token = p.tokenizer.Next() // Advance past the expr we just parsed.

		if token == TokenRSqBracket {
			break
		}

		if token != TokenComma {
			return nil, fmt.Errorf("array items should be delimited by commas, but found %q", p.tokenizer.Text())
		}
	}

	return exprs, nil
}

func (p *Parser) parseFieldRef(token Token) (*FieldRef, error) {
	fieldRefText := p.tokenizer.Text()
	if token != TokenIdent {
		return nil, fmt.Errorf("expected an identifier, but got %q", fieldRefText)
	}

	if !strings.HasPrefix(fieldRefText, ".") {
		return nil, fmt.Errorf("field references must start with '.'")
	}

	if len(fieldRefText) == 1 { // It is just the '.'
		return nil, fmt.Errorf("missing field name")
	}

	return &FieldRef{
		Field: strings.TrimPrefix(fieldRefText, "."),
	}, nil
}

func (p *Parser) parseFunction(token Token) (*Function, error) {
	funcName := p.tokenizer.Text()
	if token != TokenIdent {
		return nil, fmt.Errorf("expected an identifier, but got %q", funcName)
	}

	funcValidator, found := LookupFuncValidator(funcName)
	if !found {
		return nil, fmt.Errorf("unrecognized function: %q", funcName)
	}

	if p.tokenizer.Next() != TokenLParen {
		return nil, fmt.Errorf("expected an opening parenthesis for a function start, but got %q", p.tokenizer.Text())
	}

	// Loop for collecting function arguments. Keeps going until it hits RParen.
	// Expects a comma between each argument.
	args := []Expr{}
	for {
		expr, err := p.parseExpr(p.tokenizer.Next())
		if err != nil {
			break
		}

		args = append(args, expr)

		nextToken := p.tokenizer.Next()
		if nextToken == TokenRParen {
			// Function is done, stop.
			break
		} else if nextToken == TokenComma {
			// Possibly more coming, continue.
			continue
		} else {
			return nil, fmt.Errorf("expected either a closing parenthesis or comma, but got %q", p.tokenizer.Text())
		}
	}

	if funcValidator.ExpectedNumArgs() != len(args) {
		return nil, fmt.Errorf("expected %d args, got %d", funcValidator.ExpectedNumArgs(), len(args))
	}

	return &Function{
		Name: funcName,
		Args: args,
	}, nil
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
