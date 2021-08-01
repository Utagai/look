package config

// TODO: I don't think this needs to be an internal test.

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData = []byte(`
  value of foo: "bar" on iteration 8 and enabled: true
  value of foo: "baz" on iteration 9 and enabled: true
  value of foo: "qux" on iteration 10 and enabled: false
  value of foo: "bar" on iteration 11 and enabled: true
`)

func TestCustomFieldsReader(t *testing.T) {
	r, err := newCustomFieldsReader(bytes.NewBuffer(testData), []ParseField{
		{
			Type:      FieldTypeString,
			FieldName: "foo",
			Regex:     `"\w+"`,
		},
		{
			Type:      FieldTypeNumber,
			FieldName: "iter",
			Regex:     `\d+`,
		},
		{
			Type:      FieldTypeBool,
			FieldName: "enabled",
			Regex:     `true|false`,
		},
	})
	assert.NoError(t, err)

	actualJSON, err := ioutil.ReadAll(r)
	if err != io.EOF {
		assert.NoError(t, err)
	}

	assert.Equal(t, `[{"foo": "bar","iter": 8,"enabled": true},{"foo": "baz","iter": 9,"enabled": true},{"foo": "qux","iter": 10,"enabled": false},{"foo": "bar","iter": 11,"enabled": true}]`, string(actualJSON))
}

func TestCustomFieldsReaderOnLargeInput(t *testing.T) {
	sizeIncreaseFactor := 100
	largeTestData := make([]byte, 0, len(testData)*sizeIncreaseFactor)
	for i := 0; i < sizeIncreaseFactor; i++ {
		largeTestData = append(largeTestData, testData...)
	}

	r, err := newCustomFieldsReader(bytes.NewBuffer(largeTestData), []ParseField{
		{
			Type:      FieldTypeString,
			FieldName: "foo",
			Regex:     `"\w+"`,
		},
		{
			Type:      FieldTypeNumber,
			FieldName: "iter",
			Regex:     `\d+`,
		},
		{
			Type:      FieldTypeBool,
			FieldName: "enabled",
			Regex:     `true|false`,
		},
	})
	assert.NoError(t, err)

	actualJSON, err := ioutil.ReadAll(r)
	if err != io.EOF {
		assert.NoError(t, err)
	}

	expectedJSON := []byte(`{"foo": "bar","iter": 8,"enabled": true},{"foo": "baz","iter": 9,"enabled": true},{"foo": "qux","iter": 10,"enabled": false},{"foo": "bar","iter": 11,"enabled": true}`)
	largeExpectedJSON := make([]byte, 0, len(expectedJSON)*sizeIncreaseFactor)
	largeExpectedJSON = append(largeExpectedJSON, '[')
	for i := 0; i < sizeIncreaseFactor; i++ {
		largeExpectedJSON = append(largeExpectedJSON, expectedJSON...)
		if i < sizeIncreaseFactor-1 {
			largeExpectedJSON = append(largeExpectedJSON, ',')
		}
	}
	largeExpectedJSON = append(largeExpectedJSON, ']')
	assert.Equal(t, string(largeExpectedJSON), string(actualJSON))
}

func TestNothingMatches(t *testing.T) {
	testData := []byte(`
  foo
  bar
  baz
  `)
	r, err := newCustomFieldsReader(bytes.NewBuffer(testData), []ParseField{
		{
			Type:      FieldTypeString,
			FieldName: "foo",
			Regex:     `"\w+"`,
		},
		{
			Type:      FieldTypeNumber,
			FieldName: "iter",
			Regex:     `\d+`,
		},
		{
			Type:      FieldTypeBool,
			FieldName: "enabled",
			Regex:     `true|false`,
		},
	})
	assert.NoError(t, err)

	actualJSON, err := ioutil.ReadAll(r)
	if err != io.EOF {
		assert.NoError(t, err)
	}

	assert.Equal(t, actualJSON, []byte("[]"))
}
