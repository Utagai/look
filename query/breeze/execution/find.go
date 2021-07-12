package execution

import (
	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func executeFind(find *breeze.Find, stream datum.DatumStream) *FindStream {
	return &FindStream{
		Find:   find,
		source: stream,
	}
}

// FindStream is an implementation of datum.Stream for the find stage.
type FindStream struct {
	*breeze.Find
	source datum.DatumStream
}

// Next implements the datum.DatumStream interface.
func (fs *FindStream) Next() (datum.Datum, error) {
	// Keep iterating the stream until something passes the checks.
outer:
	for {
		datum, err := fs.source.Next()
		if err != nil {
			return nil, err
		}

		for _, check := range fs.Find.Checks {
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
