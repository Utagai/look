package execution

// TODO: Rename this file to op.go.

import (
	"fmt"
	"strings"

	"github.com/utagai/look/query/breeze"
)

func executeUnaryOp(left interface{}, op breeze.UnaryOp) bool {
	switch op {
	case breeze.UnaryOpExists:
		return left != nil
	case breeze.UnaryOpExistsNot:
		return left == nil
	default:
		panic(fmt.Sprintf("unrecognized BinaryOp: %q", op))
	}
}

func executeBinaryOp(left interface{}, right breeze.Concrete, op breeze.BinaryCmpOp) (bool, error) {
	rightIf, err := right.Interface()
	if err != nil {
		return false, err
	}

	switch op {
	case breeze.BinaryOpEquals:
		return Compare(left, rightIf) == Equal, nil
	case breeze.BinaryOpGeq:
		return Compare(left, rightIf) == Greater, nil
	case breeze.BinaryOpContains:
		return strings.Contains(left.(string), right.GetStringRepr()), nil
	default:
		panic(fmt.Sprintf("unrecognized BinaryOp: %q", op))
	}
}
