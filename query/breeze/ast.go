package breeze

import (
	"fmt"
	"strconv"
	"strings"
)

// BinaryOp enumerates the kinds of binary operations in breeze.
// TODO: This is confusing when used alongside BinaryCmpOp.
type BinaryOp string

const (
	// BinaryOpPlus is the equality operation.
	BinaryOpPlus BinaryOp = "+"
	// BinaryOpMinus is the equality operation.
	BinaryOpMinus BinaryOp = "-"
	// BinaryOpMultiply is the equality operation.
	BinaryOpMultiply BinaryOp = "*"
	// BinaryOpDivide is the equality operation.
	BinaryOpDivide BinaryOp = "/"
)

// ExprKind denotes the kind of expression.
type ExprKind string

const (
	// ExprKindTerm is for the terms that make up an expression, joined together
	// by operators.
	ExprKindTerm = "TERM"
	// ExprKindBinary is for binary expressions that tie together expressions with
	// an operator.
	ExprKindBinary = "BINARY"
)

// Expr is a breeze expression.
type Expr interface {
	ExprKind() ExprKind
	GetStringRepr() string
}

// BinaryExpr is a binary expression that combines two expressions together by
// an operator.
type BinaryExpr struct {
	Left  Expr
	Right Expr
	Op    BinaryOp
}

// ExprKind implements the Expr interface.
func (b *BinaryExpr) ExprKind() ExprKind {
	return ExprKindBinary
}

// GetStringRepr implements the Expr interface.
func (b *BinaryExpr) GetStringRepr() string {
	return fmt.Sprintf("%s %s %s", b.Left.GetStringRepr(), b.Op, b.Right.GetStringRepr())
}

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

// ValueKind implements the Value interface.
func (c *Const) ValueKind() ValueKind {
	return ValueKindConst
}

// GetStringRepr implements the Value interface.
func (c *Const) GetStringRepr() string {
	return c.Stringified
}

// ExprKind implements the Expr interface.
func (c *Const) ExprKind() ExprKind {
	return ExprKindTerm
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

// FieldRef is a reference to a field of a datum.
type FieldRef struct {
	Field string
}

// ValueKind implements the Value interface.
func (f *FieldRef) ValueKind() ValueKind {
	return ValueKindFieldRef
}

// GetStringRepr implements the Value interface.
func (f *FieldRef) GetStringRepr() string {
	return f.Field
}

// ExprKind implements the Expr interface.
func (f *FieldRef) ExprKind() ExprKind {
	return ExprKindTerm
}

// Function is a breeze function.
type Function struct {
	Name string
	Args []Expr
}

// ValueKind implements the Value interface.
func (f *Function) ValueKind() ValueKind {
	return ValueKindFunc
}

// GetStringRepr implements the Value interface.
// TODO: Test this?
func (f *Function) GetStringRepr() string {
	argsStrSlice := make([]string, len(f.Args))
	for i := range f.Args {
		argsStrSlice[i] = f.Args[i].GetStringRepr()
	}

	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(argsStrSlice, ","))
}

// ExprKind implements the Expr interface.
func (f *Function) ExprKind() ExprKind {
	return ExprKindTerm
}

// ValueKind enumerates the kinds of values in breeze expressions.
// TODO: Move this and Value to the top, above const and fieldref.
type ValueKind string

const (
	// ValueKindConst represents a breeze Const.
	ValueKindConst = "const"
	// ValueKindFieldRef represents a reference to a field.
	ValueKindFieldRef = "fieldref"
	// ValueKindFunc represents a evaluatable function.
	ValueKindFunc = "func"
)

// Value is simply a value in breeze. It could be a constant, field reference, or
// function.
type Value interface {
	Expr
	// TODO: Rename this to ValueKind() for consistency with Expr.
	ValueKind() ValueKind
	GetStringRepr() string
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

// BinaryCmpOp enumerates the kinds of binary comparison operations in breeze.
type BinaryCmpOp string

const (
	// BinaryOpEquals is the equality operation.
	BinaryOpEquals BinaryCmpOp = "="
	// BinaryOpGeq is the greater-than-or-equal operation.
	BinaryOpGeq = ">"
	// BinaryOpContains is the contains operation.
	BinaryOpContains = "contains"
)

// BinaryCheck is a filter condition for filter that uses a binary operation.
type BinaryCheck struct {
	Field string
	Value *Const
	Op    BinaryCmpOp
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

// Group is a stage that performs grouping of data and aggregates computations
// over them.
// TODO: Rename to reduce, or is reduce the wrong word?
type Group struct {
	AggFunc        AggregateFunc
	GroupByField   *string
	AggregateField string
}

// Name implements the Stage interface.
func (g *Group) Name() string {
	return "group"
}

// AggregateFunc is an aggregate function.
type AggregateFunc string

// The various kinds of aggregate functions in breeze.
const (
	AggFuncSum    AggregateFunc = "sum"
	AggFuncAvg    AggregateFunc = "avg"
	AggFuncCount  AggregateFunc = "count"
	AggFuncMin    AggregateFunc = "min"
	AggFuncMax    AggregateFunc = "max"
	AggFuncMode   AggregateFunc = "mode"
	AggFuncStdDev AggregateFunc = "stddev"
)

// Map is a stage that performs transformations on a per-field basis.
type Map struct {
	Assignments []FieldAssignment
}

// Name implements the Stage interface.
func (m *Map) Name() string {
	return "map"
}

// FieldAssignment is a remapping of a field, found in maps.
type FieldAssignment struct {
	Field string
	// TODO: We should possibly be able to handle non-binary exprs.
	Assignment Expr
}
