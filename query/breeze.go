package query

import (
	"fmt"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
	"github.com/utagai/look/query/breeze/execution"
)

type LiquidQueryExecutor struct{}

func NewLiquidQueryExecutor() *LiquidQueryExecutor {
	return &LiquidQueryExecutor{}
}

func (s *LiquidQueryExecutor) Find(q string, datums []datum.Datum) ([]datum.Datum, error) {
	p := breeze.NewParser(q)
	stages, err := p.Parse()
	if err != nil {
		parseErr := err.(*breeze.ParseError)
		return nil, fmt.Errorf("%w:\n%v", ErrUnableToParseQuery, parseErr.ErrorDescription())
	}

	var stream datum.DatumStream = datum.NewDatumSliceStream(datums)
	stream, err = execution.Execute(stream, stages)
	if err != nil {
		return nil, fmt.Errorf("failed to execute: %w", err)
	}

	return datum.StreamToSlice(stream)
}
