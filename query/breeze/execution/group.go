package execution

import (
	"fmt"
	"io"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func executeGroup(group *breeze.Group, stream datum.Stream) *GroupStream {
	return &GroupStream{
		Group:  group,
		source: stream,
	}
}

// GroupStream is an implementation of datum.DatumStream for the group stage.
type GroupStream struct {
	*breeze.Group
	source        datum.Stream
	groupedSource datum.Stream
}

// Next implements the datum.DatumStream interface.
func (ss *GroupStream) Next() (datum.Datum, error) {
	if ss.groupedSource == nil {
		if err := ss.groupSource(); err != nil {
			return nil, fmt.Errorf("failed to group data: %w", err)
		}
	}

	return ss.groupedSource.Next()
}

func (ss *GroupStream) groupSource() error {
	sourcesToAggregate := ss.splitSource()

	aggregateResults := make([]datum.Datum, len(sourcesToAggregate))
	for i, source := range sourcesToAggregate {
		aggregateResults[i] = ss.aggregateStream(source)
	}

	ss.groupedSource = datum.NewSliceStream(aggregateResults)

	return nil
}

func (ss *GroupStream) splitSource() []datum.Stream {
	if ss.GroupByField == nil {
		// If there isn't a group by condition then we are simply aggregating over
		// the entire input, so return just the original input:
		return []datum.Stream{ss.source}
	}

	// Otherwise, we need to split apart the input stream by the group by field
	// and create N separate streams, each of which should then be independently
	// aggregated over (we do not do the aggregation here).
	table := newTable()
	for sourceDatum, err := ss.source.Next(); err == nil; sourceDatum, err = ss.source.Next() {
		groupByFieldValue, ok := sourceDatum[*ss.GroupByField]
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

		if vals, ok := table.GetOK(groupByFieldValue); ok {
			table.Set(
				groupByFieldValue,
				append(vals.([]datum.Datum), datum.Datum{ss.AggregateField: aggregateFieldValue}),
			)
		} else {
			table.Set(groupByFieldValue, []datum.Datum{{ss.AggregateField: aggregateFieldValue}})
		}
	}

	groupByFieldValues := table.Keys()
	splitSources := make([]datum.Stream, len(groupByFieldValues))
	for i, groupByFieldValue := range groupByFieldValues {
		aggregateFieldValues := table.Get(groupByFieldValue)
		splitSources[i] = datum.NewSliceStream(aggregateFieldValues.([]datum.Datum))
	}

	return splitSources
}

func (ss *GroupStream) getAggregator() aggregator {
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

func (ss *GroupStream) aggregateStream(input datum.Stream) datum.Datum {
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
