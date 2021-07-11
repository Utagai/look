package execution_test

import (
	"fmt"
	"testing"

	"github.com/utagai/look/query/liquid/execution"
)

// Not a real test, but some code I used to generate all of the comparison
// tests.
func TestGenerateGoTestCases(t *testing.T) {
	t.Skip()
	values := []interface{}{
		"foo", "bar",

		"2", "3",
		2, 3,

		3.14, 6.28,
		"3.14", "6.28",

		true, false,
		0, 1,
		"0", "1",
	}

	tcs := make([]testCase, 0, len(values)*len(values))
	for _, a := range values {
		for _, b := range values {
			tcs = append(tcs, testCase{
				a:        a,
				b:        b,
				expected: execution.Compare(a, b),
			})
		}
	}

	for _, tc := range tcs {
		expectedLabel := "execution.Equal"
		switch tc.expected {
		case execution.Equal:
			expectedLabel = "execution.Equal"
		case execution.Lesser:
			expectedLabel = "execution.Lesser"
		case execution.Greater:
			expectedLabel = "execution.Greater"
		default:
			panic(fmt.Sprintf("unrecognized comparison result: %q", tc.expected))
		}

		fmt.Printf("\n{\n\ta: %#v,\n\tb: %#v,\n\texpected: %s,\n},", tc.a, tc.b, expectedLabel)
	}
}
