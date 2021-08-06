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

func newCustomFieldsReader(t *testing.T, src io.Reader) io.Reader {
	r, err := config.NewCustomFieldsReader(src, []config.ParseField{
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

	return r
}

func runTest(t *testing.T, cfr io.Reader, expectedLines []line) {
	actualJSON, err := ioutil.ReadAll(cfr)
	if err != io.EOF {
		require.NoError(t, err)
	}

	actualLines := make([]line, len(testDataLines))
	require.NoError(t, json.Unmarshal(actualJSON, &actualLines))

	require.Equal(t,
		expectedLines,
		actualLines,
	)
}

func TestCustomFieldsReader(t *testing.T) {
	runTest(t, newCustomFieldsReader(t, newTestReader(t, testDataLines)), testDataLines)
}

func TestCustomFieldsReaderOnLargeInput(t *testing.T) {
	sizeIncreaseFactor := 100
	largeTestDataLines := make([]line, 0, len(testDataLines)*sizeIncreaseFactor)
	for i := 0; i < sizeIncreaseFactor; i++ {
		largeTestDataLines = append(largeTestDataLines, testDataLines...)
	}

	runTest(t, newCustomFieldsReader(t, newTestReader(t, largeTestDataLines)), largeTestDataLines)
}

func TestNothingMatches(t *testing.T) {
	testData := []byte(`
  foo
  bar
  baz
  `)
	runTest(t, newCustomFieldsReader(t, bytes.NewBuffer(testData)), []line{})
}
