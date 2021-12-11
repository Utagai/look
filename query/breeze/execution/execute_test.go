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

/*
Cases to think about:
	* sort
	* map
	* project
	* expressions
		- scalars
		- field ref
		- functions
		- combination of above
	* array
		- empty
		- singleton
		- N scalars
		- N function evals
		- N field refs
		- N combos
		- nested arrays
**/

type executionTestCase struct {
	name            string
	input           []datum.Datum
	query           string
	expectedResult  []datum.Datum
	expectedExecErr error
	// TODO: Unclear if we'll actually need this.
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

	result, err := execution.Execute(datum.NewDatumSliceStream(tc.input), stages)
	if tc.expectedExecErr != nil {
		require.ErrorIs(t, err, tc.expectedExecErr)
		return
	}
	require.NoError(t, err)

	actualDatums, err := datum.StreamToSlice(result)
	if tc.expectedStreamErr != nil {
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
