package execution

import (
	"fmt"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

// Execute executes the given series of stages on the given datum stream.
func Execute(stream datum.DatumStream, stages []breeze.Stage) (datum.DatumStream, error) {
	for _, stage := range stages {
		var newStream datum.DatumStream
		switch ts := stage.(type) {
		case *breeze.Filter:
			newStream = executeFilter(ts, stream)
		case *breeze.Sort:
			newStream = executeSort(ts, stream)
		case *breeze.Reduce:
			newStream = executeReduce(ts, stream)
		case *breeze.Map:
			newStream = executeMap(ts, stream)
		default:
			return nil, fmt.Errorf("unrecognized query stage: %q", stage.Name())
		}

		return Execute(newStream, stages[1:])
	}

	return stream, nil
}
