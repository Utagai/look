package config

import (
	"errors"
	"fmt"
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

type ParseField struct {
	Type      FieldType
	FieldName string
	Regex     string
}

func GetCustomParseFields(args []string) ([]ParseField, error) {
	seenFields := make(map[string]struct{}, len(args))
	parseFields := make([]ParseField, len(args))

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

	return parseFields, nil
}

const (
	typeSeparator  = ":"
	regexSeparator = "="
)

func parseParseField(arg string) (ParseField, error) {
	colonPos := strings.LastIndex(arg, typeSeparator)
	if colonPos == -1 {
		// In this case, the colon is missing, and this is an invalid parse field.
		return ParseField{}, errors.New("parse fields are of the format <fieldname>:<type>[=<regex>]")
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
			return ParseField{}, err
		}
	} else {
		// In this case, there is no regex so the user wants to default to the regex
		// for the type.
		typ, err = toFieldType(typeAndMaybeRegex)
		if err != nil {
			return ParseField{}, err
		}
		regex = defaultRegexForType(typ)
	}

	return ParseField{
		Type:      typ,
		FieldName: fieldName,
		Regex:     regex,
	}, nil
}
