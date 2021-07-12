package execution

import (
	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/liquid"
)

func executeCheck(check *liquid.Check, datum datum.Datum) bool {
	fieldValue, ok := datum[check.Field]
	if !ok {
		// If the field does not exist on this datum, evaluate to false.
		return false
	}

	return executeBinaryOp(fieldValue, check.Value, check.Op)
}
