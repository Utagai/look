package execution

import (
	"fmt"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/liquid"
)

// Execute executes the given series of stages on the given datum stream.
func Execute(stream datum.DatumStream, stages []liquid.Stage) (datum.DatumStream, error) {
	for _, stage := range stages {
		var newStream datum.DatumStream
		switch ts := stage.(type) {
		case *liquid.Find:
			newStream = executeFind(ts, stream)
		case *liquid.Sort:
			newStream = executeSort(ts, stream)
		default:
			return nil, fmt.Errorf("unrecognized query stage: %q", stage.Name())
		}

		return Execute(newStream, stages[1:])
	}

	return stream, nil
}
