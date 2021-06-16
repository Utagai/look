package data

import "strings"

// MemoryData is a data that lives entirely in memory.
type MemoryData struct {
	data []Datum
}

func NewMemoryData(data []Datum) *MemoryData {
	return &MemoryData{
		data: data,
	}
}

var _ Data = (*MemoryData)(nil)

func (md *MemoryData) Find(q string) Data {
	newDatums := make([]Datum, 0, md.Length())
	for _, datum := range md.data {
		if strings.Contains(datum.String(), q) {
			newDatums = append(newDatums, datum)
		}
	}
	return NewMemoryData(newDatums)
}

func (md *MemoryData) At(index int) (Datum, bool) {
	if index >= len(md.data) || index < 0 {
		return nil, false
	}

	return md.data[index], true
}

func (md *MemoryData) Length() int {
	return len(md.data)
}
