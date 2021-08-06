package config_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/utagai/look/config"
)

var testDataLines = []line{
	{
		Foo:     "bar",
		Iter:    8,
		Enabled: true,
	},
	{
		Foo:     "baz",
		Iter:    9,
		Enabled: true,
	},
	{
		Foo:     "qux",
		Iter:    10,
		Enabled: false,
	},
	{
		Foo:     "bar",
		Iter:    11,
		Enabled: true,
	},
}

type line struct {
	Foo     string `json:"foo"`
	Iter    int    `json:"iter"`
	Enabled bool   `json:"enabled"`
}

func newTestReader(t *testing.T, testData []line) io.Reader {
	var r bytes.Buffer
	for _, line := range testData {
		r.WriteString(fmt.Sprintf("value of foo: %q on iteration %d and enabled: %t\n", line.Foo, line.Iter, line.Enabled))
	}

	return &r
}

func TestCustomFieldsReader(t *testing.T) {
	r, err := config.NewCustomFieldsReader(newTestReader(t, testDataLines), []config.ParseField{
		{
			Type:      config.FieldTypeString,
			FieldName: "foo",
			Regex:     `"\w+"`,
		},
		{
			Type:      config.FieldTypeNumber,
			FieldName: "iter",
			Regex:     `\d+`,
		},
		{
			Type:      config.FieldTypeBool,
			FieldName: "enabled",
			Regex:     `true|false`,
		},
	})
	require.NoError(t, err)

	actualJSON, err := ioutil.ReadAll(r)
	if err != io.EOF {
		require.NoError(t, err)
	}

	actualLines := make([]line, len(testDataLines))
	require.NoError(t, json.Unmarshal(actualJSON, &actualLines))

	require.Equal(t,
		testDataLines,
		actualLines,
	)
}

func TestCustomFieldsReaderOnLargeInput(t *testing.T) {
	sizeIncreaseFactor := 100
	largeTestDataLines := make([]line, 0, len(testDataLines)*sizeIncreaseFactor)
	for i := 0; i < sizeIncreaseFactor; i++ {
		largeTestDataLines = append(largeTestDataLines, testDataLines...)
	}

	r, err := config.NewCustomFieldsReader(newTestReader(t, largeTestDataLines), []config.ParseField{
		{
			Type:      config.FieldTypeString,
			FieldName: "foo",
			Regex:     `"\w+"`,
		},
		{
			Type:      config.FieldTypeNumber,
			FieldName: "iter",
			Regex:     `\d+`,
		},
		{
			Type:      config.FieldTypeBool,
			FieldName: "enabled",
			Regex:     `true|false`,
		},
	})
	require.NoError(t, err)

	actualJSON, err := ioutil.ReadAll(r)
	if err != io.EOF {
		require.NoError(t, err)
	}

	actualLines := make([]line, len(testDataLines)*sizeIncreaseFactor)
	require.NoError(t, json.Unmarshal(actualJSON, &actualLines))

	require.Equal(t,
		largeTestDataLines,
		actualLines,
	)
}

func TestNothingMatches(t *testing.T) {
	testData := []byte(`
  foo
  bar
  baz
  `)
	r, err := config.NewCustomFieldsReader(bytes.NewBuffer(testData), []config.ParseField{
		{
			Type:      config.FieldTypeString,
			FieldName: "foo",
			Regex:     `"\w+"`,
		},
		{
			Type:      config.FieldTypeNumber,
			FieldName: "iter",
			Regex:     `\d+`,
		},
		{
			Type:      config.FieldTypeBool,
			FieldName: "enabled",
			Regex:     `true|false`,
		},
	})
	require.NoError(t, err)

	actualJSON, err := ioutil.ReadAll(r)
	if err != io.EOF {
		require.NoError(t, err)
	}

	actualLines := []line{}
	require.NoError(t, json.Unmarshal(actualJSON, &actualLines))

	require.Equal(t,
		[]line{},
		actualLines,
	)
}
