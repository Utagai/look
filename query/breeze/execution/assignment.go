package execution

import (
	"fmt"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func executeAssignment(assignment breeze.FieldAssignment, datum datum.Datum) error {
	fieldToAssign := assignment.Field
	newValue, err := evaluateExpr(assignment.Assignment, datum)
	if err != nil {
		return fmt.Errorf("failed to evaluate assignment: %w", err)
	}
	datum[fieldToAssign] = newValue
	return nil
}
