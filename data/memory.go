package data

import (
	"context"
	"strings"
)

// MemoryData is a data that lives entirely in memory.
type MemoryData struct {
	data []Datum
}

var _ Data = (*MemoryData)(nil)

func NewMemoryData(data []Datum) *MemoryData {
	return &MemoryData{
		data: data,
	}
}

func (md *MemoryData) Find(_ context.Context, q string) (Data, error) {
	newDatums := make([]Datum, 0, len(md.data))
	for _, datum := range md.data {
		if strings.Contains(datum.String(), q) {
			newDatums = append(newDatums, datum)
		}
	}
	return NewMemoryData(newDatums), nil
}

func (md *MemoryData) At(_ context.Context, index int) (Datum, error) {
	if index >= len(md.data) || index < 0 {
		return nil, ErrOutOfBounds
	}

	return md.data[index], nil
}

func (md *MemoryData) Length(context.Context) (int, error) {
	return len(md.data), nil
}
