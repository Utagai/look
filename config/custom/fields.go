package custom

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrNoMatch is an error returned when a given line produces no match from
// which to generate its JSON representation.
var ErrNoMatch = errors.New("no match for line")

// errNoMatchForField is an internal error that serves a similar purpose to
// ErrNoMatch, but is used in internal helpers to indicate when a particular
// field has failed to match, but understands that future fields may match and
// therefore the line may not necessarily return ErrNoMatch.
var errNoMatchForField = errors.New("no match for current field")

// Fields represents the custom fields defined by the user.
type Fields struct {
	CustomFields []Field
}

// NewFields is a constructor for Fields.
func NewFields(customFields []Field) (*Fields, error) {
	return &Fields{
		CustomFields: customFields,
	}, nil
}

// ToJSON converts the given line into a slice of bytes representing its
// converted JSON representation.
// The returned JSON may not be 'complete', in that it may not have a field for
// each specified custom field.
func (cf *Fields) ToJSON(line string) ([]byte, error) {
	jsonMap := make(map[string]interface{}, len(cf.CustomFields))
	var err error
	remainingLine := line
	atLeastOneNonNullValue := false
	for _, f := range cf.CustomFields {
		remainingLine, err = cf.addToMap(jsonMap, f, remainingLine)
		if err != nil && err != errNoMatchForField {
			return nil, err // Fatal.
		}
		atLeastOneNonNullValue = atLeastOneNonNullValue || err == nil
	}

	// If none of the values were found, then that means the entire line failed to
	// match anything.
	if !atLeastOneNonNullValue {
		return nil, ErrNoMatch
	}

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		return jsonBytes, fmt.Errorf("failed to serialize to JSON: %w", err)
	}
	return jsonBytes, nil
}

func (cf *Fields) addToMap(jsonMap map[string]interface{}, f Field, line string) (string, error) {
	matchIndices := f.Regex.FindStringSubmatchIndex(line)
	if len(matchIndices) <= 0 {
		// If this parse field finds no match, then it does not exist and should
		// be null.
		jsonMap[f.FieldName] = nil
		return line, errNoMatchForField
	}

	// matchIndices refers to a slice of index _pairs_, so the first two refers
	// to the entire match, even though what we actually want is the capture
	// group:
	matchIndices = matchIndices[2:]

	// Extract the match itself:
	match := line[matchIndices[0]:matchIndices[1]]
	// Advance the line for future regex matches:
	line = line[matchIndices[1]:]

	var value interface{}
	switch f.Type {
	case FieldTypeBool:
		value = match == "true"
	case FieldTypeNumber:
		f64, err := strconv.ParseFloat(match, 64)
		if err != nil {
			return line, fmt.Errorf("failed to parse numeric value: %w", err)
		}
		value = f64
	case FieldTypeString:
		value = strings.Trim(match, "\"")
	default:
		panic(fmt.Sprintf("unrecognized field type: %q", f.Type))
	}

	jsonMap[f.FieldName] = value
	return line, nil
}
