package execution_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
	"github.com/utagai/look/query/breeze/execution"
)

/*
* This package tests execution, but does so in an integration style. It generally
* works by taking an input query, as a string, and passing it through the parser
* and then executing the AST that it creates. In other words, this is not testing
* execution in isolation.
*
* The downside there is of course that failures in the parser may manifest here as
* well. Ideally however, cases like this are easy enough to determine as being a
* parser issue and don't add extra time trying to determine where it is coming
* from. Plus, always starting debugging at the parser layer is a good rule of
* thumb that can mitigate this problem.
*
* The reason for doing things this way is two-fold:
*   * It is literally easier to write tests this way instead of trying to hand-write
*     the AST for each test case.
*   * It is immune to _good_ or _neutral_ changes to the AST/parser behavior causing
*     these tests to be re-written. They should 'just work'. This is the flip-side
*     of the issue described above, where good or neutral changes that should not
*     cause failures or issues can cause development friction.
**/

type executionTestCase struct {
	name              string
	input             []datum.Datum
	query             string
	expectedResult    []datum.Datum
	expectedExecErr   error
	expectedStreamErr error
}

func runExecutionTestCases(t *testing.T, tcs []executionTestCase) {
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			runExecutionTestCase(t, tc)
		})
	}
}

func runExecutionTestCase(t *testing.T, tc executionTestCase) {
	parser := breeze.NewParser(tc.query)
	stages, err := parser.Parse()
	require.NoError(t, err)

	result, err := execution.Execute(datum.NewSliceStream(tc.input), stages)
	if tc.expectedExecErr != nil {
		require.Error(t, err)
		require.ErrorIs(t, err, tc.expectedExecErr)
		return
	}
	require.NoError(t, err)

	actualDatums, err := datum.StreamToSlice(result)
	if tc.expectedStreamErr != nil {
		require.Error(t, err)
		require.ErrorIs(t, err, tc.expectedStreamErr)
		return
	}
	require.NoError(t, err)

	require.Equal(t, tc.expectedResult, actualDatums)
}

func TestFilter(t *testing.T) {
	tcs := []executionTestCase{
		{
			name: "equality filter against scalar number",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
				{
					"a": 3,
				},
			},
			query: "filter a = 2",
			expectedResult: []datum.Datum{
				{
					"a": 2,
				},
			},
		},
		{
			name: "equality filter against scalar bool",
			input: []datum.Datum{
				{
					"a": true,
				},
				{
					"a": true,
				},
				{
					"a": false,
				},
			},
			query: "filter a = false",
			expectedResult: []datum.Datum{
				{
					"a": false,
				},
			},
		},
		{
			name: "equality filter against scalar string",
			input: []datum.Datum{
				{
					"a": "hello",
				},
				{
					"a": "world",
				},
				{
					"a": "!",
				},
			},
			query: "filter a = \"world\"", // Using 'world' fails at parse. Do we want to support '?
			expectedResult: []datum.Datum{
				{
					"a": "world",
				},
			},
		},
		{
			name: "equality filter against scalar null",
			input: []datum.Datum{
				{
					"a": nil,
				},
				{
					"a": "hello",
				},
				{
					"a": "world",
				},
			},
			query: "filter a = null", // Using 'world' fails at parse. Do we want to support '?
			expectedResult: []datum.Datum{
				{
					"a": nil,
				},
			},
		},
		{
			name: "equality filter against array",
			input: []datum.Datum{
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": []int{4, 5, 6},
				},
				{
					"a": []int{7, 8, 9},
				},
			},
			query: "filter a = [4,5,6]", // Using 'world' fails at parse. Do we want to support '?
			expectedResult: []datum.Datum{
				{
					"a": []int{4, 5, 6},
				},
			},
		},
		{
			name: "equality filter against field ref",
			input: []datum.Datum{
				{
					"a": 3,
					"b": 1,
				},
				{
					"a": 4,
					"b": 4,
				},
				{
					"a": 5,
					"b": 7,
				},
			},
			query: "filter a = .b",
			expectedResult: []datum.Datum{
				{
					"a": 4,
					"b": 4,
				},
			},
		},
		{
			name: "equality filter against function evaluation",
			input: []datum.Datum{
				{
					"a": 2,
				},
				{
					"a": 3,
				},
				{
					"a": 4,
				},
			},
			query: "filter a = pow(2,2)",
			expectedResult: []datum.Datum{
				{
					"a": 4,
				},
			},
		},
		{
			name: "filter returning multiple results",
			input: []datum.Datum{
				{
					"a": 2,
				},
				{
					"a": 2,
				},
				{
					"a": 3,
				},
			},
			query: "filter a = 2",
			expectedResult: []datum.Datum{
				{
					"a": 2,
				},
				{
					"a": 2,
				},
			},
		},
		{
			name: "filter returning zero results",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
				{
					"a": 3,
				},
			},
			query:          "filter a = 4",
			expectedResult: []datum.Datum{},
		},
		{
			name: "empty filter returns all",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
				{
					"a": 3,
				},
			},
			query: "filter",
			expectedResult: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
				{
					"a": 3,
				},
			},
		},
		{
			name: "filter on nested array",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": []interface{}{3, "hi there", []int{1, 2, 3}},
				},
				{
					"a": 3,
				},
			},
			query: "filter a = [3, \"hi there\", [1,2,3]]",
			expectedResult: []datum.Datum{
				{
					"a": []interface{}{3, "hi there", []int{1, 2, 3}},
				},
			},
		},
		{
			name:           "filter on empty input",
			input:          []datum.Datum{},
			query:          "filter a = 3",
			expectedResult: []datum.Datum{},
		},
		{
			name: "filter on empty datums",
			input: []datum.Datum{
				{},
				{},
				{},
			},
			query:          "filter a = 3",
			expectedResult: []datum.Datum{},
		},
		{
			name: "filter on empty datums with empty condition",
			input: []datum.Datum{
				{},
				{},
				{},
			},
			query: "filter",
			expectedResult: []datum.Datum{
				{},
				{},
				{},
			},
		},
	}

	runExecutionTestCases(t, tcs)
}

func TestSort(t *testing.T) {
	tcs := []executionTestCase{
		{
			name: "sort across scalar numbers",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 3,
				},
				{
					"a": 2,
				},
			},
			query: "sort a",
			expectedResult: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
				{
					"a": 3,
				},
			},
		},
		{
			name: "sort across scalar strings",
			input: []datum.Datum{
				{
					"a": "def",
				},
				{
					"a": "abc",
				},
				{
					"a": "ghi",
				},
			},
			query: "sort a",
			expectedResult: []datum.Datum{
				{
					"a": "abc",
				},
				{
					"a": "def",
				},
				{
					"a": "ghi",
				},
			},
		},
		{
			name: "sort across scalar strings with lexicographic tie-break",
			input: []datum.Datum{
				{
					"a": "def",
				},
				{
					"a": "abc",
				},
				{
					"a": "abd",
				},
			},
			query: "sort a",
			expectedResult: []datum.Datum{
				{
					"a": "abc",
				},
				{
					"a": "abd",
				},
				{
					"a": "def",
				},
			},
		},
		{
			name: "sort with null always sorts null lowest",
			input: []datum.Datum{
				{
					"a": -3,
				},
				{
					"a": -2,
				},
				{
					"a": nil,
				},
				{
					"a": false,
				},
				{
					"a": "",
				},
			},
			query: "sort a",
			expectedResult: []datum.Datum{
				{
					"a": nil,
				},
				{
					"a": "",
				},
				{
					"a": -3,
				},
				{
					"a": -2,
				},
				{
					"a": false,
				},
			},
		},
		{
			name: "sort across strictly ordered arrays",
			input: []datum.Datum{
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": []int{7, 8, 9},
				},
				{
					"a": []int{4, 5, 6},
				},
			},
			query: "sort a",
			expectedResult: []datum.Datum{
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": []int{4, 5, 6},
				},
				{
					"a": []int{7, 8, 9},
				},
			},
		},
		{
			name: "sort across arrays with lexicographic tie-break",
			input: []datum.Datum{
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": []int{1, 3, 4},
				},
				{
					"a": []int{1, 2, 4},
				},
			},
			query: "sort a",
			expectedResult: []datum.Datum{
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": []int{1, 2, 4},
				},
				{
					"a": []int{1, 3, 4},
				},
			},
		},
		{
			name: "sort across arrays with string elements and double lexicographic tie-break",
			input: []datum.Datum{
				{
					"a": []string{"abc", "def", "ghi"},
				},
				{
					"a": []string{"abc", "def", "ghj"},
				},
				{
					"a": []string{"abd", "def", "ghi"},
				},
			},
			query: "sort a",
			expectedResult: []datum.Datum{
				{
					"a": []string{"abc", "def", "ghi"},
				},
				{
					"a": []string{"abc", "def", "ghj"},
				},
				{
					"a": []string{"abd", "def", "ghi"},
				},
			},
		},
		{
			name: "sort across heterogeneously typed scalars",
			input: []datum.Datum{
				{
					"a": 3,
				},
				{
					"a": "2",
				},
				{
					"a": true,
				},
			},
			query: "sort a",
			expectedResult: []datum.Datum{
				{
					"a": true,
				},
				{
					"a": "2",
				},
				{
					"a": 3,
				},
			},
		},
		{
			name: "sort across arrays and scalars",
			input: []datum.Datum{
				{
					"a": 3,
				},
				{
					"a": []int{3, 2, 1},
				},
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": "0",
				},
			},
			query: "sort a",
			expectedResult: []datum.Datum{
				{
					"a": "0",
				},
				{
					"a": 3,
				},
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": []int{3, 2, 1},
				},
			},
		},
		{
			name:           "sort empty input",
			input:          []datum.Datum{},
			query:          "sort a",
			expectedResult: []datum.Datum{},
		},
		{
			name: "sort with single datum",
			input: []datum.Datum{
				{
					"a": 3,
				},
			},
			query: "sort a",
			expectedResult: []datum.Datum{
				{
					"a": 3,
				},
			},
		},
		{
			name: "sort with equal values",
			input: []datum.Datum{
				{
					"a": 3,
					"b": 1,
				},
				{
					"a": 3,
					"b": 2,
				},
				{
					"a": 3,
					"b": 3,
				},
			},
			query: "sort a",
			expectedResult: []datum.Datum{
				{
					"a": 3,
					"b": 1,
				},
				{
					"a": 3,
					"b": 2,
				},
				{
					"a": 3,
					"b": 3,
				},
			},
		},
	}

	runExecutionTestCases(t, tcs)
}

func TestAggregations(t *testing.T) {
	tcs := []executionTestCase{
		{
			name: "sum with number",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
				{
					"a": 3,
				},
			},
			query: "group sum a",
			expectedResult: []datum.Datum{
				{
					"a": 6.0,
				},
			},
		},
		{
			name: "sum with bool gives true",
			input: []datum.Datum{
				{
					"a": true,
				},
				{
					"a": true,
				},
				{
					"a": false,
				},
			},
			query: "group sum a",
			expectedResult: []datum.Datum{
				{
					"a": true,
				},
			},
		},
		{
			name: "sum with bool gives false",
			input: []datum.Datum{
				{
					"a": false,
				},
				{
					"a": false,
				},
				{
					"a": false,
				},
			},
			query: "group sum a",
			expectedResult: []datum.Datum{
				{
					"a": false,
				},
			},
		},
		{
			name:  "sum no datums",
			input: []datum.Datum{},
			query: "group sum a",
			expectedResult: []datum.Datum{
				{
					"a": 0.0,
				},
			},
		},
		{
			name: "min with number",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
				{
					"a": 3,
				},
			},
			query: "group min a",
			expectedResult: []datum.Datum{
				{
					"a": 1,
				},
			},
		},
		{
			name: "max with number",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
				{
					"a": 3,
				},
			},
			query: "group max a",
			expectedResult: []datum.Datum{
				{
					"a": 3,
				},
			},
		},
		{
			name: "min with bool gives false",
			input: []datum.Datum{
				{
					"a": true,
				},
				{
					"a": true,
				},
				{
					"a": false,
				},
			},
			query: "group min a",
			expectedResult: []datum.Datum{
				{
					"a": false,
				},
			},
		},
		{
			name: "min with only true gives true",
			input: []datum.Datum{
				{
					"a": true,
				},
				{
					"a": true,
				},
				{
					"a": true,
				},
			},
			query: "group min a",
			expectedResult: []datum.Datum{
				{
					"a": true,
				},
			},
		},
		{
			name: "max with bool gives false",
			input: []datum.Datum{
				{
					"a": true,
				},
				{
					"a": true,
				},
				{
					"a": false,
				},
			},
			query: "group max a",
			expectedResult: []datum.Datum{
				{
					"a": true,
				},
			},
		},
		{
			name: "max with only false gives false",
			input: []datum.Datum{
				{
					"a": false,
				},
				{
					"a": false,
				},
				{
					"a": false,
				},
			},
			query: "group max a",
			expectedResult: []datum.Datum{
				{
					"a": false,
				},
			},
		},
		{
			name: "max with only false gives false",
			input: []datum.Datum{
				{
					"a": false,
				},
				{
					"a": false,
				},
				{
					"a": false,
				},
			},
			query: "group max a",
			expectedResult: []datum.Datum{
				{
					"a": false,
				},
			},
		},
		{
			name: "min with all equal",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 1,
				},
				{
					"a": 1,
				},
			},
			query: "group min a",
			expectedResult: []datum.Datum{
				{
					"a": 1,
				},
			},
		},
		{
			name: "max with all equal",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 1,
				},
				{
					"a": 1,
				},
			},
			query: "group max a",
			expectedResult: []datum.Datum{
				{
					"a": 1,
				},
			},
		},
		{
			name:  "max with no datums",
			input: []datum.Datum{},
			query: "group max a",
			expectedResult: []datum.Datum{
				{
					"a": nil,
				},
			},
		},
		{
			name:  "min with no datums",
			input: []datum.Datum{},
			query: "group min a",
			expectedResult: []datum.Datum{
				{
					"a": nil,
				},
			},
		},
		{
			name: "min with string",
			input: []datum.Datum{
				{
					"a": "abc",
				},
				{
					"a": "zzz",
				},
				{
					"a": "def",
				},
			},
			query: "group min a",
			expectedResult: []datum.Datum{
				{
					"a": "abc",
				},
			},
		},
		{
			name: "max with string",
			input: []datum.Datum{
				{
					"a": "abc",
				},
				{
					"a": "zzz",
				},
				{
					"a": "def",
				},
			},
			query: "group max a",
			expectedResult: []datum.Datum{
				{
					"a": "zzz",
				},
			},
		},
		{
			name: "min with null",
			input: []datum.Datum{
				{
					"a": "abc",
				},
				{
					"a": nil,
				},
				{
					"a": "def",
				},
			},
			query: "group min a",
			expectedResult: []datum.Datum{
				{
					"a": nil,
				},
			},
		},
		{
			name: "max with nil",
			input: []datum.Datum{
				{
					"a": 5,
				},
				{
					"a": 3,
				},
				{
					"a": nil,
				},
			},
			query: "group max a",
			expectedResult: []datum.Datum{
				{
					"a": 5,
				},
			},
		},
		{
			name: "max with only nil gives nil",
			input: []datum.Datum{
				{
					"a": nil,
				},
				{
					"a": nil,
				},
				{
					"a": nil,
				},
			},
			query: "group max a",
			expectedResult: []datum.Datum{
				{
					"a": nil,
				},
			},
		},
		{
			name: "min with array",
			input: []datum.Datum{
				{
					"a": []int{4, 5, 6},
				},
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": []int{7, 8, 9},
				},
			},
			query: "group min a",
			expectedResult: []datum.Datum{
				{
					"a": []int{1, 2, 3},
				},
			},
		},
		{
			name: "max with array",
			input: []datum.Datum{
				{
					"a": []int{4, 5, 6},
				},
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": []int{7, 8, 9},
				},
			},
			query: "group max a",
			expectedResult: []datum.Datum{
				{
					"a": []int{7, 8, 9},
				},
			},
		},
		{
			name: "avg with number",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
				{
					"a": 3,
				},
			},
			query: "group avg a",
			expectedResult: []datum.Datum{
				{
					"a": 2.0,
				},
			},
		},
		{
			name: "avg with mixed bool gives number between 0 and 1",
			input: []datum.Datum{
				{
					"a": true,
				},
				{
					"a": true,
				},
				{
					"a": false,
				},
			},
			query: "group avg a",
			expectedResult: []datum.Datum{
				{
					"a": 0.6666666666666666,
				},
			},
		},
		{
			name: "avg with all false gives 0",
			input: []datum.Datum{
				{
					"a": false,
				},
				{
					"a": false,
				},
				{
					"a": false,
				},
			},
			query: "group avg a",
			expectedResult: []datum.Datum{
				{
					"a": 0.0,
				},
			},
		},
		{
			name:  "avg with no datums",
			input: []datum.Datum{},
			query: "group avg a",
			expectedResult: []datum.Datum{
				{
					"a": 0,
				},
			},
		},
		{
			name: "stddev with number",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
				{
					"a": 3,
				},
			},
			query: "group stddev a",
			expectedResult: []datum.Datum{
				{
					"a": 0.816496580927726,
				},
			},
		},
		{
			name: "stddev with mixed bool gives result as stddev of 0s and 1s",
			input: []datum.Datum{
				{
					"a": true,
				},
				{
					"a": true,
				},
				{
					"a": false,
				},
			},
			query: "group stddev a",
			expectedResult: []datum.Datum{
				{
					"a": 0.4714045207910317,
				},
			},
		},
		{
			name: "stddev with all false gives 0",
			input: []datum.Datum{
				{
					"a": false,
				},
				{
					"a": false,
				},
				{
					"a": false,
				},
			},
			query: "group stddev a",
			expectedResult: []datum.Datum{
				{
					"a": 0.0,
				},
			},
		},
		{
			name:  "stddev with no datums gives NaN",
			input: []datum.Datum{},
			query: "group stddev a",
			expectedResult: []datum.Datum{
				{
					"a": "NaN",
				},
			},
		},
		{
			name: "stddev with 1 datum gives NaN",
			input: []datum.Datum{
				{
					"a": 1,
				},
			},
			query: "group stddev a",
			expectedResult: []datum.Datum{
				{
					"a": "NaN",
				},
			},
		},
		{
			name: "stddev with 2 datums does not give NaN",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
			},
			query: "group stddev a",
			expectedResult: []datum.Datum{
				{
					"a": 0.5,
				},
			},
		},
		{
			name: "stddev with identical values gives 0",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 1,
				},
			},
			query: "group stddev a",
			expectedResult: []datum.Datum{
				{
					"a": 0.0,
				},
			},
		},
		{
			name: "count N datums",
			input: []datum.Datum{
				{
					"a": false,
				},
				{
					"a": false,
				},
				{
					"a": false,
				},
			},
			query: "group count a",
			expectedResult: []datum.Datum{
				{
					"a": uint(3),
				},
			},
		},
		{
			name: "count N datums on nonexistent field is 0",
			input: []datum.Datum{
				{
					"a": false,
				},
				{
					"a": false,
				},
				{
					"a": false,
				},
			},
			query: "group count b",
			expectedResult: []datum.Datum{
				{
					"b": uint(0),
				},
			},
		},
		{
			name: "count N datums on sometimes missing field",
			input: []datum.Datum{
				{
					"a": false,
				},
				{
					"b": false,
				},
				{
					"a": false,
				},
			},
			query: "group count b",
			expectedResult: []datum.Datum{
				{
					"b": uint(1),
				},
			},
		},
		{
			name: "count single datum",
			input: []datum.Datum{
				{
					"a": false,
				},
			},
			query: "group count a",
			expectedResult: []datum.Datum{
				{
					"a": uint(1),
				},
			},
		},
		{
			name:  "count no datums",
			input: []datum.Datum{},
			query: "group count a",
			expectedResult: []datum.Datum{
				{
					"a": uint(0),
				},
			},
		},
		{
			name: "count non-numeric type",
			input: []datum.Datum{
				{
					"a": "hello",
				},
				{
					"a": "world",
				},
				{
					"a": "!",
				},
			},
			query: "group count a",
			expectedResult: []datum.Datum{
				{
					"a": uint(3),
				},
			},
		},
		{
			name: "mode simple win",
			input: []datum.Datum{
				{
					"a": 0,
				},
				{
					"a": 2,
				},
				{
					"a": 2,
				},
			},
			query: "group mode a",
			expectedResult: []datum.Datum{
				{
					"a": 2.0,
				},
			},
		},
		{
			name: "mode all datums equal",
			input: []datum.Datum{
				{
					"a": 2,
				},
				{
					"a": 2,
				},
				{
					"a": 2,
				},
			},
			query: "group mode a",
			expectedResult: []datum.Datum{
				{
					"a": 2.0,
				},
			},
		},
		{
			name: "mode only cares about the specified field",
			input: []datum.Datum{
				{
					"a": 1,
					"b": 2,
				},
				{
					"a": 5,
					"b": 2,
				},
				{
					"a": 5,
					"b": 2,
				},
				{
					"a": 3,
					"b": 2,
				},
			},
			query: "group mode a",
			expectedResult: []datum.Datum{
				{
					"a": 5.0,
				},
			},
		},
		{
			name:  "mode no datums",
			input: []datum.Datum{},
			query: "group mode a",
			expectedResult: []datum.Datum{
				{
					"a": nil,
				},
			},
		},
		{
			name: "mode single datum",
			input: []datum.Datum{
				{
					"a": 0,
				},
			},
			query: "group mode a",
			expectedResult: []datum.Datum{
				{
					"a": 0.0,
				},
			},
		},
		{
			name: "mode strings",
			input: []datum.Datum{
				{
					"a": "hello",
				},
				{
					"a": "hello",
				},
				{
					"a": "world",
				},
			},
			query: "group mode a",
			expectedResult: []datum.Datum{
				{
					"a": "hello",
				},
			},
		},
		{
			name: "mode booleans",
			input: []datum.Datum{
				{
					"a": true,
				},
				{
					"a": false,
				},
				{
					"a": true,
				},
			},
			query: "group mode a",
			expectedResult: []datum.Datum{
				{
					"a": true,
				},
			},
		},
		{
			name: "mode nil",
			input: []datum.Datum{
				{
					"a": nil,
				},
				{
					"a": false,
				},
				{
					"a": nil,
				},
			},
			query: "group mode a",
			expectedResult: []datum.Datum{
				{
					"a": nil,
				},
			},
		},
	}

	runExecutionTestCases(t, tcs)
}

func TestAggregationsUnsupportedTypes(t *testing.T) {
	tcs := []executionTestCase{
		{
			name: "sum with string",
			input: []datum.Datum{
				{
					"a": "a",
				},
				{
					"a": "b",
				},
				{
					"a": "c",
				},
			},
			query: "group sum a",
			expectedResult: []datum.Datum{
				{
					"a": 0.0,
				},
			},
		},
		{
			name: "sum with null",
			input: []datum.Datum{
				{
					"a": nil,
				},
				{
					"a": nil,
				},
				{
					"a": nil,
				},
			},
			query: "group sum a",
			expectedResult: []datum.Datum{
				{
					"a": 0.0,
				},
			},
		},
		{
			name: "sum with array",
			input: []datum.Datum{
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": []int{4, 5, 6},
				},
				{
					"a": []int{7, 8, 9},
				},
			},
			query: "group sum a",
			expectedResult: []datum.Datum{
				{
					"a": 0.0,
				},
			},
		},
		{
			name: "sum with mix of supported and unsupported types",
			input: []datum.Datum{
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": 3,
				},
				{
					"a": 4.2,
				},
			},
			query: "group sum a",
			expectedResult: []datum.Datum{
				{
					"a": 7.2,
				},
			},
		},
		{
			name: "avg with string",
			input: []datum.Datum{
				{
					"a": "a",
				},
				{
					"a": "b",
				},
				{
					"a": "c",
				},
			},
			query: "group avg a",
			expectedResult: []datum.Datum{
				{
					"a": 0,
				},
			},
		},
		{
			name: "avg with null",
			input: []datum.Datum{
				{
					"a": nil,
				},
				{
					"a": nil,
				},
				{
					"a": nil,
				},
			},
			query: "group avg a",
			expectedResult: []datum.Datum{
				{
					"a": 0,
				},
			},
		},
		{
			name: "avg with array",
			input: []datum.Datum{
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": []int{4, 5, 6},
				},
				{
					"a": []int{7, 8, 9},
				},
			},
			query: "group avg a",
			expectedResult: []datum.Datum{
				{
					"a": 0,
				},
			},
		},
		{
			name: "avg with mix of supported and unsupported types",
			input: []datum.Datum{
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": 3,
				},
				{
					"a": 4.2,
				},
			},
			query: "group avg a",
			expectedResult: []datum.Datum{
				{
					"a": 3.6,
				},
			},
		},
		{
			name: "stddev with string",
			input: []datum.Datum{
				{
					"a": "a",
				},
				{
					"a": "b",
				},
				{
					"a": "c",
				},
			},
			query: "group stddev a",
			expectedResult: []datum.Datum{
				{
					"a": "NaN",
				},
			},
		},
		{
			name: "stddev with null",
			input: []datum.Datum{
				{
					"a": nil,
				},
				{
					"a": nil,
				},
				{
					"a": nil,
				},
			},
			query: "group stddev a",
			expectedResult: []datum.Datum{
				{
					"a": "NaN",
				},
			},
		},
		{
			name: "stddev with array",
			input: []datum.Datum{
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": []int{4, 5, 6},
				},
				{
					"a": []int{7, 8, 9},
				},
			},
			query: "group stddev a",
			expectedResult: []datum.Datum{
				{
					"a": "NaN",
				},
			},
		},
		{
			name: "stddev with mix of supported and unsupported types",
			input: []datum.Datum{
				{
					"a": []int{1, 2, 3},
				},
				{
					"a": 3,
				},
				{
					"a": 4.2,
				},
			},
			query: "group stddev a",
			expectedResult: []datum.Datum{
				{
					"a": 0.6000000000000001,
				},
			},
		},
	}

	runExecutionTestCases(t, tcs)
}

func TestFunctions(t *testing.T) {
	tcs := []executionTestCase{
		{
			name: "pow successful execution",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
			},
			query: "map res = pow(.a, 2)",
			expectedResult: []datum.Datum{
				{
					"a":   1,
					"res": 1.0,
				},
				{
					"a":   2,
					"res": 4.0,
				},
			},
		},
		{
			name: "pow wrong type",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
			},
			query: `map res = pow(.a, "foo")`,
			expectedResult: []datum.Datum{
				{
					"a":   1,
					"res": "[TYPE ERR: expected number, got 'foo' (string)]",
				},
				{
					"a":   2,
					"res": "[TYPE ERR: expected number, got 'foo' (string)]",
				},
			},
		},
		{
			name: "regex successful execution",
			input: []datum.Datum{
				{
					"a": "hello world",
				},
				{
					"a": "goodbye world",
				},
			},
			query: `map res = regex(.a, "hello")`,
			expectedResult: []datum.Datum{
				{
					"a":   "hello world",
					"res": true,
				},
				{
					"a":   "goodbye world",
					"res": false,
				},
			},
		},
		{
			name: "regex wrong type",
			input: []datum.Datum{
				{
					"a": 1,
				},
				{
					"a": 2,
				},
			},
			query: `map res = regex(.a, "foo")`,
			expectedResult: []datum.Datum{
				{
					"a":   1,
					"res": "[TYPE ERR: expected string, got '1' (number)]",
				},
				{
					"a":   2,
					"res": "[TYPE ERR: expected string, got '2' (number)]",
				},
			},
		},
	}

	runExecutionTestCases(t, tcs)
}
