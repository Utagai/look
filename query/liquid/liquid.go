package liquid

import (
	"fmt"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query"
)

type LiquidQueryExecutor struct{}

func NewLiquidQueryExecutor() *LiquidQueryExecutor {
	return &LiquidQueryExecutor{}
}

func (s *LiquidQueryExecutor) Find(q string, datums []datum.Datum) ([]datum.Datum, error) {
	p := NewParser(q)
	stages, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", query.ErrUnableToParseQuery, err)
	}

	var stream datum.DatumStream = datum.NewDatumSliceStream(datums)
	for _, stage := range stages {
		stream, err = stage.Execute(stream)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return datum.StreamToSlice(stream)
}
