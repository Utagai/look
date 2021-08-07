package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// FieldType is the type annotation for the field to parse.
type FieldType string

const (
	// FieldTypeString corresponds to JSON strings.
	FieldTypeString = "string"
	// FieldTypeNumber corresponds to JSON numbers (integers/floats).
	FieldTypeNumber = "number"
	// FieldTypeBool corresponds to JSON booleans.
	FieldTypeBool = "bool"
)

var fieldTypes = []FieldType{
	FieldTypeString,
	FieldTypeNumber,
	FieldTypeBool,
}

func newUnrecognizedFieldTypeErr(unrecognizedField string) error {
	return fmt.Errorf("unrecognized parse field type: %q (valid: %v)", unrecognizedField, fieldTypes)
}

func toFieldType(typStr string) (FieldType, error) {
	switch typStr {
	case "string":
		return FieldTypeString, nil
	case "number":
		return FieldTypeNumber, nil
	case "bool":
		return FieldTypeBool, nil
	default:
		return "", newUnrecognizedFieldTypeErr(typStr)
	}

}

func defaultRegexForType(typ FieldType) string {
	switch typ {
	case FieldTypeString:
		return `"\w+"`
	case FieldTypeNumber:
		return `\d+`
	case FieldTypeBool:
		return `true|false`
	default:
		panic(newUnrecognizedFieldTypeErr(string(typ)))
	}
}

// Field represents a JSON field to parse from the input lines.
type Field struct {
	Type      FieldType
	FieldName string
	Regex     string
}

// GetCustomFields returns a CustomFields for the given arguments.
// TODO: OK, so I think what we actually want here is these functions (and all
// other custom_fields_* things to be moved to its own package. The tests can be
// internal and we can only expose the reader or something.)
func GetCustomFields(args []string) (*CustomFields, error) {
	seenFields := make(map[string]struct{}, len(args))
	parseFields := make([]Field, len(args))

	for i, arg := range args {
		parseField, err := parseParseField(arg)
		if err != nil {
			return nil, err
		}

		if _, ok := seenFields[parseField.FieldName]; ok {
			return nil, fmt.Errorf("duplicate field: %q", parseField.FieldName)
		}
		seenFields[parseField.FieldName] = struct{}{}

		parseFields[i] = parseField
	}

	return NewCustomFields(parseFields)
}

const (
	typeSeparator  = ":"
	regexSeparator = "="
)

func parseParseField(arg string) (Field, error) {
	colonPos := strings.LastIndex(arg, typeSeparator)
	if colonPos == -1 {
		// In this case, the colon is missing, and this is an invalid parse field.
		return Field{}, errors.New("parse fields are of the format <fieldname>:<type>[=<regex>]")
	}
	fieldName, typeAndMaybeRegex := arg[:colonPos], arg[colonPos+len(typeSeparator):]

	equalsPos := strings.LastIndex(typeAndMaybeRegex, regexSeparator)
	var typ FieldType
	var regex string
	var err error
	if equalsPos != -1 {
		var typeStr string
		typeStr, regex = typeAndMaybeRegex[:equalsPos], typeAndMaybeRegex[equalsPos+len(regexSeparator):]
		typ, err = toFieldType(typeStr)
		if err != nil {
			return Field{}, err
		}
	} else {
		// In this case, there is no regex so the user wants to default to the regex
		// for the type.
		typ, err = toFieldType(typeAndMaybeRegex)
		if err != nil {
			return Field{}, err
		}
		regex = defaultRegexForType(typ)
	}

	// TODO: I think we should store regexp.Regexp in ParseField instead.
	_, err = regexp.Compile(regex)
	if err != nil {
		return Field{}, fmt.Errorf("invalid regex: %w", err)
	}

	return Field{
		Type:      typ,
		FieldName: fieldName,
		Regex:     regex,
	}, nil
}
