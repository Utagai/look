package breeze

import (
	"fmt"
	"strconv"
	"strings"
)

// BinaryOp enumerates the kinds of binary operations in breeze.
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
	// BinaryOpEquals is the equality operation.
	BinaryOpEquals BinaryOp = "="
	// BinaryOpGeq is the greater-than-or-equal operation.
	BinaryOpGeq BinaryOp = ">"
	// BinaryOpContains is the contains operation.
	BinaryOpContains BinaryOp = "contains"
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

// ConcreteKind enumerates the kinds of concrete values in Breeze.
type ConcreteKind string

const (
	// ConcreteKindScalar refers to scalar values.
	ConcreteKindScalar = "scalar"
	// ConcreteKindArray refers to arrays.
	ConcreteKindArray = "array"
	// ConcreteKindDoc = "doc"
)

// Concrete represents a breeze expression & value that can be evaluated to a
// single, 'concrete' Go value for query processing. In other words, this is an
// AST node that can be realized as a true Go value for use in actual
// computation.
type Concrete interface {
	Value
	ConcreteKind() ConcreteKind
	Interface() (interface{}, error)
}

// ScalarKind enumerates the kinds of scalars in breeze.
type ScalarKind string

const (
	// ScalarKindString represents a string constant.
	ScalarKindString ScalarKind = "string"
	// ScalarKindNumber represents a number constant.
	ScalarKindNumber ScalarKind = "number"
	// ScalarKindBool represents a bool constant.
	ScalarKindBool ScalarKind = "bool"
	// ScalarKindNull represents the special null constant value.
	ScalarKindNull ScalarKind = "null"
)

// Scalar is a single-dimensional constant value.
type Scalar struct {
	Kind        ScalarKind
	Stringified string
}

// ValueKind implements the Value interface.
func (s *Scalar) ValueKind() ValueKind {
	return ValueKindScalar
}

// GetStringRepr implements the Expr interface.
func (s *Scalar) GetStringRepr() string {
	return s.Stringified
}

// ExprKind implements the Expr interface.
func (s *Scalar) ExprKind() ExprKind {
	return ExprKindTerm
}

// ConcreteKind implements the Concrete interface.
func (s *Scalar) ConcreteKind() ConcreteKind {
	return ConcreteKindScalar
}

// Interface implements the Concrete interface.
func (s *Scalar) Interface() (interface{}, error) {
	switch s.Kind {
	case ScalarKindString:
		return s.Stringified, nil
	case ScalarKindNumber:
		f64, _ := strconv.ParseFloat(s.Stringified, 64)
		return f64, nil
	case ScalarKindBool:
		return s.Stringified == "true", nil
	case ScalarKindNull:
		return nil, nil
	default:
		panic(fmt.Sprintf("unexpected const kind: %q", s.Kind))
	}
}

// Array is a breeze array.
type Array []Expr

// ExprKind implements the Expr interface.
func (a Array) ExprKind() ExprKind {
	return ExprKindTerm
}

// ValueKind implements the Value interface.
func (a Array) ValueKind() ValueKind {
	return ValueKindArray
}

// GetStringRepr implements the Value interface.
func (a Array) GetStringRepr() string {
	memberStrs := make([]string, len(a))
	for i := range a {
		memberStrs[i] = a[i].GetStringRepr()
	}

	return fmt.Sprintf("[%s]", strings.Join(memberStrs, ","))
}

// ConcreteKind implements the Concrete interface.
func (a Array) ConcreteKind() ConcreteKind {
	return ConcreteKindArray
}

// Interface implements the Concrete interface.
func (a Array) Interface() (interface{}, error) {
	goArr := make([]interface{}, len(a))
	for i := range a {
		c, ok := a[i].(Concrete)
		if !ok {
			return nil, fmt.Errorf("element %d (%q) is not a const but must be for Array to be const", i, a[i].GetStringRepr())
		}

		var err error
		goArr[i], err = c.Interface()
		if err != nil {
			return nil, fmt.Errorf("failed to realize element %d (%q): %w", i, a[i].GetStringRepr(), err)
		}
	}

	return goArr, nil
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
type ValueKind string

const (
	// ValueKindScalar represents a breeze Const.
	ValueKindScalar = "scalar"
	// ValueKindArray represents a breeze Array.
	ValueKindArray = "array"
	// ValueKindFieldRef represents a reference to a field.
	ValueKindFieldRef = "fieldref"
	// ValueKindFunc represents a evaluatable function.
	ValueKindFunc = "func"
)

// Value is simply a value in breeze. It could be a constant, field reference, or
// function.
type Value interface {
	Expr
	ValueKind() ValueKind
}

// Filter is a stage that applies a series of filters to the data.
type Filter struct {
	Exprs []Expr
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
	Field      string
	Assignment Expr
}
