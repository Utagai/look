package custom_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/utagai/look/config/custom"
)

type testCase struct {
	args                []string
	expectedParseFields []custom.Field
	expectedErr         error
}

func TestGetCustomParseArgs(t *testing.T) {
	for _, tc := range []testCase{
		{
			args: []string{`foo:number=LOL`},
			expectedParseFields: []custom.Field{
				{
					Type:      custom.FieldTypeNumber,
					Regex:     regexp.MustCompile(`LOL`),
					FieldName: "foo",
				},
			},
		},
		{
			args: []string{`foo:bool=hello!`},
			expectedParseFields: []custom.Field{
				{
					Type:      custom.FieldTypeBool,
					Regex:     regexp.MustCompile(`hello!`),
					FieldName: "foo",
				},
			},
		},
		{
			args: []string{`foo:string=whatever`},
			expectedParseFields: []custom.Field{
				{
					Type:      custom.FieldTypeString,
					Regex:     regexp.MustCompile(`whatever`),
					FieldName: "foo",
				},
			},
		},
		{
			args: []string{`foo:number`},
			expectedParseFields: []custom.Field{
				{
					Type:      custom.FieldTypeNumber,
					Regex:     regexp.MustCompile(`(\d+)`),
					FieldName: "foo",
				},
			},
		},
		{
			args: []string{`foo:bool`},
			expectedParseFields: []custom.Field{
				{
					Type:      custom.FieldTypeBool,
					Regex:     regexp.MustCompile(`(true|false)`),
					FieldName: "foo",
				},
			},
		},
		{
			args: []string{`foo:string`},
			expectedParseFields: []custom.Field{
				{
					Type:      custom.FieldTypeString,
					Regex:     regexp.MustCompile(`"(\w+)"`),
					FieldName: "foo",
				},
			},
		},
		{
			args: []string{`bar:string`},
			expectedParseFields: []custom.Field{
				{
					Type:      custom.FieldTypeString,
					Regex:     regexp.MustCompile(`"(\w+)"`),
					FieldName: "bar",
				},
			},
		},
		{
			args:        []string{`bar:notatype`},
			expectedErr: errors.New("unrecognized parse field type: \"notatype\" (valid: [string number bool])"),
		},
		{
			args:        []string{`bar|string`},
			expectedErr: errors.New("parse fields are of the format <fieldname>:<type>[=<regex>]"),
		},
		{
			args: []string{`foo:string`, `bar:number`, `baz:bool`},
			expectedParseFields: []custom.Field{
				{
					Type:      custom.FieldTypeString,
					Regex:     regexp.MustCompile(`"(\w+)"`),
					FieldName: "foo",
				},
				{
					Type:      custom.FieldTypeNumber,
					Regex:     regexp.MustCompile(`(\d+)`),
					FieldName: "bar",
				},
				{
					Type:      custom.FieldTypeBool,
					Regex:     regexp.MustCompile(`(true|false)`),
					FieldName: "baz",
				},
			},
		},
		{
			args:        []string{`bar:string`, `bar:number`},
			expectedErr: errors.New("duplicate field: \"bar\""),
		},
		{
			args: []string{`foo:string`, `bar:number=\d\d\d`, `baz:bool`},
			expectedParseFields: []custom.Field{
				{
					Type:      custom.FieldTypeString,
					Regex:     regexp.MustCompile(`"(\w+)"`),
					FieldName: "foo",
				},
				{
					Type:      custom.FieldTypeNumber,
					Regex:     regexp.MustCompile(`\d\d\d`),
					FieldName: "bar",
				},
				{
					Type:      custom.FieldTypeBool,
					Regex:     regexp.MustCompile(`(true|false)`),
					FieldName: "baz",
				},
			},
		},
		{
			args:                []string{},
			expectedParseFields: []custom.Field{},
		},
	} {
		testGetCustomParseArgs(t, tc)
	}
}

func testGetCustomParseArgs(t *testing.T, tc testCase) {
	customFields, err := custom.ParseFields(tc.args)
	if tc.expectedErr != nil {
		assert.EqualError(t, err, tc.expectedErr.Error())
	} else {
		parseFields := customFields.ParseFields
		if err != nil {
			t.Fatalf("did not expect an error, but got %v", err)
		}
		assert.Equal(t, tc.expectedParseFields, parseFields)
	}
}
