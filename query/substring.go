package query

import (
	"strings"

	"github.com/utagai/look/datum"
)

type SubstringQueryExecutor struct{}

func NewSubstringQueryExecutor() *SubstringQueryExecutor {
	return &SubstringQueryExecutor{}
}

func (s *SubstringQueryExecutor) Find(q string, datums []datum.Datum) ([]datum.Datum, error) {
	newDatums := []datum.Datum{}
	for _, datum := range datums {
		if strings.Contains(datum.String(), q) {
			newDatums = append(newDatums, datum)
		}
	}

	return newDatums, nil
}
