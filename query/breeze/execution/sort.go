package execution

import (
	"fmt"
	"sort"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func executeSort(sort *breeze.Sort, stream datum.DatumStream) *SortStream {
	return &SortStream{
		Sort:   sort,
		source: stream,
	}
}

// SortStream is an implementation of datum.Stream for the sort stage.
type SortStream struct {
	*breeze.Sort
	sortedDatums []datum.Datum
	source       datum.DatumStream
	sortedSource datum.DatumStream
}

// Next implements the DatumStream interface.
func (ss *SortStream) Next() (datum.Datum, error) {
	if ss.sortedDatums == nil {
		if err := ss.sortStream(); err != nil {
			return nil, fmt.Errorf("failed to sort data: %w", err)
		}
	}
	return ss.sortedSource.Next()
}

func (ss *SortStream) sortStream() error {
	// TODO: Dangerous if source is large.
	datums, err := datum.StreamToSlice(ss.source)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}
	sortable := sortableDatums{datums: datums, fieldName: ss.Field}
	sort.Sort(sortable)
	ss.sortedDatums = sortable.datums
	ss.sortedSource = datum.NewDatumSliceStream(sortable.datums)

	return nil
}

// TODO: This maybe should not be done in memory.
type sortableDatums struct {
	datums    []datum.Datum
	fieldName string
}

func (ds sortableDatums) Len() int {
	return len(ds.datums)
}

func (ds sortableDatums) Less(i, j int) bool {
	ithDoc := ds.datums[i]
	jthDoc := ds.datums[j]

	ithValue, ok := ithDoc[ds.fieldName]
	if !ok {
		// Treat documents where the field does not exist as being less than.
		return true
	}

	jthValue, ok := jthDoc[ds.fieldName]
	if !ok {
		// Ditto above.
		return true
	}

	return Compare(ithValue, jthValue) == Lesser
}

func (ds sortableDatums) Swap(i, j int) {
	ds.datums[i], ds.datums[j] = ds.datums[j], ds.datums[i]
}
