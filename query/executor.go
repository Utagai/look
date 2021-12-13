package query

import (
	"errors"

	"github.com/utagai/look/datum"
)

var ErrUnableToParseQuery = errors.New("unable to parse the given query")

// Executor executes a given query string on a set of datums and returns a
// resulting set of datums.
// TODO: These execute calls should be debounced.
// TODO: Find is not a great name for this. We can do more than just find
// things.
type Executor interface {
	Find(q string, datums []datum.Datum) ([]datum.Datum, error)
}
