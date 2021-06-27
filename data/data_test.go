package data_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utagai/look/data"
	"github.com/utagai/look/datum"
)

const numTestDatums = 100

var getTestDatum = func(i int) datum.Datum {
	return datum.Datum{
		"foo": datum.Datum{
			"a": i,
			"b": i + 1,
		},
		"bar": true,
		"baz": fmt.Sprintf("hello world: %d!", i),
	}
}

var testDatums = func() []datum.Datum {
	datums := make([]datum.Datum, numTestDatums)
	for i := 0; i < numTestDatums; i++ {
		datums[i] = getTestDatum(i)
	}

	return datums
}()

func newMemoryData() data.Data {
	return data.NewMemoryData(testDatums)
}

func TestMemoryDataFindSingle(t *testing.T) {
	ctx := context.Background()
	memData := newMemoryData()

	for i := 0; i < numTestDatums; i++ {
		expectedIndex := i
		foundData, err := memData.Find(ctx, fmt.Sprintf("hello world: %d!", expectedIndex))
		assert.NoError(t, err)
		length, err := foundData.Length(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, length)
		datum, err := foundData.At(ctx, 0)
		assert.NoError(t, err)
		assert.Equal(t, getTestDatum(expectedIndex), datum)
	}
}

func TestMemoryDataLength(t *testing.T) {
	ctx := context.Background()
	memData := newMemoryData()

	initialLength, err := memData.Length(ctx)
	assert.NoError(t, err)
	assert.Equal(t, numTestDatums, initialLength)

	foundData, err := memData.Find(ctx, "hello world: 1")
	assert.NoError(t, err)
	queriedLength, err := foundData.Length(ctx)
	assert.NoError(t, err)
	// 11 is the count of numbers from [1, 100) where the first digit == 1:
	// [1, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19] (note bounds excluding '100').
	assert.Equal(t, 11, queriedLength)
}

func TestMemoryDataAt(t *testing.T) {
	ctx := context.Background()
	memData := newMemoryData()

	for i := 0; i < numTestDatums; i++ {
		datum, err := memData.At(ctx, i)
		assert.NoError(t, err)
		assert.Equal(t, getTestDatum(i), datum)
	}
}
