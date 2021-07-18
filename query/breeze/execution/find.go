package execution

import (
	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func executeFilter(filter *breeze.Filter, stream datum.DatumStream) *FilterStream {
	return &FilterStream{
		Filter: filter,
		source: stream,
	}
}

// FilterStream is an implementation of datum.Stream for the filter stage.
type FilterStream struct {
	*breeze.Filter
	source datum.DatumStream
}

// Next implements the datum.DatumStream interface.
func (fs *FilterStream) Next() (datum.Datum, error) {
	// Keep iterating the stream until something passes the checks.
outer:
	for {
		datum, err := fs.source.Next()
		if err != nil {
			return nil, err
		}

		for _, check := range fs.Filter.Checks {
			// If we failed, move onto the next datum.
			if !executeCheck(check, datum) {
				continue outer
			}
		}

		// If we get here, we have successfully evaluated against every check, and
		// we can be returned.
		return datum, nil
	}
}
