package breeze_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utagai/look/query/breeze"
)

func TestParser(t *testing.T) {
	type testCase struct {
		query  string
		errMsg string
		stages []breeze.Stage
	}

	tcs := []testCase{
		{
			query: "find foo = 4.2 bar >7 car=\"hello\" | sort bar",
			stages: []breeze.Stage{
				&breeze.Find{
					Checks: []*breeze.Check{
						{
							Field: "foo",
							Value: &breeze.Const{
								Kind:        breeze.ConstKindNumber,
								Stringified: "4.2",
							},
							Op: breeze.BinaryOpEquals,
						},
						{
							Field: "bar",
							Value: &breeze.Const{
								Kind:        breeze.ConstKindNumber,
								Stringified: "7",
							},
							Op: breeze.BinaryOpGeq,
						},
						{
							Field: "car",
							Value: &breeze.Const{
								Kind:        breeze.ConstKindString,
								Stringified: "hello",
							},
							Op: breeze.BinaryOpEquals,
						},
					},
				},
				&breeze.Sort{
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
			parser := breeze.NewParser(tc.query)
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
