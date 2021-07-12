package query

import (
	"fmt"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/liquid"
	"github.com/utagai/look/query/liquid/execution"
)

type LiquidQueryExecutor struct{}

func NewLiquidQueryExecutor() *LiquidQueryExecutor {
	return &LiquidQueryExecutor{}
}

func (s *LiquidQueryExecutor) Find(q string, datums []datum.Datum) ([]datum.Datum, error) {
	p := liquid.NewParser(q)
	stages, parseErr := p.Parse()
	if parseErr != nil {
		return nil, fmt.Errorf("%w:\n%v", ErrUnableToParseQuery, parseErr.ErrorDescription())
	}

	var stream datum.DatumStream = datum.NewDatumSliceStream(datums)
	stream, err := execution.Execute(stream, stages)
	if err != nil {
		return nil, fmt.Errorf("failed to execute: %w", err)
	}

	return datum.StreamToSlice(stream)
}
