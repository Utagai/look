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
		agg = &sum{}
	case breeze.AggFuncCount:
		agg = &sum{}
	case breeze.AggFuncMin:
		agg = &sum{}
	case breeze.AggFuncMax:
		agg = &sum{}
	case breeze.AggFuncMode:
		agg = &sum{}
	case breeze.AggFuncStdDev:
		agg = &sum{}
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

type aggregator interface {
	ingest(value interface{})
	aggregate() interface{}
}

type sum struct {
	numberTotal float64
	numNumbers  int

	// sum can handle both floats and bools (for bools, we are taking an OR of the
	// entire set). We track the counts of each so we know which one to give back
	// when Aggregate() is called.
	boolTotal bool
	numBools  int
}

// TODO: This can overflow, and so can some other agg functions.
func (s *sum) ingest(value interface{}) {
	if floatValue, ok := convertPotentialNumber(value); ok {
		s.numberTotal += floatValue
	} else if boolValue, ok := convertPotentialBool(value); ok {
		s.boolTotal = s.boolTotal || boolValue
	}
	// Otherwise, do nothing.
}

func (s *sum) aggregate() interface{} {
	if s.numNumbers >= s.numBools {
		return s.numberTotal
	}
	return s.boolTotal
}
