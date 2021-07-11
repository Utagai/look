package execution_test

import (
	"fmt"
	"testing"

	"github.com/utagai/look/query/liquid/execution"
)

func comparisonToSymbol(cmp execution.Comparison) string {
	switch cmp {
	case execution.Lesser:
		return "<"
	case execution.Greater:
		return ">"
	case execution.Equal:
		return "=="
	default:
		panic(fmt.Sprintf("unrecognized Comparison result: %q", cmp))
	}
}

type testCase struct {
	a        interface{}
	b        interface{}
	expected execution.Comparison
}

func (tc testCase) failureMessage(actual execution.Comparison) string {
	return fmt.Sprintf(
		"expected a %s b, but got a %s b",
		comparisonToSymbol(tc.expected),
		comparisonToSymbol(actual),
	)
}

func runTestCases(t *testing.T, tcs []testCase) {
	t.Helper()
	for _, tc := range tcs {
		t.Run(
			fmt.Sprintf("%#v %s %#v", tc.a, tc.expected, tc.b),
			func(t *testing.T) {
				actual := execution.Compare(tc.a, tc.b)
				if actual != tc.expected {
					t.Fatal(tc.failureMessage(actual))
				}
			},
		)
	}
}

func TestComparison(t *testing.T) {
	tcs := []testCase{
		{
			a:        "foo",
			b:        "foo",
			expected: execution.Equal,
		},
		{
			a:        "foo",
			b:        "bar",
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        "2",
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        "3",
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        2,
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        3,
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        3.14,
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        6.28,
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        "3.14",
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        "6.28",
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        true,
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        false,
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        0,
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        1,
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        "0",
			expected: execution.Greater,
		},
		{
			a:        "foo",
			b:        "1",
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        "bar",
			b:        "bar",
			expected: execution.Equal,
		},
		{
			a:        "bar",
			b:        "2",
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        "3",
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        2,
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        3,
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        3.14,
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        6.28,
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        "3.14",
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        "6.28",
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        true,
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        false,
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        0,
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        1,
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        "0",
			expected: execution.Greater,
		},
		{
			a:        "bar",
			b:        "1",
			expected: execution.Greater,
		},
		{
			a:        "2",
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        "2",
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        "2",
			b:        "2",
			expected: execution.Equal,
		},
		{
			a:        "2",
			b:        "3",
			expected: execution.Lesser,
		},
		{
			a:        "2",
			b:        2,
			expected: execution.Equal,
		},
		{
			a:        "2",
			b:        3,
			expected: execution.Lesser,
		},
		{
			a:        "2",
			b:        3.14,
			expected: execution.Lesser,
		},
		{
			a:        "2",
			b:        6.28,
			expected: execution.Lesser,
		},
		{
			a:        "2",
			b:        "3.14",
			expected: execution.Lesser,
		},
		{
			a:        "2",
			b:        "6.28",
			expected: execution.Lesser,
		},
		{
			a:        "2",
			b:        true,
			expected: execution.Greater,
		},
		{
			a:        "2",
			b:        false,
			expected: execution.Greater,
		},
		{
			a:        "2",
			b:        0,
			expected: execution.Greater,
		},
		{
			a:        "2",
			b:        1,
			expected: execution.Greater,
		},
		{
			a:        "2",
			b:        "0",
			expected: execution.Greater,
		},
		{
			a:        "2",
			b:        "1",
			expected: execution.Greater,
		},
		{
			a:        "3",
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        "3",
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        "3",
			b:        "2",
			expected: execution.Greater,
		},
		{
			a:        "3",
			b:        "3",
			expected: execution.Equal,
		},
		{
			a:        "3",
			b:        2,
			expected: execution.Greater,
		},
		{
			a:        "3",
			b:        3,
			expected: execution.Equal,
		},
		{
			a:        "3",
			b:        3.14,
			expected: execution.Lesser,
		},
		{
			a:        "3",
			b:        6.28,
			expected: execution.Lesser,
		},
		{
			a:        "3",
			b:        "3.14",
			expected: execution.Lesser,
		},
		{
			a:        "3",
			b:        "6.28",
			expected: execution.Lesser,
		},
		{
			a:        "3",
			b:        true,
			expected: execution.Greater,
		},
		{
			a:        "3",
			b:        false,
			expected: execution.Greater,
		},
		{
			a:        "3",
			b:        0,
			expected: execution.Greater,
		},
		{
			a:        "3",
			b:        1,
			expected: execution.Greater,
		},
		{
			a:        "3",
			b:        "0",
			expected: execution.Greater,
		},
		{
			a:        "3",
			b:        "1",
			expected: execution.Greater,
		},
		{
			a:        2,
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        2,
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        2,
			b:        "2",
			expected: execution.Equal,
		},
		{
			a:        2,
			b:        "3",
			expected: execution.Lesser,
		},
		{
			a:        2,
			b:        2,
			expected: execution.Equal,
		},
		{
			a:        2,
			b:        3,
			expected: execution.Lesser,
		},
		{
			a:        2,
			b:        3.14,
			expected: execution.Lesser,
		},
		{
			a:        2,
			b:        6.28,
			expected: execution.Lesser,
		},
		{
			a:        2,
			b:        "3.14",
			expected: execution.Lesser,
		},
		{
			a:        2,
			b:        "6.28",
			expected: execution.Lesser,
		},
		{
			a:        2,
			b:        true,
			expected: execution.Greater,
		},
		{
			a:        2,
			b:        false,
			expected: execution.Greater,
		},
		{
			a:        2,
			b:        0,
			expected: execution.Greater,
		},
		{
			a:        2,
			b:        1,
			expected: execution.Greater,
		},
		{
			a:        2,
			b:        "0",
			expected: execution.Greater,
		},
		{
			a:        2,
			b:        "1",
			expected: execution.Greater,
		},
		{
			a:        3,
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        3,
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        3,
			b:        "2",
			expected: execution.Greater,
		},
		{
			a:        3,
			b:        "3",
			expected: execution.Equal,
		},
		{
			a:        3,
			b:        2,
			expected: execution.Greater,
		},
		{
			a:        3,
			b:        3,
			expected: execution.Equal,
		},
		{
			a:        3,
			b:        3.14,
			expected: execution.Lesser,
		},
		{
			a:        3,
			b:        6.28,
			expected: execution.Lesser,
		},
		{
			a:        3,
			b:        "3.14",
			expected: execution.Lesser,
		},
		{
			a:        3,
			b:        "6.28",
			expected: execution.Lesser,
		},
		{
			a:        3,
			b:        true,
			expected: execution.Greater,
		},
		{
			a:        3,
			b:        false,
			expected: execution.Greater,
		},
		{
			a:        3,
			b:        0,
			expected: execution.Greater,
		},
		{
			a:        3,
			b:        1,
			expected: execution.Greater,
		},
		{
			a:        3,
			b:        "0",
			expected: execution.Greater,
		},
		{
			a:        3,
			b:        "1",
			expected: execution.Greater,
		},
		{
			a:        3.14,
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        3.14,
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        3.14,
			b:        "2",
			expected: execution.Greater,
		},
		{
			a:        3.14,
			b:        "3",
			expected: execution.Greater,
		},
		{
			a:        3.14,
			b:        2,
			expected: execution.Greater,
		},
		{
			a:        3.14,
			b:        3,
			expected: execution.Greater,
		},
		{
			a:        3.14,
			b:        3.14,
			expected: execution.Equal,
		},
		{
			a:        3.14,
			b:        6.28,
			expected: execution.Lesser,
		},
		{
			a:        3.14,
			b:        "3.14",
			expected: execution.Equal,
		},
		{
			a:        3.14,
			b:        "6.28",
			expected: execution.Lesser,
		},
		{
			a:        3.14,
			b:        true,
			expected: execution.Greater,
		},
		{
			a:        3.14,
			b:        false,
			expected: execution.Greater,
		},
		{
			a:        3.14,
			b:        0,
			expected: execution.Greater,
		},
		{
			a:        3.14,
			b:        1,
			expected: execution.Greater,
		},
		{
			a:        3.14,
			b:        "0",
			expected: execution.Greater,
		},
		{
			a:        3.14,
			b:        "1",
			expected: execution.Greater,
		},
		{
			a:        6.28,
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        6.28,
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        6.28,
			b:        "2",
			expected: execution.Greater,
		},
		{
			a:        6.28,
			b:        "3",
			expected: execution.Greater,
		},
		{
			a:        6.28,
			b:        2,
			expected: execution.Greater,
		},
		{
			a:        6.28,
			b:        3,
			expected: execution.Greater,
		},
		{
			a:        6.28,
			b:        3.14,
			expected: execution.Greater,
		},
		{
			a:        6.28,
			b:        6.28,
			expected: execution.Equal,
		},
		{
			a:        6.28,
			b:        "3.14",
			expected: execution.Greater,
		},
		{
			a:        6.28,
			b:        "6.28",
			expected: execution.Equal,
		},
		{
			a:        6.28,
			b:        true,
			expected: execution.Greater,
		},
		{
			a:        6.28,
			b:        false,
			expected: execution.Greater,
		},
		{
			a:        6.28,
			b:        0,
			expected: execution.Greater,
		},
		{
			a:        6.28,
			b:        1,
			expected: execution.Greater,
		},
		{
			a:        6.28,
			b:        "0",
			expected: execution.Greater,
		},
		{
			a:        6.28,
			b:        "1",
			expected: execution.Greater,
		},
		{
			a:        "3.14",
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        "3.14",
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        "3.14",
			b:        "2",
			expected: execution.Greater,
		},
		{
			a:        "3.14",
			b:        "3",
			expected: execution.Greater,
		},
		{
			a:        "3.14",
			b:        2,
			expected: execution.Greater,
		},
		{
			a:        "3.14",
			b:        3,
			expected: execution.Greater,
		},
		{
			a:        "3.14",
			b:        3.14,
			expected: execution.Equal,
		},
		{
			a:        "3.14",
			b:        6.28,
			expected: execution.Lesser,
		},
		{
			a:        "3.14",
			b:        "3.14",
			expected: execution.Equal,
		},
		{
			a:        "3.14",
			b:        "6.28",
			expected: execution.Lesser,
		},
		{
			a:        "3.14",
			b:        true,
			expected: execution.Greater,
		},
		{
			a:        "3.14",
			b:        false,
			expected: execution.Greater,
		},
		{
			a:        "3.14",
			b:        0,
			expected: execution.Greater,
		},
		{
			a:        "3.14",
			b:        1,
			expected: execution.Greater,
		},
		{
			a:        "3.14",
			b:        "0",
			expected: execution.Greater,
		},
		{
			a:        "3.14",
			b:        "1",
			expected: execution.Greater,
		},
		{
			a:        "6.28",
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        "6.28",
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        "6.28",
			b:        "2",
			expected: execution.Greater,
		},
		{
			a:        "6.28",
			b:        "3",
			expected: execution.Greater,
		},
		{
			a:        "6.28",
			b:        2,
			expected: execution.Greater,
		},
		{
			a:        "6.28",
			b:        3,
			expected: execution.Greater,
		},
		{
			a:        "6.28",
			b:        3.14,
			expected: execution.Greater,
		},
		{
			a:        "6.28",
			b:        6.28,
			expected: execution.Equal,
		},
		{
			a:        "6.28",
			b:        "3.14",
			expected: execution.Greater,
		},
		{
			a:        "6.28",
			b:        "6.28",
			expected: execution.Equal,
		},
		{
			a:        "6.28",
			b:        true,
			expected: execution.Greater,
		},
		{
			a:        "6.28",
			b:        false,
			expected: execution.Greater,
		},
		{
			a:        "6.28",
			b:        0,
			expected: execution.Greater,
		},
		{
			a:        "6.28",
			b:        1,
			expected: execution.Greater,
		},
		{
			a:        "6.28",
			b:        "0",
			expected: execution.Greater,
		},
		{
			a:        "6.28",
			b:        "1",
			expected: execution.Greater,
		},
		{
			a:        true,
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        true,
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        true,
			b:        "2",
			expected: execution.Lesser,
		},
		{
			a:        true,
			b:        "3",
			expected: execution.Lesser,
		},
		{
			a:        true,
			b:        2,
			expected: execution.Lesser,
		},
		{
			a:        true,
			b:        3,
			expected: execution.Lesser,
		},
		{
			a:        true,
			b:        3.14,
			expected: execution.Lesser,
		},
		{
			a:        true,
			b:        6.28,
			expected: execution.Lesser,
		},
		{
			a:        true,
			b:        "3.14",
			expected: execution.Lesser,
		},
		{
			a:        true,
			b:        "6.28",
			expected: execution.Lesser,
		},
		{
			a:        true,
			b:        true,
			expected: execution.Equal,
		},
		{
			a:        true,
			b:        false,
			expected: execution.Greater,
		},
		{
			a:        true,
			b:        0,
			expected: execution.Greater,
		},
		{
			a:        true,
			b:        1,
			expected: execution.Equal,
		},
		{
			a:        true,
			b:        "0",
			expected: execution.Greater,
		},
		{
			a:        true,
			b:        "1",
			expected: execution.Equal,
		},
		{
			a:        false,
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        false,
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        false,
			b:        "2",
			expected: execution.Lesser,
		},
		{
			a:        false,
			b:        "3",
			expected: execution.Lesser,
		},
		{
			a:        false,
			b:        2,
			expected: execution.Lesser,
		},
		{
			a:        false,
			b:        3,
			expected: execution.Lesser,
		},
		{
			a:        false,
			b:        3.14,
			expected: execution.Lesser,
		},
		{
			a:        false,
			b:        6.28,
			expected: execution.Lesser,
		},
		{
			a:        false,
			b:        "3.14",
			expected: execution.Lesser,
		},
		{
			a:        false,
			b:        "6.28",
			expected: execution.Lesser,
		},
		{
			a:        false,
			b:        true,
			expected: execution.Lesser,
		},
		{
			a:        false,
			b:        false,
			expected: execution.Equal,
		},
		{
			a:        false,
			b:        0,
			expected: execution.Equal,
		},
		{
			a:        false,
			b:        1,
			expected: execution.Lesser,
		},
		{
			a:        false,
			b:        "0",
			expected: execution.Equal,
		},
		{
			a:        false,
			b:        "1",
			expected: execution.Lesser,
		},
		{
			a:        0,
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        0,
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        0,
			b:        "2",
			expected: execution.Lesser,
		},
		{
			a:        0,
			b:        "3",
			expected: execution.Lesser,
		},
		{
			a:        0,
			b:        2,
			expected: execution.Lesser,
		},
		{
			a:        0,
			b:        3,
			expected: execution.Lesser,
		},
		{
			a:        0,
			b:        3.14,
			expected: execution.Lesser,
		},
		{
			a:        0,
			b:        6.28,
			expected: execution.Lesser,
		},
		{
			a:        0,
			b:        "3.14",
			expected: execution.Lesser,
		},
		{
			a:        0,
			b:        "6.28",
			expected: execution.Lesser,
		},
		{
			a:        0,
			b:        true,
			expected: execution.Lesser,
		},
		{
			a:        0,
			b:        false,
			expected: execution.Equal,
		},
		{
			a:        0,
			b:        0,
			expected: execution.Equal,
		},
		{
			a:        0,
			b:        1,
			expected: execution.Lesser,
		},
		{
			a:        0,
			b:        "0",
			expected: execution.Equal,
		},
		{
			a:        0,
			b:        "1",
			expected: execution.Lesser,
		},
		{
			a:        1,
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        1,
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        1,
			b:        "2",
			expected: execution.Lesser,
		},
		{
			a:        1,
			b:        "3",
			expected: execution.Lesser,
		},
		{
			a:        1,
			b:        2,
			expected: execution.Lesser,
		},
		{
			a:        1,
			b:        3,
			expected: execution.Lesser,
		},
		{
			a:        1,
			b:        3.14,
			expected: execution.Lesser,
		},
		{
			a:        1,
			b:        6.28,
			expected: execution.Lesser,
		},
		{
			a:        1,
			b:        "3.14",
			expected: execution.Lesser,
		},
		{
			a:        1,
			b:        "6.28",
			expected: execution.Lesser,
		},
		{
			a:        1,
			b:        true,
			expected: execution.Equal,
		},
		{
			a:        1,
			b:        false,
			expected: execution.Greater,
		},
		{
			a:        1,
			b:        0,
			expected: execution.Greater,
		},
		{
			a:        1,
			b:        1,
			expected: execution.Equal,
		},
		{
			a:        1,
			b:        "0",
			expected: execution.Greater,
		},
		{
			a:        1,
			b:        "1",
			expected: execution.Equal,
		},
		{
			a:        "0",
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        "0",
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        "0",
			b:        "2",
			expected: execution.Lesser,
		},
		{
			a:        "0",
			b:        "3",
			expected: execution.Lesser,
		},
		{
			a:        "0",
			b:        2,
			expected: execution.Lesser,
		},
		{
			a:        "0",
			b:        3,
			expected: execution.Lesser,
		},
		{
			a:        "0",
			b:        3.14,
			expected: execution.Lesser,
		},
		{
			a:        "0",
			b:        6.28,
			expected: execution.Lesser,
		},
		{
			a:        "0",
			b:        "3.14",
			expected: execution.Lesser,
		},
		{
			a:        "0",
			b:        "6.28",
			expected: execution.Lesser,
		},
		{
			a:        "0",
			b:        true,
			expected: execution.Lesser,
		},
		{
			a:        "0",
			b:        false,
			expected: execution.Equal,
		},
		{
			a:        "0",
			b:        0,
			expected: execution.Equal,
		},
		{
			a:        "0",
			b:        1,
			expected: execution.Lesser,
		},
		{
			a:        "0",
			b:        "0",
			expected: execution.Equal,
		},
		{
			a:        "0",
			b:        "1",
			expected: execution.Lesser,
		},
		{
			a:        "1",
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        "1",
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        "1",
			b:        "2",
			expected: execution.Lesser,
		},
		{
			a:        "1",
			b:        "3",
			expected: execution.Lesser,
		},
		{
			a:        "1",
			b:        2,
			expected: execution.Lesser,
		},
		{
			a:        "1",
			b:        3,
			expected: execution.Lesser,
		},
		{
			a:        "1",
			b:        3.14,
			expected: execution.Lesser,
		},
		{
			a:        "1",
			b:        6.28,
			expected: execution.Lesser,
		},
		{
			a:        "1",
			b:        "3.14",
			expected: execution.Lesser,
		},
		{
			a:        "1",
			b:        "6.28",
			expected: execution.Lesser,
		},
		{
			a:        "1",
			b:        true,
			expected: execution.Equal,
		},
		{
			a:        "1",
			b:        false,
			expected: execution.Greater,
		},
		{
			a:        "1",
			b:        0,
			expected: execution.Greater,
		},
		{
			a:        "1",
			b:        1,
			expected: execution.Equal,
		},
		{
			a:        "1",
			b:        "0",
			expected: execution.Greater,
		},
		{
			a:        "1",
			b:        "1",
			expected: execution.Equal,
		},
	}

	runTestCases(t, tcs)
}
