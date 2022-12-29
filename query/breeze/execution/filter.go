package execution

import (
	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func executeFilter(filter *breeze.Filter, stream datum.Stream) *FilterStream {
	return &FilterStream{
		Filter: filter,
		source: stream,
	}
}

// FilterStream is an implementation of datum.Stream for the filter stage.
type FilterStream struct {
	*breeze.Filter
	source datum.Stream
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

		for _, expr := range fs.Filter.Exprs {
			val, err := evaluateExpr(expr, datum)
			if err != nil {
				return nil, err
			}

			if val == nil {
				continue outer
			}

			if res, ok := val.(bool); ok {
				if !res {
					// If any of these results return false, then the whole
					// filter condition for the datum returns false as they are
					// all implicitly AND'd together.
					continue outer
				}
			} else {
				continue outer
			}
		}

		// If we get here, we have successfully evaluated against every check, and
		// we can be returned.
		return datum, nil
	}
}
