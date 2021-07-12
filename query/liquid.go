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
	stages, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnableToParseQuery, err)
	}

	var stream datum.DatumStream = datum.NewDatumSliceStream(datums)
	stream, err = execution.Execute(stream, stages)

	return datum.StreamToSlice(stream)
}
