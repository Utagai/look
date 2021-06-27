package query

import "github.com/utagai/look/datum"

// Executor executes a given query string on a set of datums and returns a
// resulting set of datums.
type Executor interface {
	Find(q string, datums []datum.Datum) ([]datum.Datum, error)
}
