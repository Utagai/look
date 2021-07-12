package execution

import (
	"fmt"
	"strings"

	"github.com/utagai/look/query/breeze"
)

func executeBinaryOp(left interface{}, right *breeze.Const, op breeze.BinaryOp) bool {
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
