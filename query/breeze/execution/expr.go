package execution

import (
	"fmt"
	"math"
	"strconv"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func evaluateValueOrExpr(valOrExpr breeze.ValueOrExpr, datum datum.Datum) (interface{}, error) {
	breezeConst, err := evaluateValueOrExprToConst(valOrExpr, datum)
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
		panic(fmt.Errorf("unrecognized kind: %q", breezeConst.Kind))
	}
}

func evaluateValueOrExprToConst(valOrExpr breeze.ValueOrExpr, datum datum.Datum) (*breeze.Const, error) {
	if valOrExpr.Value != nil {
		// The simple case where this is already just a single value.
		return evaluateValue(valOrExpr.Value, datum)
	}

	// Otherwise, we must evaluate this expr:
	return evaluateExpr(valOrExpr.Expr, datum)
}

func evaluateValue(value breeze.Value, datum datum.Datum) (*breeze.Const, error) {
	// TODO: Later, this can return things like function or field reference. In
	// these cases, we can't just cast into a native go type.
	switch value.GetKind() {
	case breeze.ValueKindConst:
		constValue := value.(*breeze.Const)
		return constValue, nil
	case breeze.ValueKindFieldRef:
		return evaluateFieldRef(value.(*breeze.FieldRef), datum), nil
	case breeze.ValueKindFunc:
		return evaluateFunction(value.(*breeze.Function), datum)
	default:
		panic(fmt.Sprintf("unrecognized value kind: %v", value.GetKind()))
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
		evaluatedArg, err := evaluateValue(function.Args[i], datum)
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

func evaluateExpr(expr breeze.Expr, datum datum.Datum) (*breeze.Const, error) {
	panic("TODO")
}
