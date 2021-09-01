package execution

import (
	"fmt"
	"io"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func executeGroup(group *breeze.Group, stream datum.DatumStream) *GroupStream {
	return &GroupStream{
		Group:  group,
		source: stream,
	}
}

// GroupStream is an implementation of datum.DatumStream for the group stage.
type GroupStream struct {
	*breeze.Group
	source        datum.DatumStream
	groupedSource datum.DatumStream
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
	var agg aggregator
	switch ss.AggFunc {
	case breeze.AggFuncSum:
		agg = &sum{}
	case breeze.AggFuncAvg:
		agg = &avg{}
	case breeze.AggFuncCount:
		agg = &count{}
	case breeze.AggFuncMin:
		agg = &min{}
	case breeze.AggFuncMax:
		agg = &max{}
	case breeze.AggFuncMode:
		agg = &mode{}
	case breeze.AggFuncStdDev:
		agg = &stddev{}
	default:
		panic(fmt.Sprintf("unrecognized aggregate function: %q", ss.AggFunc))
	}

	for datum, err := ss.source.Next(); err != io.EOF; datum, err = ss.source.Next() {
		fieldValue, ok := datum[ss.Field]
		if !ok {
			// If the field doesn't exist, ignore the document.
			continue
		}
		agg.ingest(fieldValue)
	}

	ss.groupedSource = datum.NewUnaryStream(datum.Datum{
		ss.Field: agg.aggregate(),
	})

	return nil
}
