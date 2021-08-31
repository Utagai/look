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

type avg struct {
	total     sum
	numValues int
}

func (a *avg) ingest(value interface{}) {
	var ingestibleValue float64 = 0
	switch ta := value.(type) {
	case bool:
		ingestibleValue = 0
		if ta {
			ingestibleValue = 1
		}
	case float64:
		ingestibleValue = ta
	default:
		ingestibleValue = 0
	}
	a.total.ingest(ingestibleValue)
	a.numValues++
}

func (a *avg) aggregate() interface{} {
	totalSum := a.total.aggregate()
	switch tsum := totalSum.(type) {
	case float64:
		if a.numValues == 0 {
			return 0
		}
		return tsum / float64(a.numValues)
	default:
		panic("TODO")
	}
}

type count struct {
	numValues uint
}

func (c *count) ingest(_ interface{}) {
	c.numValues++
}

func (c *count) aggregate() interface{} {
	return c.numValues
}

// min and max can be implemented with the same general type, but since the
// amount of code duplication reduction will be small, we choose to duplicate
// it with simpler code.
type min struct {
	minimumVal interface{}
}

func (m *min) ingest(v interface{}) {
	if m.minimumVal == nil {
		m.minimumVal = v
		return
	}

	switch Compare(v, m.minimumVal) {
	case Lesser:
		m.minimumVal = v
	}
}

func (m *min) aggregate() interface{} {
	return m.minimumVal
}

type max struct {
	maximumVal interface{}
}

func (m *max) ingest(v interface{}) {
	switch Compare(v, m.maximumVal) {
	case Greater:
		m.maximumVal = v
	}
}

func (m *max) aggregate() interface{} {
	return m.maximumVal
}
