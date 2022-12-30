package execution

import (
	"fmt"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func evaluateExpr(expr breeze.Expr, datum datum.Datum) (interface{}, error) {
	breezeConcrete, err := evaluateExprToConcrete(expr, datum)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate to constant: %w", err)
	}

	return breezeConcrete.Interface()
}

func evaluateExprToConcrete(expr breeze.Expr, datum datum.Datum) (breeze.Concrete, error) {
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

func evaluateValue(value breeze.Value, datum datum.Datum) (breeze.Concrete, error) {
	switch value.ValueKind() {
	case breeze.ValueKindScalar:
		scalarValue := value.(*breeze.Scalar)
		return scalarValue, nil
	case breeze.ValueKindFieldRef:
		return evaluateFieldRef(value.(*breeze.FieldRef), datum), nil
	case breeze.ValueKindFunc:
		return evaluateFunction(value.(*breeze.Function), datum)
	case breeze.ValueKindArray:
		arrValue := value.(breeze.Array)
		concreteArr := make([]breeze.Expr, len(arrValue))
		for i := range arrValue {
			var err error
			concreteArr[i], err = evaluateExprToConcrete(arrValue[i], datum)
			if err != nil {
				return nil, err
			}
		}
		return breeze.Array(concreteArr), nil
	default:
		panic(fmt.Sprintf("unrecognized value kind: %v", value.ValueKind()))
	}
}

func evaluateFieldRef(fieldRef *breeze.FieldRef, datum datum.Datum) breeze.Concrete {
	val, ok := datum[fieldRef.Field]
	if !ok {
		return &breeze.Missing{}
	}
	return goValueToConcrete(val)
}

func evaluateFunction(function *breeze.Function, datum datum.Datum) (breeze.Concrete, error) {
	// Evaluate the arguments.
	evaluatedArgs := make([]breeze.Concrete, len(function.Args))
	for i := range function.Args {
		evaluatedArg, err := evaluateExprToConcrete(function.Args[i], datum)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate argument %d (%q): %w", i, function.Args[i].GetStringRepr(), err)
		}
		evaluatedArgs[i] = evaluatedArg
	}

	return executeFunction(function, evaluatedArgs)
}

func goValueToConcrete(val interface{}) breeze.Concrete {
	if val == nil {
		return &breeze.Scalar{
			Kind:        breeze.ScalarKindNull,
			Stringified: "null",
		}
	}

	switch tval := val.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return &breeze.Scalar{
			Kind:        breeze.ScalarKindNumber,
			Stringified: fmt.Sprintf("%v", tval),
		}
	case int, int8, int16, int32, int64:
		return &breeze.Scalar{
			Kind:        breeze.ScalarKindNumber,
			Stringified: fmt.Sprintf("%v", tval),
		}
	case float32, float64:
		return &breeze.Scalar{
			Kind:        breeze.ScalarKindNumber,
			Stringified: fmt.Sprintf("%v", tval),
		}
	case string:
		return &breeze.Scalar{
			Kind:        breeze.ScalarKindString,
			Stringified: tval,
		}
	case bool:
		return &breeze.Scalar{
			Kind:        breeze.ScalarKindBool,
			Stringified: fmt.Sprintf("%t", tval),
		}
	}

	panic(fmt.Sprintf("unrecognized base type: %T", val))
}

func evaluateBinaryExpr(expr *breeze.BinaryExpr, datum datum.Datum) (breeze.Concrete, error) {
	leftConst, err := evaluateExprToConcrete(expr.Left, datum)
	if err != nil {
		return &breeze.Scalar{}, err
	}
	rightConst, err := evaluateExprToConcrete(expr.Right, datum)
	if err != nil {
		return &breeze.Scalar{}, err
	}

	return evaluateOp(leftConst, rightConst, expr.Op, datum)
}
