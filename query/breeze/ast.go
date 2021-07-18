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
	default:
		panic(fmt.Sprintf("unexpected const kind: %q", c.Kind))
	}
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

// Check is a filter condition for find.
type Check struct {
	Field string
	Value *Const
	Op    BinaryOp
}

func (c *Check) String() string {
	return fmt.Sprintf("%s %s %v", c.Field, c.Op, c.Value)
}

// Filter is a stage that applies a series of filters to the data.
type Filter struct {
	Checks []*Check
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
