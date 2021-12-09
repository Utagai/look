package execution

import (
	"fmt"
	"strconv"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func evaluateExpr(expr breeze.Expr, datum datum.Datum) (interface{}, error) {
	breezeConst, err := evaluateExprToConst(expr, datum)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate to constant: %w", err)
	}

	switch breezeConst.Kind {
	case breeze.ConstKindBool:
		return breezeConst.Stringified == "true", nil
	case breeze.ConstKindNull:
		return nil, nil
	case breeze.ConstKindString:
		return breezeConst.Stringified, nil
	case breeze.ConstKindNumber:
		num, err := strconv.ParseFloat(breezeConst.Stringified, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q as a number: %w", breezeConst.Stringified, err)
		}
		return num, nil
	default:
		panic(fmt.Errorf("unrecognized const kind: %q", breezeConst.Kind))
	}
}

func evaluateExprToConst(expr breeze.Expr, datum datum.Datum) (*breeze.Const, error) {
	switch expr.ExprKind() {
	case breeze.ExprKindTerm:
		// The simple case where this is already just a single value.
		return evaluateValue(expr.(breeze.Value), datum)
	case breeze.ExprKindBinary:
		// Otherwise, we must evaluate this expr:
		binaryExpr := expr.(*breeze.BinaryExpr)
		return evaluateBinaryExpr(binaryExpr, datum)
	}
	panic(fmt.Sprintf("unrecognized expr kind: %q", expr.ExprKind()))
}

func evaluateValue(value breeze.Value, datum datum.Datum) (*breeze.Const, error) {
	// TODO: Later, this can return things like function or field reference. In
	// these cases, we can't just cast into a native go type.
	switch value.ValueKind() {
	case breeze.ValueKindConst:
		constValue := value.(*breeze.Const)
		return constValue, nil
	case breeze.ValueKindFieldRef:
		return evaluateFieldRef(value.(*breeze.FieldRef), datum), nil
	case breeze.ValueKindFunc:
		return evaluateFunction(value.(*breeze.Function), datum)
	default:
		panic(fmt.Sprintf("unrecognized value kind: %v", value.ValueKind()))
	}
}

func evaluateFieldRef(fieldRef *breeze.FieldRef, datum datum.Datum) *breeze.Const {
	return goValueToConst(datum[fieldRef.Field])
}

func evaluateFunction(function *breeze.Function, datum datum.Datum) (*breeze.Const, error) {
	// Evaluate the arguments.
	evaluatedArgs := make([]*breeze.Const, len(function.Args))
	for i := range function.Args {
		evaluatedArg, err := evaluateExprToConst(function.Args[i], datum)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate argument %d (%q): %w", i, function.Args[i].GetStringRepr(), err)
		}
		evaluatedArgs[i] = evaluatedArg
	}

	return executeFunction(function, evaluatedArgs)
}

func goValueToConst(val interface{}) *breeze.Const {
	if val == nil {
		return &breeze.Const{
			Kind:        breeze.ConstKindNull,
			Stringified: "null",
		}
	}

	switch tval := val.(type) {
	case uint, uint8, uint16, uint32, uint64:
	case int, int8, int16, int32, int64:
	case float32, float64:
		return &breeze.Const{
			Kind:        breeze.ConstKindNumber,
			Stringified: fmt.Sprintf("%v", tval),
		}
	case string:
		return &breeze.Const{
			Kind:        breeze.ConstKindString,
			Stringified: tval,
		}
	case bool:
		return &breeze.Const{
			Kind:        breeze.ConstKindBool,
			Stringified: fmt.Sprintf("%t", tval),
		}
	}

	panic(fmt.Sprintf("unrecognized base type: %T", val))
}

func evaluateBinaryExpr(expr *breeze.BinaryExpr, datum datum.Datum) (*breeze.Const, error) {
	leftConst, err := evaluateExprToConst(expr.Left, datum)
	if err != nil {
		return nil, err
	}
	rightConst, err := evaluateExprToConst(expr.Right, datum)
	if err != nil {
		return nil, err
	}

	return evaluateOp(*leftConst, *rightConst, expr.Op, datum)
}
