package data

import (
	"context"
	"fmt"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query"
)

// MemoryData is a data that lives entirely in memory.
type MemoryData struct {
	data     []datum.Datum
	executor query.Executor
}

var _ Data = (*MemoryData)(nil)

func NewMemoryData(data []datum.Datum) *MemoryData {
	return &MemoryData{
		data:     data,
		executor: query.NewSubstringQueryExecutor(),
	}
}

func (md *MemoryData) Find(_ context.Context, q string) (Data, error) {
	datums, err := md.executor.Find(q, md.data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute %q: %w", q, err)
	}
	return NewMemoryData(datums), nil
}

func (md *MemoryData) At(_ context.Context, index int) (datum.Datum, error) {
	if index >= len(md.data) || index < 0 {
		return nil, ErrOutOfBounds
	}

	return md.data[index], nil
}

func (md *MemoryData) Length(context.Context) (int, error) {
	return len(md.data), nil
}
