package data_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utagai/look/data"
	"github.com/utagai/look/datum"
	"github.com/utagai/look/query"
)

const (
	numTestDatums  = 100
	testMongoDBURI = "mongodb://localhost:27017"
)

var getTestDatum = func(i int) datum.Datum {
	return datum.Datum{
		"foo": datum.Datum{
			"a": int64(i),
			"b": int64(i) + 1,
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

// Generally, we'd prefer returning the struct instead of an interface, but we
// can't create a slice of these constructors unless we make them both return
// the same type (the shared interface), so here we are.
func newMemoryData(_ *testing.T) data.Data {
	return data.NewMemoryData(testDatums, query.NewSubstringQueryExecutor())
}

func newMongoDBData(t *testing.T) data.Data {
	mdbData, err := data.NewMongoDBData(testMongoDBURI, "data_test", t.Name(), testDatums)
	assert.NoError(t, err)
	return mdbData
}

type dataTestCase struct {
	constructor func(*testing.T) data.Data
	// Given a number identifying a test datum, return a query that fetches it.
	getSingleQuery func(int) string
	// In the data's query language, give a query that returns all datums where
	// the baz field has 'hello world: 1' as a prefix.
	multipleQuery string
}

func TestData(t *testing.T) {
	dataTestCases := map[string]dataTestCase{
		"memory": {
			constructor:    newMemoryData,
			getSingleQuery: func(i int) string { return fmt.Sprintf("hello world: %d!", i) },
			multipleQuery:  "hello world: 1",
		},
		"mongodb": func() dataTestCase {
			return dataTestCase{
				constructor: newMongoDBData,
				getSingleQuery: func(i int) string {
					return fmt.Sprintf(`[{"$match": {"baz": "hello world: %d!"}}]`, i)
				},
				multipleQuery: `[{"$match": {"baz": {"$regex": "hello world: 1"}}}]`,
			}
		}(),
	}

	// We don't have an explicit test for .Length(), but these tests inevitably
	// test it.
	tests := map[string]func(*testing.T, data.Data, dataTestCase){
		"single find":   testDataFindSingle,
		"multiple find": testDataFindMultiple,
		"at":            testDataAt,
	}

	for dataName, dataTestCase := range dataTestCases {
		t.Run(dataName, func(t *testing.T) {
			dataConstructor := dataTestCase.constructor
			for testName, test := range tests {
				t.Run(testName, func(t *testing.T) {
					test(t, dataConstructor(t), dataTestCase)
				})
			}
		})
	}
}

func testDataFindSingle(t *testing.T, data data.Data, tc dataTestCase) {
	ctx := context.Background()

	for i := 0; i < numTestDatums; i++ {
		expectedIndex := i
		foundData, err := data.Find(ctx, tc.getSingleQuery(i))
		assert.NoError(t, err)
		length, err := foundData.Length(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, length)
		datum, err := foundData.At(ctx, 0)
		assert.NoError(t, err)
		assert.Equal(t, getTestDatum(expectedIndex), datum)
	}
}

func testDataFindMultiple(t *testing.T, data data.Data, tc dataTestCase) {
	ctx := context.Background()

	initialLength, err := data.Length(ctx)
	assert.NoError(t, err)
	assert.Equal(t, numTestDatums, initialLength)

	foundData, err := data.Find(ctx, tc.multipleQuery)
	assert.NoError(t, err)
	queriedLength, err := foundData.Length(ctx)
	assert.NoError(t, err)
	// 11 is the count of numbers from [1, 100) where the first digit == 1:
	// [1, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19] (note bounds excluding '100').
	assert.Equal(t, 11, queriedLength)
}

func testDataAt(t *testing.T, data data.Data, _ dataTestCase) {
	ctx := context.Background()

	for i := 0; i < numTestDatums; i++ {
		datum, err := data.At(ctx, i)
		assert.NoError(t, err)
		assert.Equal(t, getTestDatum(i), datum)
	}
}
