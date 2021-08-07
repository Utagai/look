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
	newLine("bar", 8, true),
	newLine("baz", 9, true),
	newLine("qux", 10, false),
	newLine("bar", 11, true),
}

type line struct {
	Foo     *string `json:"foo"`
	Iter    *int    `json:"iter"`
	Enabled *bool   `json:"enabled"`
}

func newLine(foo string, iter int, enabled bool) line {
	return line{
		Foo:     &foo,
		Iter:    &iter,
		Enabled: &enabled,
	}
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func newTestReader(t *testing.T, testData []line) io.Reader {
	var r bytes.Buffer
	for _, line := range testData {
		// Just to catch programmer mistakes.
		if line.Foo == nil || line.Iter == nil || line.Enabled == nil {
			t.Fatalf("newTestReader should only be called with non-nil field values for each line")
		}

		r.WriteString(fmt.Sprintf("value of foo: %#v on iteration %v and enabled: %v\n", *line.Foo, *line.Iter, *line.Enabled))
	}

	return &r
}

func newCustomFieldsReader(t *testing.T, src io.Reader) io.Reader {
	customFields, err := config.NewCustomFields([]config.ParseField{
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
	r, err := config.NewCustomFieldsReader(src, customFields)
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

// This test is effectively testing how we handle heterogeneous schemas -- for
// example, if one line matches everything but another line only matches some of
// the fields, and another matches none.
func TestPartialMatches(t *testing.T) {
	testData := []byte(`
  value of foo: "bar" on iteration 8 and enabled: true
  value of foo: "bar" on and enabled: true
  value of foo: "baz" on iteration 9 and enabled: true
  no match
  value of foo: "qux" on iteration 10 and enabled: false
  `)

	runTest(t, newCustomFieldsReader(t, bytes.NewBuffer(testData)), []line{
		newLine("bar", 8, true),
		{
			Foo:     strPtr("bar"),
			Iter:    nil,
			Enabled: boolPtr(true),
		},
		newLine("baz", 9, true),
		newLine("bar", 11, true),
	})
}
