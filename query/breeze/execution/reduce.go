package execution

import (
	"fmt"
	"io"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func executeReduce(reduce *breeze.Reduce, stream datum.DatumStream) *ReduceStream {
	return &ReduceStream{
		Reduce: reduce,
		source: stream,
	}
}

// ReduceStream is an implementation of datum.DatumStream for the reduce stage.
type ReduceStream struct {
	*breeze.Reduce
	source         datum.DatumStream
	reduceedSource datum.DatumStream
}

// Next implements the datum.DatumStream interface.
func (ss *ReduceStream) Next() (datum.Datum, error) {
	if ss.reduceedSource == nil {
		if err := ss.reduceSource(); err != nil {
			return nil, fmt.Errorf("failed to reduce data: %w", err)
		}
	}

	return ss.reduceedSource.Next()
}

func (ss *ReduceStream) reduceSource() error {
	sourcesToAggregate := ss.splitSource()

	aggregateResults := make([]datum.Datum, len(sourcesToAggregate))
	for i, source := range sourcesToAggregate {
		aggregateResults[i] = ss.aggregateStream(source)
	}

	ss.reduceedSource = datum.NewDatumSliceStream(aggregateResults)

	return nil
}

func (ss *ReduceStream) splitSource() []datum.DatumStream {
	if ss.ReduceByField == nil {
		// If there isn't a reduce by condition then we are simply aggregating over
		// the entire input, so return just the original input:
		return []datum.DatumStream{ss.source}
	}

	// Otherwise, we need to split apart the input stream by the reduce by field
	// and create N separate streams, each of which should then be independently
	// aggregated over (we do not do the aggregation here).
	table := newTable()
	for sourceDatum, err := ss.source.Next(); err == nil; sourceDatum, err = ss.source.Next() {
		reduceByFieldValue, ok := sourceDatum[*ss.ReduceByField]
		if !ok {
			// If the field doesn't exist, ignore the document.
			continue
		}

		aggregateFieldValue, ok := sourceDatum[ss.AggregateField]
		if !ok {
			// If the aggregate field value for this doesn't exist, we should
			// similarly also ignore it.
			continue
		}

		if vals, ok := table.GetOK(reduceByFieldValue); ok {
			table.Set(
				reduceByFieldValue,
				append(vals.([]datum.Datum), datum.Datum{ss.AggregateField: aggregateFieldValue}),
			)
		} else {
			table.Set(reduceByFieldValue, []datum.Datum{{ss.AggregateField: aggregateFieldValue}})
		}
	}

	reduceByFieldValues := table.Keys()
	splitSources := make([]datum.DatumStream, len(reduceByFieldValues))
	for i, reduceByFieldValue := range reduceByFieldValues {
		aggregateFieldValues := table.Get(reduceByFieldValue)
		splitSources[i] = datum.NewDatumSliceStream(aggregateFieldValues.([]datum.Datum))
	}

	return splitSources
}

func (ss *ReduceStream) getAggregator() aggregator {
	switch ss.AggFunc {
	case breeze.AggFuncSum:
		return &sum{}
	case breeze.AggFuncAvg:
		return &avg{}
	case breeze.AggFuncCount:
		return &count{}
	case breeze.AggFuncMin:
		return &min{}
	case breeze.AggFuncMax:
		return &max{}
	case breeze.AggFuncMode:
		return &mode{}
	case breeze.AggFuncStdDev:
		return &stddev{}
	default:
		panic(fmt.Sprintf("unrecognized aggregate function: %q", ss.AggFunc))
	}
}

func (ss *ReduceStream) aggregateStream(input datum.DatumStream) datum.Datum {
	agg := ss.getAggregator()

	for datum, err := input.Next(); err != io.EOF; datum, err = input.Next() {
		fieldValue, ok := datum[ss.AggregateField]
		if !ok {
			// If the field doesn't exist, ignore the document.
			continue
		}
		agg.ingest(fieldValue)
	}

	return datum.Datum{
		ss.AggregateField: agg.aggregate(),
	}
}
