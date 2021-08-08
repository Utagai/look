package custom

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
		return `"(\w+)"`
	case FieldTypeNumber:
		return `(\d+)`
	case FieldTypeBool:
		return `(true|false)`
	default:
		panic(newUnrecognizedFieldTypeErr(string(typ)))
	}
}

// Field represents a JSON field to parse from the input lines.
type Field struct {
	Type      FieldType
	FieldName string
	Regex     *regexp.Regexp
}

// ParseFields returns a Fields for the given arguments.
func ParseFields(args []string) (*Fields, error) {
	seenFields := make(map[string]struct{}, len(args))
	fields := make([]Field, len(args))

	for i, arg := range args {
		field, err := parseField(arg)
		if err != nil {
			return nil, err
		}

		if _, ok := seenFields[field.FieldName]; ok {
			return nil, fmt.Errorf("duplicate field: %q", field.FieldName)
		}
		seenFields[field.FieldName] = struct{}{}

		fields[i] = field
	}

	return NewFields(fields)
}

const (
	typeSeparator  = ":"
	regexSeparator = "="
)

func parseField(arg string) (Field, error) {
	colonPos := strings.Index(arg, typeSeparator)
	if colonPos == -1 {
		// In this case, the colon is missing, and this is an invalid parse field.
		return Field{}, errors.New("parse fields are of the format <fieldname>:<type>[=<regex>]")
	}
	fieldName, typeAndMaybeRegex := arg[:colonPos], arg[colonPos+len(typeSeparator):]

	equalsPos := strings.LastIndex(typeAndMaybeRegex, regexSeparator)
	var typ FieldType
	var regexStr string
	var err error
	if equalsPos != -1 {
		var typeStr string
		typeStr, regexStr = typeAndMaybeRegex[:equalsPos], typeAndMaybeRegex[equalsPos+len(regexSeparator):]
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
		regexStr = defaultRegexForType(typ)
	}

	regex, err := regexp.Compile(regexStr)
	if err != nil {
		return Field{}, fmt.Errorf("invalid regex: %w", err)
	}

	return Field{
		Type:      typ,
		FieldName: fieldName,
		Regex:     regex,
	}, nil
}
