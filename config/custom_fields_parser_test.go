package config_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utagai/look/config"
)

type testCase struct {
	args                []string
	expectedParseFields []config.ParseField
	expectedErr         error
}

func TestGetCustomParseArgs(t *testing.T) {

	for _, tc := range []testCase{
		{
			args: []string{`foo:number=LOL`},
			expectedParseFields: []config.ParseField{
				{
					Type:      config.FieldTypeNumber,
					Regex:     `LOL`,
					FieldName: "foo",
				},
			},
		},
		{
			args: []string{`foo:bool=hello!`},
			expectedParseFields: []config.ParseField{
				{
					Type:      config.FieldTypeBool,
					Regex:     `hello!`,
					FieldName: "foo",
				},
			},
		},
		{
			args: []string{`foo:string=whatever`},
			expectedParseFields: []config.ParseField{
				{
					Type:      config.FieldTypeString,
					Regex:     `whatever`,
					FieldName: "foo",
				},
			},
		},
		{
			args: []string{`foo:number`},
			expectedParseFields: []config.ParseField{
				{
					Type:      config.FieldTypeNumber,
					Regex:     `\d+`,
					FieldName: "foo",
				},
			},
		},
		{
			args: []string{`foo:bool`},
			expectedParseFields: []config.ParseField{
				{
					Type:      config.FieldTypeBool,
					Regex:     `true|false`,
					FieldName: "foo",
				},
			},
		},
		{
			args: []string{`foo:string`},
			expectedParseFields: []config.ParseField{
				{
					Type:      config.FieldTypeString,
					Regex:     `"\w+"`,
					FieldName: "foo",
				},
			},
		},
		{
			args: []string{`bar:string`},
			expectedParseFields: []config.ParseField{
				{
					Type:      config.FieldTypeString,
					Regex:     `"\w+"`,
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
			expectedParseFields: []config.ParseField{
				{
					Type:      config.FieldTypeString,
					Regex:     `"\w+"`,
					FieldName: "foo",
				},
				{
					Type:      config.FieldTypeNumber,
					Regex:     `\d+`,
					FieldName: "bar",
				},
				{
					Type:      config.FieldTypeBool,
					Regex:     `true|false`,
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
			expectedParseFields: []config.ParseField{
				{
					Type:      config.FieldTypeString,
					Regex:     `"\w+"`,
					FieldName: "foo",
				},
				{
					Type:      config.FieldTypeNumber,
					Regex:     `\d\d\d`,
					FieldName: "bar",
				},
				{
					Type:      config.FieldTypeBool,
					Regex:     `true|false`,
					FieldName: "baz",
				},
			},
		},
	} {
		testGetCustomParseArgs(t, tc)
	}
}

func testGetCustomParseArgs(t *testing.T, tc testCase) {
	parseFields, err := config.GetCustomParseFields(tc.args)
	if tc.expectedErr != nil {
		assert.EqualError(t, err, tc.expectedErr.Error())
	} else {
		if err != nil {
			t.Fatalf("did not expect an error, but got %v", err)
		}
		assert.Equal(t, tc.expectedParseFields, parseFields)
	}
}
