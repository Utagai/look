package data

import (
	"context"
	"encoding/json"
	"errors"
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
	At(ctx context.Context, index int) (Datum, error)
	Find(ctx context.Context, query string) (Data, error)
	Length(context.Context) (int, error)
}

// Datum is a JSON object.
type Datum map[string]interface{}

func (d Datum) String() string {
	jsonString, err := json.Marshal(d)
	if err != nil {
		// TODO: Should we handle this a bit better?
		panic(err)
	}

	return string(jsonString)
}
