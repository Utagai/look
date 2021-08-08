package execution_test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/utagai/look/query/breeze/execution"
)

var shouldGen = flag.Bool("generate", false, "Whether we should re-generate the comparison tests")

// Not a real test, but some code I used to generate all of the comparison
// tests.
func TestGenerateGoTestCases(t *testing.T) {
	if !*shouldGen {
		t.SkipNow()
	}

	values := []interface{}{
		"foo", "bar",

		"2", "3",
		2, 3,

		3.14, 6.28,
		"3.14", "6.28",

		true, false,
		0, 1,
		"0", "1",

		nil,
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

		// Writing to stderr makes things a bit easier since we can just redirect
		// only stderr and quickly get _just_ the Golang we want to generate.
		fmt.Printf(
			"\n{\n\ta: %s,\n\tb: %s,\n\texpected: %s,\n},",
			toGoStr(tc.a), toGoStr(tc.b), expectedLabel,
		)
	}
}

// toGoStr takes some value and returns a valid, compilable Go string
// representation of it.
func toGoStr(x interface{}) string {
	// Unfortunately, `nil` requires special handling. %#v does the wrong thing.
	if x == nil {
		return "nil"
	}

	return fmt.Sprintf("%#v", x)
}
