package liquid

import (
	"fmt"
	"strconv"

	"github.com/utagai/look/datum"
)

type ConstKind int

const (
	ConstKindString ConstKind = iota
	ConstKindNumber
)

type Const struct {
	Kind        ConstKind
	Stringified string
}

func (c *Const) String() string {
	return c.Stringified
}

type BinaryOp string

const (
	BinaryOpEquals BinaryOp = "="
	BinaryOpGeq             = ">"
)

func evaluateBinaryOp(left interface{}, right *Const, op BinaryOp) bool {
	rightString := right.Stringified
	rightFloat, _ := strconv.ParseFloat(right.Stringified, 64)
	switch op {
	case BinaryOpEquals:
		switch right.Kind {
		case ConstKindString:
			return left.(string) == rightString
		case ConstKindNumber:
			return left.(float64) == rightFloat
		default:
			panic(fmt.Sprintf("unrecognized ConstKind: %q", right.Kind))
		}
	case BinaryOpGeq:
		switch right.Kind {
		case ConstKindString:
			return left.(string) > rightString
		case ConstKindNumber:
			return left.(float64) > rightFloat
		default:
			panic(fmt.Sprintf("unrecognized ConstKind: %q", right.Kind))
		}
	default:
		panic(fmt.Sprintf("unrecognized BinaryOp: %q", op))
	}
}

type Check struct {
	Field string
	Value *Const
	Op    BinaryOp
}

func (c *Check) Evaluate(datum datum.Datum) bool {
	fieldValue, ok := datum[c.Field]
	if !ok {
		// If the field does not exist on this datum, evaluate to false.
		return false
	}

	return evaluateBinaryOp(fieldValue, c.Value, c.Op)
}

func (c *Check) String() string {
	return fmt.Sprintf("%s %s %v", c.Field, c.Op, c.Value)
}

type Find struct {
	Checks []*Check
}

type FindStream struct {
	*Find
	source datum.DatumStream
}

// Next implements the datum.DatumStream interface.
func (fs *FindStream) Next() (datum.Datum, error) {
	// Keep iterating the stream until something passes the checks.
outer:
	for {
		datum, err := fs.source.Next()
		if err != nil {
			return nil, err
		}

		for _, check := range fs.Find.Checks {
			// If we failed, move onto the next datum.
			if !check.Evaluate(datum) {
				continue outer
			}
		}

		// If we get here, we have successfully evaluated against every check, and
		// we can be returned.
		return datum, nil
	}
}

// Execute implements the Stage interface.
func (f *Find) Execute(datums datum.DatumStream) (datum.DatumStream, error) {
	return &FindStream{
		Find:   f,
		source: datums,
	}, nil
}

type SortStream struct {
	*Sort
	source datum.DatumStream
}

func (ss *SortStream) Next() (datum.Datum, error) {
	// TODO: Sort this...
	return ss.source.Next()
}

type Sort struct {
	Descending bool
	Field      string
}

func (s *Sort) Execute(datums datum.DatumStream) (datum.DatumStream, error) {
	return &SortStream{
		Sort:   s,
		source: datums,
	}, nil
}

type Stage interface {
	// TODO: This should be decoupled from the AST.
	Execute(datum.DatumStream) (datum.DatumStream, error)
}
