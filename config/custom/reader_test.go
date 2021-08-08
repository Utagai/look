package custom_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/utagai/look/config/custom"
)

const (
	TestMaxBufSizeBytes = 100
)

var testDataLines = []line{
	newLine("bar", 8.2, true),
	newLine("baz", 9, true),
	newLine("qux", -10.1, false),
	newLine("bar", 11, true),
}

type line struct {
	Foo     *string  `json:"foo"`
	Iter    *float64 `json:"iter"`
	Enabled *bool    `json:"enabled"`
}

func newLine(foo string, iter float64, enabled bool) line {
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

func newCustomFieldsReaderWithFields(t *testing.T, src io.Reader, fields *custom.Fields) io.Reader {
	r, err := custom.NewFieldsReader(src, fields, TestMaxBufSizeBytes)
	require.NoError(t, err)

	return r

}

func newCustomFieldsReader(t *testing.T, src io.Reader) io.Reader {
	fields, err := custom.ParseFields([]string{"foo:string", "iter:number", "enabled:bool"})
	require.NoError(t, err)
	return newCustomFieldsReaderWithFields(t, src, fields)
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

	cfr := newCustomFieldsReader(t, bytes.NewBuffer(testData))

	runTest(t, cfr, []line{
		newLine("bar", 8, true),
		{
			Foo:     strPtr("bar"),
			Iter:    nil,
			Enabled: boolPtr(true),
		},
		newLine("baz", 9, true),
		newLine("qux", 10, false),
	})
}

func TestOnEmptyInput(t *testing.T) {
	testData := []byte(``)

	runTest(t, newCustomFieldsReader(t, bytes.NewBuffer(testData)), []line{})
}

func TestWithNoCustomFields(t *testing.T) {
	fields, err := custom.ParseFields([]string{})
	require.NoError(t, err)
	cfr := newCustomFieldsReaderWithFields(t, newTestReader(t, testDataLines), fields)

	runTest(t, cfr, []line{})
}

func TestWithNonDefaultFields(t *testing.T) {
	fields, err := custom.ParseFields([]string{
		`foo:string=foo: "(\w+)"`,
		`iter:number`,
		`enabled:bool`,
	})
	require.NoError(t, err)
	cfr := newCustomFieldsReaderWithFields(t, newTestReader(t, testDataLines), fields)

	runTest(t, cfr, testDataLines)
}
