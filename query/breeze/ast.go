package breeze

import (
	"fmt"
	"strconv"
)

type ConstKind string

const (
	ConstKindString ConstKind = "string"
	ConstKindNumber ConstKind = "number"
	ConstKindBool   ConstKind = "bool"
)

type Const struct {
	Kind        ConstKind
	Stringified string
}

func (c *Const) String() string {
	return c.Stringified
}

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

type BinaryOp string

const (
	BinaryOpEquals   BinaryOp = "="
	BinaryOpGeq               = ">"
	BinaryOpContains          = "contains"
)

type Check struct {
	Field string
	Value *Const
	Op    BinaryOp
}

func (c *Check) String() string {
	return fmt.Sprintf("%s %s %v", c.Field, c.Op, c.Value)
}

type Find struct {
	Checks []*Check
}

func (f *Find) Name() string {
	return "find"
}

type Sort struct {
	Descending bool
	Field      string
}

func (s *Sort) Name() string {
	return "sort"
}

type Stage interface {
	Name() string
}
