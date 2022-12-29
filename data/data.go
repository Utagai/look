package data

import (
	"context"
	"errors"

	"github.com/utagai/look/datum"
)

var (
	ErrOutOfBounds = errors.New("index out of bounds")
)

// Data is an in-memory group of JSON items.
type Data interface {
	// At returns the datum at the given index. If the given index is
	// out-of-bounds (< 0 || > .Length()), then this should return (nil, false).
	At(ctx context.Context, index int) (datum.Datum, error)
	// Find executes the given query against the data. It returns
	// another Data that represents the result set from running the
	// given query.
	Find(ctx context.Context, query string) (Data, error)
	// Length returns the number of datums in this Data.
	Length(context.Context) (int, error)
}
