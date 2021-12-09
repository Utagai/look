package execution

import (
	"fmt"
	"math"
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
	// This should always exist, because if it did not, parsing would have failed
	// earlier.
	funcValidator, _ := breeze.LookupFuncValidator(function.Name)

	// Evaluate the arguments.
	evaluatedArgs := make([]*breeze.Const, len(function.Args))
	for i := range function.Args {
		evaluatedArg, err := evaluateExprToConst(function.Args[i], datum)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate argument %d (%q): %w", i, function.Args[i].GetStringRepr(), err)
		}
		evaluatedArgs[i] = evaluatedArg
	}

	if err := funcValidator.ValidateTypes(evaluatedArgs); err != nil {
		return &breeze.Const{
			Kind:        breeze.ConstKindString,
			Stringified: err.ToEmbeddedDatumErrorMessage(),
		}, nil
	}

	return executeFunction(function, evaluatedArgs)
}

func executeFunction(function *breeze.Function, args []*breeze.Const) (*breeze.Const, error) {
	switch function.Name {
	case "pow":
		base := args[0].Interface()
		exp := args[1].Interface()

		return &breeze.Const{
			Kind:        breeze.ConstKindNumber,
			Stringified: fmt.Sprintf("%f", math.Pow(base.(float64), exp.(float64))),
		}, nil
	case "hello":
		return &breeze.Const{
			Kind:        breeze.ConstKindString,
			Stringified: "hello world!",
		}, nil
	}
	return nil, fmt.Errorf("unrecognized function: %q", function.Name)
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
	leftConst, err := evaluateExpr(expr.Left, datum)
	if err != nil {
		return nil, err
	}
	rightConst, err := evaluateExpr(expr.Right, datum)
	if err != nil {
		return nil, err
	}

	// TODO: For now, we are only supporting numbers and these 4 ops, but in the
	// future, some of these ops will work across different types and there are
	// more ops as well.
	leftNum, ok := leftConst.(float64)
	if !ok {
		// TODO: These should be type mismatch error strings.
		return &breeze.Const{
			Kind:        breeze.ConstKindNull,
			Stringified: "null",
		}, nil
	}
	rightNum, ok := rightConst.(float64)
	if !ok {
		return &breeze.Const{
			Kind:        breeze.ConstKindNull,
			Stringified: "null",
		}, nil
	}

	result := 0.0
	switch expr.Op {
	case breeze.BinaryOpPlus:
		result = leftNum + rightNum
	case breeze.BinaryOpMinus:
		result = leftNum - rightNum
	case breeze.BinaryOpMultiply:
		result = leftNum * rightNum
	case breeze.BinaryOpDivide:
		result = leftNum / rightNum
	default:
		panic(fmt.Sprintf("unrecognized operator: %q", expr.Op))
	}

	return &breeze.Const{
		Kind:        breeze.ConstKindNumber,
		Stringified: fmt.Sprintf("%f", result),
	}, nil
}
