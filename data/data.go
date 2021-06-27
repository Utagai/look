package data

import (
	"context"
	"errors"

	"github.com/utagai/look/datum"
)

var (
	ErrOutOfBounds        = errors.New("index out of bounds")
	ErrUnableToParseQuery = errors.New("unable to parse the given query")
)

// Data is an in-memory group of JSON items.
// TODO: These methods should probably take a context and return an error (replacing bool when exists).
type Data interface {
	// At returns the datum at the given index. If the given index is
	// out-of-bounds (< 0 || > .Length()), then this should return (nil, false).
	At(ctx context.Context, index int) (datum.Datum, error)
	Find(ctx context.Context, query string) (Data, error)
	Length(context.Context) (int, error)
}
