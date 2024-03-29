package execution

import (
	"fmt"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

// Execute executes the given series of stages on the given datum stream.
func Execute(stream datum.Stream, stages []breeze.Stage) (datum.Stream, error) {
	for _, stage := range stages {
		var newStream datum.Stream
		switch ts := stage.(type) {
		case *breeze.Filter:
			newStream = executeFilter(ts, stream)
		case *breeze.Sort:
			newStream = executeSort(ts, stream)
		case *breeze.Group:
			newStream = executeGroup(ts, stream)
		case *breeze.Map:
			newStream = executeMap(ts, stream)
		default:
			return nil, fmt.Errorf("unrecognized query stage: %q", stage.Name())
		}

		return Execute(newStream, stages[1:])
	}

	return stream, nil
}
