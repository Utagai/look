package execution

import (
	"fmt"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func evaluateOp(left, right breeze.Const, op breeze.BinaryOp, datum datum.Datum) (breeze.Const, error) {
	switch op {
	case breeze.BinaryOpPlus:
		return evaluatePlus(left, right, op, datum)
	case breeze.BinaryOpMinus:
		return evaluateMinus(left, right, op, datum)
	case breeze.BinaryOpMultiply:
		return evaluateMultiply(left, right, op, datum)
	case breeze.BinaryOpDivide:
		return evaluateDivide(left, right, op, datum)
	default:
		panic(fmt.Sprintf("unrecognized operator: %q", op))
	}
}

func evaluatePlus(left, right breeze.Const, op breeze.BinaryOp, datum datum.Datum) (breeze.Const, error) {
	if err := checkTypes(left, breeze.ConstKindNumber, right, breeze.ConstKindNumber); err != nil {
		return err.ToEmbeddedDatumErrorMessage(), nil
	}

	return breeze.Const{
		Kind:        breeze.ConstKindNumber,
		Stringified: fmt.Sprintf("%f", left.Interface().(float64)+right.Interface().(float64)),
	}, nil
}

func evaluateMinus(left, right breeze.Const, op breeze.BinaryOp, datum datum.Datum) (breeze.Const, error) {
	if err := checkTypes(left, breeze.ConstKindNumber, right, breeze.ConstKindNumber); err != nil {
		return err.ToEmbeddedDatumErrorMessage(), nil
	}

	return breeze.Const{
		Kind:        breeze.ConstKindNumber,
		Stringified: fmt.Sprintf("%f", left.Interface().(float64)-right.Interface().(float64)),
	}, nil
}

func evaluateMultiply(left, right breeze.Const, op breeze.BinaryOp, datum datum.Datum) (breeze.Const, error) {
	if err := checkTypes(left, breeze.ConstKindNumber, right, breeze.ConstKindNumber); err != nil {
		return err.ToEmbeddedDatumErrorMessage(), nil
	}

	return breeze.Const{
		Kind:        breeze.ConstKindNumber,
		Stringified: fmt.Sprintf("%f", left.Interface().(float64)*right.Interface().(float64)),
	}, nil
}

func evaluateDivide(left, right breeze.Const, op breeze.BinaryOp, datum datum.Datum) (breeze.Const, error) {
	if err := checkTypes(left, breeze.ConstKindNumber, right, breeze.ConstKindNumber); err != nil {
		return err.ToEmbeddedDatumErrorMessage(), nil
	}

	return breeze.Const{
		Kind:        breeze.ConstKindNumber,
		Stringified: fmt.Sprintf("%f", left.Interface().(float64)/right.Interface().(float64)),
	}, nil
}

func checkTypes(
	actualLeft breeze.Const, expectedLeft breeze.ConstKind,
	actualRight breeze.Const, expectedRight breeze.ConstKind,
) *breeze.TypeMismatchErr {
	if actualLeft.Kind != expectedLeft {
		return breeze.NewTypeMismatchErr(expectedLeft, actualLeft)
	}

	if actualRight.Kind != expectedRight {
		return breeze.NewTypeMismatchErr(expectedRight, actualRight)
	}

	return nil
}
