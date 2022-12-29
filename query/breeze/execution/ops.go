package execution

import (
	"fmt"
	"strings"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func evaluateOp(left, right breeze.Concrete, op breeze.BinaryOp, datum datum.Datum) (breeze.Concrete, error) {
	switch op {
	case breeze.BinaryOpPlus:
		return evaluatePlus(left, right, op, datum)
	case breeze.BinaryOpMinus:
		return evaluateMinus(left, right, op, datum)
	case breeze.BinaryOpMultiply:
		return evaluateMultiply(left, right, op, datum)
	case breeze.BinaryOpDivide:
		return evaluateDivide(left, right, op, datum)
	case breeze.BinaryOpEquals:
		return evaluateComparisonOp(left, right, op, datum)
	case breeze.BinaryOpGeq:
		return evaluateComparisonOp(left, right, op, datum)
	case breeze.BinaryOpContains:
		return evaluateContains(left, right, datum)
	default:
		panic(fmt.Sprintf("unrecognized operator: %q", op))
	}
}

func evaluatePlus(left, right breeze.Concrete, op breeze.BinaryOp, datum datum.Datum) (*breeze.Scalar, error) {
	if err := checkTypes(left, breeze.ScalarKindNumber, right, breeze.ScalarKindNumber); err != nil {
		return err.ToEmbeddedDatumErrorMessage(), nil
	}

	leftNum, rightNum, err := getLeftAndRightNums(left, right)
	if err != nil {
		return nil, err
	}

	return &breeze.Scalar{
		Kind:        breeze.ScalarKindNumber,
		Stringified: fmt.Sprintf("%f", leftNum+rightNum),
	}, nil
}

func evaluateMinus(left, right breeze.Concrete, op breeze.BinaryOp, datum datum.Datum) (*breeze.Scalar, error) {
	if err := checkTypes(left, breeze.ScalarKindNumber, right, breeze.ScalarKindNumber); err != nil {
		return err.ToEmbeddedDatumErrorMessage(), nil
	}

	leftNum, rightNum, err := getLeftAndRightNums(left, right)
	if err != nil {
		return nil, err
	}

	return &breeze.Scalar{
		Kind:        breeze.ScalarKindNumber,
		Stringified: fmt.Sprintf("%f", leftNum+rightNum),
	}, nil
}

func evaluateMultiply(left, right breeze.Concrete, op breeze.BinaryOp, datum datum.Datum) (*breeze.Scalar, error) {
	if err := checkTypes(left, breeze.ScalarKindNumber, right, breeze.ScalarKindNumber); err != nil {
		return err.ToEmbeddedDatumErrorMessage(), nil
	}

	leftNum, rightNum, err := getLeftAndRightNums(left, right)
	if err != nil {
		return nil, err
	}

	return &breeze.Scalar{
		Kind:        breeze.ScalarKindNumber,
		Stringified: fmt.Sprintf("%f", leftNum*rightNum),
	}, nil
}

func evaluateDivide(left, right breeze.Concrete, op breeze.BinaryOp, datum datum.Datum) (*breeze.Scalar, error) {
	if err := checkTypes(left, breeze.ScalarKindNumber, right, breeze.ScalarKindNumber); err != nil {
		return err.ToEmbeddedDatumErrorMessage(), nil
	}

	leftNum, rightNum, err := getLeftAndRightNums(left, right)
	if err != nil {
		return nil, err
	}

	return &breeze.Scalar{
		Kind:        breeze.ScalarKindNumber,
		Stringified: fmt.Sprintf("%f", leftNum/rightNum),
	}, nil
}

func evaluateComparisonOp(left, right breeze.Concrete, op breeze.BinaryOp, datum datum.Datum) (*breeze.Scalar, error) {
	var desiredComparison Comparison = Equal
	switch op {
	case breeze.BinaryOpEquals:
		desiredComparison = Equal
	case breeze.BinaryOpGeq:
		desiredComparison = Lesser
	default:
		panic(fmt.Sprintf("unreachable: evaluateComparisonOp() called with %q", op))
	}

	leftIf, err := left.Interface()
	if err != nil {
		return nil, err
	}

	rightIf, err := right.Interface()
	if err != nil {
		return nil, err
	}

	return boolToConcrete(Compare(leftIf, rightIf) == desiredComparison), nil
}

func evaluateContains(left, right breeze.Concrete, datum datum.Datum) (*breeze.Scalar, error) {
	if err := checkTypes(left, breeze.ScalarKindString, right, breeze.ScalarKindString); err != nil {
		return err.ToEmbeddedDatumErrorMessage(), nil
	}

	return boolToConcrete(strings.Contains(left.GetStringRepr(), right.GetStringRepr())), nil
}

func checkTypes(
	actualLeft breeze.Concrete, expectedLeft breeze.ScalarKind,
	actualRight breeze.Concrete, expectedRight breeze.ScalarKind,
) *breeze.TypeMismatchErr {
	if actualLeft.ConcreteKind() != breeze.ConcreteKindScalar && actualLeft.(*breeze.Scalar).Kind != expectedLeft {
		return breeze.NewTypeMismatchErr(string(expectedLeft), actualLeft)
	}

	if actualRight.ConcreteKind() != breeze.ConcreteKindScalar && actualRight.(*breeze.Scalar).Kind != expectedRight {
		return breeze.NewTypeMismatchErr(string(expectedRight), actualRight)
	}

	return nil
}

func getLeftAndRightNums(left breeze.Concrete, right breeze.Concrete) (float64, float64, error) {
	leftNum, err := left.Interface()
	if err != nil {
		return 0, 0, err
	}
	rightNum, err := right.Interface()
	if err != nil {
		return 0, 0, err
	}

	return leftNum.(float64), rightNum.(float64), nil
}

func boolToConcrete(b bool) *breeze.Scalar {
	return &breeze.Scalar{
		Kind:        breeze.ScalarKindBool,
		Stringified: fmt.Sprintf("%t", b),
	}
}
