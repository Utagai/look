package execution

import (
	"fmt"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

// TODO: This should probably exist in the datum package?
func makeDatumCopy(d datum.Datum) datum.Datum {
	copy := datum.Datum{}

	// TODO: This does not handle nested datums.
	for k, v := range d {
		copy[k] = v
	}

	return copy
}

func executeAssignment(assignment breeze.FieldAssignment, datum datum.Datum) (datum.Datum, error) {
	fieldToAssign := assignment.Field
	datumCopy := makeDatumCopy(datum)
	newValue, err := evaluateExpr(assignment.Assignment, datum)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate assignment: %w", err)
	}
	datumCopy[fieldToAssign] = newValue
	// TODO: This should eventually return something...
	return datumCopy, nil
}
