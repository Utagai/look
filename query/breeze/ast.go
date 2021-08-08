package breeze

import (
	"fmt"
	"strconv"
)

// ConstKind enumerates the kinds of constants in breeze.
type ConstKind string

const (
	// ConstKindString represents a string constant.
	ConstKindString ConstKind = "string"
	// ConstKindNumber represents a number constant.
	ConstKindNumber ConstKind = "number"
	// ConstKindBool represents a bool constant.
	ConstKindBool ConstKind = "bool"
	// ConstKindNull represents the special null constant value.
	ConstKindNull ConstKind = "null"
)

// Const is a constant.
type Const struct {
	Kind        ConstKind
	Stringified string
}

func (c *Const) String() string {
	return c.Stringified
}

// Interface returns this constant, casted from its stringified version into the
// proper Go type.
func (c *Const) Interface() interface{} {
	switch c.Kind {
	case ConstKindString:
		return c.Stringified
	case ConstKindNumber:
		f64, _ := strconv.ParseFloat(c.Stringified, 64)
		return f64
	case ConstKindBool:
		return c.Stringified == "true"
	case ConstKindNull:
		return nil
	default:
		panic(fmt.Sprintf("unexpected const kind: %q", c.Kind))
	}
}

// UnaryOp enumerates the kinds of unary operations in breeze.
type UnaryOp string

const (
	// UnaryOpExists is the field existence operator.
	UnaryOpExists UnaryOp = "exists"
	// UnaryOpExistsNot is the field non-existence operator.
	UnaryOpExistsNot = "!exists"
)

// UnaryCheck is a filter condition for a filter that uses a unary operation.
type UnaryCheck struct {
	Field string
	Op    UnaryOp
}

func (c *UnaryCheck) String() string {
	return fmt.Sprintf("%s %s", c.Field, c.Op)
}

// BinaryOp enumerates the kinds of binary operations in breeze.
type BinaryOp string

const (
	// BinaryOpEquals is the equality operation.
	BinaryOpEquals BinaryOp = "="
	// BinaryOpGeq is the greater-than-or-equal operation.
	BinaryOpGeq = ">"
	// BinaryOpContains is the contains operation.
	BinaryOpContains = "contains"
)

// BinaryCheck is a filter condition for filter that uses a binary operation.
type BinaryCheck struct {
	Field string
	Value *Const
	Op    BinaryOp
}

func (c *BinaryCheck) String() string {
	return fmt.Sprintf("%s %s %v", c.Field, c.Op, c.Value)
}

// Filter is a stage that applies a series of filters to the data.
type Filter struct {
	UnaryChecks  []*UnaryCheck
	BinaryChecks []*BinaryCheck
}

// Name implements the Stage interface.
func (f *Filter) Name() string {
	return "filter"
}

// Sort is a stage that sorts the data.
type Sort struct {
	Descending bool
	Field      string
}

// Name implements the Stage interface.
func (s *Sort) Name() string {
	return "sort"
}

// Stage is a single unit of processing to apply to the data.
type Stage interface {
	Name() string
}
