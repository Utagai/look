package data

import "encoding/json"

// Data is an in-memory group of JSON items.
type Data interface {
	// At returns the datum at the given index. If the given index is
	// out-of-bounds (< 0 || > .Length()), then this should return (nil, false).
	At(index int) (Datum, bool)
	Find(q string) Data
	Length() int
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
