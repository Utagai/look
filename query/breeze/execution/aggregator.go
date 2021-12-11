package execution

import (
	"math"
)

type aggregator interface {
	ingest(v interface{})
	aggregate() interface{}
}

// TODO: We could (and possibly should) convert this to use breeze.Concrete or
// something.
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
func (s *sum) ingest(v interface{}) {
	if floatValue, ok := convertPotentialNumber(v); ok {
		s.numberTotal += floatValue
		s.numNumbers += 1
	} else if boolValue, ok := convertPotentialBool(v); ok {
		s.boolTotal = s.boolTotal || boolValue
		s.numBools += 1
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

func (a *avg) ingest(v interface{}) {
	var ingestibleValue float64 = 0
	switch ta := v.(type) {
	case bool:
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
	minSet     bool
}

func (m *min) ingest(v interface{}) {
	if !m.minSet {
		m.minimumVal = v
		m.minSet = true
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

type mode struct {
	valueCounts *table
}

func (m *mode) ingest(v interface{}) {
	if m.valueCounts == nil {
		m.valueCounts = newTable()
	}

	untypedCnt, ok := m.valueCounts.GetOK(v)
	if ok {
		cnt := untypedCnt.(uint)
		m.valueCounts.Set(v, cnt+1)
	} else {
		m.valueCounts.Set(v, uint(1))
	}
}

func (m *mode) aggregate() interface{} {
	if m.valueCounts == nil {
		m.valueCounts = newTable()
	}

	keys := m.valueCounts.Keys()

	maxCount := -1
	var maxCountKey interface{} = nil
	for _, key := range keys {
		untypedCnt := m.valueCounts.Get(key)
		cnt := int(untypedCnt.(uint))
		if cnt > maxCount {
			maxCount = cnt
			maxCountKey = key
		}
	}

	return maxCountKey
}

// stddev is implemented via Welford's Algorithm for streamed variance.
type stddev struct {
	count uint
	mean  float64
	m2    float64
}

func (s *stddev) ingest(v interface{}) {
	floatVal, ok := v.(float64)
	if !ok {
		// If this is not a float value, then we actually want it to _not_ affect
		// the variance. Setting it to 0 or something similar would mean non-floats
		// would affect the variance, which would be incorrect. So, we skip it.
		return
	}
	s.count++
	delta := floatVal - s.mean
	s.mean += delta / float64(s.count)
	delta2 := floatVal - s.mean
	s.m2 += delta * delta2
}

func (s *stddev) aggregate() interface{} {
	if s.count < 2 {
		return "NaN"
	}

	return math.Sqrt(s.m2 / float64(s.count))
}
