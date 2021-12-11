package execution_test

// THESE TESTS ARE GENERATED.
// YOU HAVE BEEN WARNED.

import (
	"fmt"
	"testing"

	"github.com/utagai/look/query/breeze/execution"
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

type cmpTestCase struct {
	a        interface{}
	b        interface{}
	expected execution.Comparison
}

func (tc cmpTestCase) failureMessage(actual execution.Comparison) string {
	return fmt.Sprintf(
		"expected a %s b, but got a %s b",
		comparisonToSymbol(tc.expected),
		comparisonToSymbol(actual),
	)
}

func runCmpTestCases(t *testing.T, tcs []cmpTestCase) {
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

// THESE TESTS ARE GENERATED.
// YOU HAVE BEEN WARNED.
func TestComparison(t *testing.T) {
	tcs := []cmpTestCase{
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
			a:        "foo",
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        "foo",
			b:        nil,
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
			a:        "bar",
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        "bar",
			b:        nil,
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
			a:        "2",
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        "2",
			b:        nil,
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
			a:        "3",
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        "3",
			b:        nil,
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
			a:        2,
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        2,
			b:        nil,
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
			a:        3,
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        3,
			b:        nil,
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
			a:        3.14,
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        3.14,
			b:        nil,
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
			a:        6.28,
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        6.28,
			b:        nil,
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
			a:        "3.14",
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        "3.14",
			b:        nil,
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
			a:        "6.28",
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        "6.28",
			b:        nil,
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
			a:        true,
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        true,
			b:        nil,
			expected: execution.Greater,
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
			a:        false,
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        false,
			b:        nil,
			expected: execution.Greater,
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
			a:        0,
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        0,
			b:        nil,
			expected: execution.Greater,
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
			a:        1,
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        1,
			b:        nil,
			expected: execution.Greater,
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
			a:        "0",
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        "0",
			b:        nil,
			expected: execution.Greater,
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
		{
			a:        "1",
			b:        []int{},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []int{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []int{1, 2, 3, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []int{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []int{2, 3},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []float64{},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []float64{1, 1, 2},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []float64{7, 8, 9},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []float64{3.14, 6.28},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []string{},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []string{"foo", "bar"},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []string{"1", "2", "3"},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Incomparable,
		},
		{
			a:        "1",
			b:        nil,
			expected: execution.Greater,
		},
		{
			a:        []int{},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []int{},
			b:        []int{},
			expected: execution.Equal,
		},
		{
			a:        []int{},
			b:        []int{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []int{},
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        []int{},
			b:        []int{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []int{},
			b:        []int{2, 3},
			expected: execution.Lesser,
		},
		{
			a:        []int{},
			b:        []float64{},
			expected: execution.Equal,
		},
		{
			a:        []int{},
			b:        []float64{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []int{},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []int{},
			b:        []float64{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []int{},
			b:        []float64{3.14, 6.28},
			expected: execution.Lesser,
		},
		{
			a:        []int{},
			b:        []string{},
			expected: execution.Equal,
		},
		{
			a:        []int{},
			b:        []string{"foo", "bar"},
			expected: execution.Lesser,
		},
		{
			a:        []int{},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Lesser,
		},
		{
			a:        []int{},
			b:        []string{"1", "2", "3"},
			expected: execution.Lesser,
		},
		{
			a:        []int{},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Lesser,
		},
		{
			a:        []int{},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 1, 2},
			b:        []int{},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 1, 2},
			b:        []int{1, 1, 2},
			expected: execution.Equal,
		},
		{
			a:        []int{1, 1, 2},
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        []int{1, 1, 2},
			b:        []int{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []int{1, 1, 2},
			b:        []int{2, 3},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 1, 2},
			b:        []float64{},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 1, 2},
			b:        []float64{1, 1, 2},
			expected: execution.Equal,
		},
		{
			a:        []int{1, 1, 2},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []int{1, 1, 2},
			b:        []float64{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []int{1, 1, 2},
			b:        []float64{3.14, 6.28},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 1, 2},
			b:        []string{},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 1, 2},
			b:        []string{"foo", "bar"},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 1, 2},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Lesser,
		},
		{
			a:        []int{1, 1, 2},
			b:        []string{"1", "2", "3"},
			expected: execution.Lesser,
		},
		{
			a:        []int{1, 1, 2},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 1, 2},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []int{},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []int{1, 1, 2},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []int{1, 2, 3, 4},
			expected: execution.Equal,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []int{7, 8, 9},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []int{2, 3},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []float64{},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []float64{1, 1, 2},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []float64{7, 8, 9},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []float64{3.14, 6.28},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []string{},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []string{"foo", "bar"},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []string{"1", "2", "3"},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Greater,
		},
		{
			a:        []int{1, 2, 3, 4},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []int{7, 8, 9},
			b:        []int{},
			expected: execution.Greater,
		},
		{
			a:        []int{7, 8, 9},
			b:        []int{1, 1, 2},
			expected: execution.Greater,
		},
		{
			a:        []int{7, 8, 9},
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        []int{7, 8, 9},
			b:        []int{7, 8, 9},
			expected: execution.Equal,
		},
		{
			a:        []int{7, 8, 9},
			b:        []int{2, 3},
			expected: execution.Greater,
		},
		{
			a:        []int{7, 8, 9},
			b:        []float64{},
			expected: execution.Greater,
		},
		{
			a:        []int{7, 8, 9},
			b:        []float64{1, 1, 2},
			expected: execution.Greater,
		},
		{
			a:        []int{7, 8, 9},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []int{7, 8, 9},
			b:        []float64{7, 8, 9},
			expected: execution.Equal,
		},
		{
			a:        []int{7, 8, 9},
			b:        []float64{3.14, 6.28},
			expected: execution.Greater,
		},
		{
			a:        []int{7, 8, 9},
			b:        []string{},
			expected: execution.Greater,
		},
		{
			a:        []int{7, 8, 9},
			b:        []string{"foo", "bar"},
			expected: execution.Greater,
		},
		{
			a:        []int{7, 8, 9},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Lesser,
		},
		{
			a:        []int{7, 8, 9},
			b:        []string{"1", "2", "3"},
			expected: execution.Greater,
		},
		{
			a:        []int{7, 8, 9},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Greater,
		},
		{
			a:        []int{7, 8, 9},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []int{2, 3},
			b:        []int{},
			expected: execution.Greater,
		},
		{
			a:        []int{2, 3},
			b:        []int{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []int{2, 3},
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        []int{2, 3},
			b:        []int{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []int{2, 3},
			b:        []int{2, 3},
			expected: execution.Equal,
		},
		{
			a:        []int{2, 3},
			b:        []float64{},
			expected: execution.Greater,
		},
		{
			a:        []int{2, 3},
			b:        []float64{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []int{2, 3},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []int{2, 3},
			b:        []float64{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []int{2, 3},
			b:        []float64{3.14, 6.28},
			expected: execution.Lesser,
		},
		{
			a:        []int{2, 3},
			b:        []string{},
			expected: execution.Greater,
		},
		{
			a:        []int{2, 3},
			b:        []string{"foo", "bar"},
			expected: execution.Lesser,
		},
		{
			a:        []int{2, 3},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Lesser,
		},
		{
			a:        []int{2, 3},
			b:        []string{"1", "2", "3"},
			expected: execution.Lesser,
		},
		{
			a:        []int{2, 3},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Lesser,
		},
		{
			a:        []int{2, 3},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{},
			b:        []int{},
			expected: execution.Equal,
		},
		{
			a:        []float64{},
			b:        []int{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []float64{},
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        []float64{},
			b:        []int{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []float64{},
			b:        []int{2, 3},
			expected: execution.Lesser,
		},
		{
			a:        []float64{},
			b:        []float64{},
			expected: execution.Equal,
		},
		{
			a:        []float64{},
			b:        []float64{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []float64{},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []float64{},
			b:        []float64{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []float64{},
			b:        []float64{3.14, 6.28},
			expected: execution.Lesser,
		},
		{
			a:        []float64{},
			b:        []string{},
			expected: execution.Equal,
		},
		{
			a:        []float64{},
			b:        []string{"foo", "bar"},
			expected: execution.Lesser,
		},
		{
			a:        []float64{},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Lesser,
		},
		{
			a:        []float64{},
			b:        []string{"1", "2", "3"},
			expected: execution.Lesser,
		},
		{
			a:        []float64{},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Lesser,
		},
		{
			a:        []float64{},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []int{},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []int{1, 1, 2},
			expected: execution.Equal,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []int{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []int{2, 3},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []float64{},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []float64{1, 1, 2},
			expected: execution.Equal,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []float64{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []float64{3.14, 6.28},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []string{},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []string{"foo", "bar"},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Lesser,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []string{"1", "2", "3"},
			expected: execution.Lesser,
		},
		{
			a:        []float64{1, 1, 2},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 1, 2},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []int{},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []int{1, 1, 2},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []int{1, 2, 3, 4},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []int{7, 8, 9},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []int{2, 3},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []float64{},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []float64{1, 1, 2},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Equal,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []float64{7, 8, 9},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []float64{3.14, 6.28},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []string{},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []string{"foo", "bar"},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []string{"1", "2", "3"},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Greater,
		},
		{
			a:        []float64{1, 2, 3.14, 4},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []int{},
			expected: execution.Greater,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []int{1, 1, 2},
			expected: execution.Greater,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []int{7, 8, 9},
			expected: execution.Equal,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []int{2, 3},
			expected: execution.Greater,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []float64{},
			expected: execution.Greater,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []float64{1, 1, 2},
			expected: execution.Greater,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []float64{7, 8, 9},
			expected: execution.Equal,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []float64{3.14, 6.28},
			expected: execution.Greater,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []string{},
			expected: execution.Greater,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []string{"foo", "bar"},
			expected: execution.Greater,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Lesser,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []string{"1", "2", "3"},
			expected: execution.Greater,
		},
		{
			a:        []float64{7, 8, 9},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Greater,
		},
		{
			a:        []float64{7, 8, 9},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []int{},
			expected: execution.Greater,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []int{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []int{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []int{2, 3},
			expected: execution.Greater,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []float64{},
			expected: execution.Greater,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []float64{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []float64{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []float64{3.14, 6.28},
			expected: execution.Equal,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []string{},
			expected: execution.Greater,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []string{"foo", "bar"},
			expected: execution.Lesser,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Lesser,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []string{"1", "2", "3"},
			expected: execution.Lesser,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Lesser,
		},
		{
			a:        []float64{3.14, 6.28},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []string{},
			b:        []int{},
			expected: execution.Equal,
		},
		{
			a:        []string{},
			b:        []int{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []string{},
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        []string{},
			b:        []int{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []string{},
			b:        []int{2, 3},
			expected: execution.Lesser,
		},
		{
			a:        []string{},
			b:        []float64{},
			expected: execution.Equal,
		},
		{
			a:        []string{},
			b:        []float64{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []string{},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []string{},
			b:        []float64{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []string{},
			b:        []float64{3.14, 6.28},
			expected: execution.Lesser,
		},
		{
			a:        []string{},
			b:        []string{},
			expected: execution.Equal,
		},
		{
			a:        []string{},
			b:        []string{"foo", "bar"},
			expected: execution.Lesser,
		},
		{
			a:        []string{},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Lesser,
		},
		{
			a:        []string{},
			b:        []string{"1", "2", "3"},
			expected: execution.Lesser,
		},
		{
			a:        []string{},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Lesser,
		},
		{
			a:        []string{},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []int{},
			expected: execution.Greater,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []int{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []int{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []int{2, 3},
			expected: execution.Greater,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []float64{},
			expected: execution.Greater,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []float64{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []float64{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []float64{3.14, 6.28},
			expected: execution.Greater,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []string{},
			expected: execution.Greater,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []string{"foo", "bar"},
			expected: execution.Equal,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Lesser,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []string{"1", "2", "3"},
			expected: execution.Lesser,
		},
		{
			a:        []string{"foo", "bar"},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Lesser,
		},
		{
			a:        []string{"foo", "bar"},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []int{},
			expected: execution.Greater,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []int{1, 1, 2},
			expected: execution.Greater,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []int{7, 8, 9},
			expected: execution.Greater,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []int{2, 3},
			expected: execution.Greater,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []float64{},
			expected: execution.Greater,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []float64{1, 1, 2},
			expected: execution.Greater,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []float64{7, 8, 9},
			expected: execution.Greater,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []float64{3.14, 6.28},
			expected: execution.Greater,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []string{},
			expected: execution.Greater,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []string{"foo", "bar"},
			expected: execution.Greater,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Equal,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []string{"1", "2", "3"},
			expected: execution.Greater,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Greater,
		},
		{
			a:        []string{"baz", "qux", "alpha"},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []int{},
			expected: execution.Greater,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []int{1, 1, 2},
			expected: execution.Greater,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []int{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []int{2, 3},
			expected: execution.Greater,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []float64{},
			expected: execution.Greater,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []float64{1, 1, 2},
			expected: execution.Greater,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []float64{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []float64{3.14, 6.28},
			expected: execution.Greater,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []string{},
			expected: execution.Greater,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []string{"foo", "bar"},
			expected: execution.Greater,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Lesser,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []string{"1", "2", "3"},
			expected: execution.Equal,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Greater,
		},
		{
			a:        []string{"1", "2", "3"},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        "foo",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        "bar",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        "2",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        "3",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        2,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        3,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        3.14,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        6.28,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        "3.14",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        "6.28",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        true,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        false,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        0,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        1,
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        "0",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        "1",
			expected: execution.Incomparable,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []int{},
			expected: execution.Greater,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []int{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []int{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []int{2, 3},
			expected: execution.Greater,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []float64{},
			expected: execution.Greater,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []float64{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []float64{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []float64{3.14, 6.28},
			expected: execution.Greater,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []string{},
			expected: execution.Greater,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []string{"foo", "bar"},
			expected: execution.Greater,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Lesser,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []string{"1", "2", "3"},
			expected: execution.Lesser,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Equal,
		},
		{
			a:        []string{"0", "3.14", "6.28"},
			b:        nil,
			expected: execution.Incomparable,
		},
		{
			a:        nil,
			b:        "foo",
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        "bar",
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        "2",
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        "3",
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        2,
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        3,
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        3.14,
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        6.28,
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        "3.14",
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        "6.28",
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        true,
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        false,
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        0,
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        1,
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        "0",
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        "1",
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []int{},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []int{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []int{1, 2, 3, 4},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []int{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []int{2, 3},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []float64{},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []float64{1, 1, 2},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []float64{1, 2, 3.14, 4},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []float64{7, 8, 9},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []float64{3.14, 6.28},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []string{},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []string{"foo", "bar"},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []string{"baz", "qux", "alpha"},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []string{"1", "2", "3"},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        []string{"0", "3.14", "6.28"},
			expected: execution.Lesser,
		},
		{
			a:        nil,
			b:        nil,
			expected: execution.Equal,
		},
	}
	runCmpTestCases(t, tcs)
}
