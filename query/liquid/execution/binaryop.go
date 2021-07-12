package execution

import (
	"fmt"
	"strings"

	"github.com/utagai/look/query/liquid"
)

func executeBinaryOp(left interface{}, right *liquid.Const, op liquid.BinaryOp) bool {
	switch op {
	case liquid.BinaryOpEquals:
		return Compare(left, right.Interface()) == Equal
	case liquid.BinaryOpGeq:
		return Compare(left, right.Interface()) == Greater
	case liquid.BinaryOpContains:
		return strings.Contains(left.(string), right.Stringified)
	default:
		panic(fmt.Sprintf("unrecognized BinaryOp: %q", op))
	}
}
