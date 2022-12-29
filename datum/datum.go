package datum

import (
	"encoding/json"
)

// Datum is a JSON object.
type Datum map[string]interface{}

func (d Datum) String() string {
	jsonString, err := json.Marshal(d)
	if err != nil {
		// Datums are always created from JSON initially, so this should never error.
		panic(err)
	}

	return string(jsonString)
}
