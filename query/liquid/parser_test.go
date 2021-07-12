package liquid_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utagai/look/query/liquid"
)

func TestParser(t *testing.T) {
	type testCase struct {
		query  string
		errMsg string
		stages []liquid.Stage
	}

	tcs := []testCase{
		{
			query: "find foo = 4.2 bar >7 car=\"hello\" | sort bar",
			stages: []liquid.Stage{
				&liquid.Find{
					Checks: []*liquid.Check{
						{
							Field: "foo",
							Value: &liquid.Const{
								Kind:        liquid.ConstKindNumber,
								Stringified: "4.2",
							},
							Op: liquid.BinaryOpEquals,
						},
						{
							Field: "bar",
							Value: &liquid.Const{
								Kind:        liquid.ConstKindNumber,
								Stringified: "7",
							},
							Op: liquid.BinaryOpGeq,
						},
						{
							Field: "car",
							Value: &liquid.Const{
								Kind:        liquid.ConstKindString,
								Stringified: "hello",
							},
							Op: liquid.BinaryOpEquals,
						},
					},
				},
				&liquid.Sort{
					Descending: false,
					Field:      "bar",
				},
			},
		},
		{
			query:  "find f",
			errMsg: "failed to parse: failed to parse check: failed to parse op: expected a binary operator, but reached end of query",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.query, func(t *testing.T) {
			parser := liquid.NewParser(tc.query)
			stages, err := parser.Parse()
			if tc.errMsg == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.stages, stages)
			} else {
				assert.EqualError(t, err, tc.errMsg)
			}
		})
	}
}
