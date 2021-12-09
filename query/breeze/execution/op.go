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

func executeBinaryOp(left interface{}, right breeze.Const, op breeze.BinaryCmpOp) bool {
	switch op {
	case breeze.BinaryOpEquals:
		return Compare(left, right.Interface()) == Equal
	case breeze.BinaryOpGeq:
		return Compare(left, right.Interface()) == Greater
	case breeze.BinaryOpContains:
		return strings.Contains(left.(string), right.Stringified)
	default:
		panic(fmt.Sprintf("unrecognized BinaryOp: %q", op))
	}
}
