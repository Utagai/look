package execution

import (
	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func executeUnaryCheck(check *breeze.UnaryCheck, datum datum.Datum) bool {
	// Allow nil field values because of the (!)exists unary operators.
	fieldValue := datum[check.Field]

	return executeUnaryOp(fieldValue, check.Op)
}

func executeBinaryCheck(check *breeze.BinaryCheck, datum datum.Datum) bool {
	fieldValue, ok := datum[check.Field]
	if !ok {
		// If the field does not exist on this datum, evaluate to false.
		return false
	}

	return executeBinaryOp(fieldValue, check.Value, check.Op)
}
