package execution

import (
	"fmt"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func executeMap(filter *breeze.Map, stream datum.Stream) *MapStream {
	return &MapStream{
		Map:    filter,
		source: stream,
	}
}

// MapStream is an implementation of datum.Stream for the filter stage.
type MapStream struct {
	*breeze.Map
	source datum.Stream
}

// Next implements the datum.DatumStream interface.
func (fs *MapStream) Next() (datum.Datum, error) {
	// Keep iterating the stream and performing assignments per datum.
	datum, err := fs.source.Next()
	if err != nil {
		return nil, err
	}

	for _, assignment := range fs.Map.Assignments {
		// If we failed, move onto the next datum.
		if err := executeAssignment(assignment, datum); err != nil {
			return nil, fmt.Errorf("failed to execute assignment: %w", err)
		}
	}

	// If we get here, we have successfully re-assigned everything and can
	// return.
	return datum, nil
}
