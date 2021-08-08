package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ErrNoMatch is an error returned when a given line produces no match from
// which to generate its JSON representation.
var ErrNoMatch = errors.New("no match for line")

// CustomFields represents the custom fields defined by the user.
// TODO: Does this need to be exported?
type CustomFields struct {
	ParseFields []Field
	TotalRegex  *regexp.Regexp
}

// NewCustomFields is a constructor for CustomFields.
func NewCustomFields(parseFields []Field) (*CustomFields, error) {
	var totalRegexStr strings.Builder
	// Leading extra characters.
	totalRegexStr.WriteString(".*?")
	for _, pf := range parseFields {
		totalRegexStr.WriteString(fmt.Sprintf("(%s)", pf.Regex))
		totalRegexStr.WriteString(".*?")
	}

	totalRegex, err := regexp.Compile(totalRegexStr.String())
	if err != nil {
		return nil, fmt.Errorf("failed to compile the complete regex: %w", err)
	}

	return &CustomFields{
		ParseFields: parseFields,
		TotalRegex:  totalRegex,
	}, nil
}

// ToJSON converts the given line into a slice of bytes representing its
// converted JSON representation.
// The returned JSON may not be 'complete', in that it may not have a field for
// each specified custom field.
func (cf *CustomFields) ToJSON(line string) ([]byte, error) {
	matches := cf.TotalRegex.FindStringSubmatch(line)
	if len(matches) == 0 {
		return nil, ErrNoMatch
	}
	// Discard the initial match of the entire string -- we only care about the
	// capture groups.
	matches = matches[1:]

	jsonMap := make(map[string]interface{}, len(cf.ParseFields))
	for i, pf := range cf.ParseFields {
		match := matches[i]
		if len(match) == 0 {
			// An empty string match implies either the regex matches the empty
			// string. This is by definition an optional regex, so set this to null
			// and move on.
			jsonMap[pf.FieldName] = nil
			continue
		}

		var value interface{}
		switch pf.Type {
		case FieldTypeBool:
			value = match == "true"
		case FieldTypeNumber:
			f64, err := strconv.ParseFloat(match, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse numeric value: %w", err)
			}
			value = f64
		case FieldTypeString:
			value = strings.Trim(match, "\"")
		default:
			panic(fmt.Sprintf("unrecognized field type: %q", pf.Type))
		}

		jsonMap[pf.FieldName] = value
	}

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		return jsonBytes, fmt.Errorf("failed to serialize to JSON: %w", err)
	}
	return jsonBytes, nil
}
