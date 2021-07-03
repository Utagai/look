package datum

import (
	"encoding/json"
)

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
