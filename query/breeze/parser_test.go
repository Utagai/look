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
			query: "filter foo = 4.2 bar >7 car=\"hello\" | sort bar",
			stages: []breeze.Stage{
				&breeze.Filter{
					BinaryChecks: []*breeze.BinaryCheck{
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
					UnaryChecks: []*breeze.UnaryCheck{},
				},
				&breeze.Sort{
					Descending: false,
					Field:      "bar",
				},
			},
		},
		{
			query: "filter foo = 4.2",
			stages: []breeze.Stage{
				&breeze.Filter{
					BinaryChecks: []*breeze.BinaryCheck{
						{
							Field: "foo",
							Value: &breeze.Const{
								Kind:        breeze.ConstKindNumber,
								Stringified: "4.2",
							},
							Op: breeze.BinaryOpEquals,
						},
					},
					UnaryChecks: []*breeze.UnaryCheck{},
				},
			},
		},
		{
			query: "filter foo| = 4.2",
			stages: []breeze.Stage{
				&breeze.Filter{
					BinaryChecks: []*breeze.BinaryCheck{
						{
							Field: "foo|",
							Value: &breeze.Const{
								Kind:        breeze.ConstKindNumber,
								Stringified: "4.2",
							},
							Op: breeze.BinaryOpEquals,
						},
					},
					UnaryChecks: []*breeze.UnaryCheck{},
				},
			},
		},
		{
			query: "filter foo = 4.2 bar = 6.4",
			stages: []breeze.Stage{
				&breeze.Filter{
					BinaryChecks: []*breeze.BinaryCheck{
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
								Stringified: "6.4",
							},
							Op: breeze.BinaryOpEquals,
						},
					},
					UnaryChecks: []*breeze.UnaryCheck{},
				},
			},
		},
		{
			query: "sort foo",
			stages: []breeze.Stage{
				&breeze.Sort{
					Descending: false,
					Field:      "foo",
				},
			},
		},
		{
			query: "sort foo desc",
			stages: []breeze.Stage{
				&breeze.Sort{
					Descending: true,
					Field:      "foo",
				},
			},
		},
		{
			query: "sort foo asc",
			stages: []breeze.Stage{
				&breeze.Sort{
					Descending: false,
					Field:      "foo",
				},
			},
		},
		{
			query: "sort foo | filter foo = 4.2",
			stages: []breeze.Stage{
				&breeze.Sort{
					Descending: false,
					Field:      "foo",
				},
				&breeze.Filter{
					BinaryChecks: []*breeze.BinaryCheck{
						{
							Field: "foo",
							Value: &breeze.Const{
								Kind:        breeze.ConstKindNumber,
								Stringified: "4.2",
							},
							Op: breeze.BinaryOpEquals,
						},
					},
					UnaryChecks: []*breeze.UnaryCheck{},
				},
			},
		},
		{
			query: "sort foo asc | filter foo = 4.2",
			stages: []breeze.Stage{
				&breeze.Sort{
					Descending: false,
					Field:      "foo",
				},
				&breeze.Filter{
					BinaryChecks: []*breeze.BinaryCheck{
						{
							Field: "foo",
							Value: &breeze.Const{
								Kind:        breeze.ConstKindNumber,
								Stringified: "4.2",
							},
							Op: breeze.BinaryOpEquals,
						},
					},
					UnaryChecks: []*breeze.UnaryCheck{},
				},
			},
		},
		{
			query: "filter",
			stages: []breeze.Stage{
				&breeze.Filter{
					BinaryChecks: []*breeze.BinaryCheck{},
					UnaryChecks:  []*breeze.UnaryCheck{},
				},
			},
		},
		{
			query: "filter foo exists",
			stages: []breeze.Stage{
				&breeze.Filter{
					BinaryChecks: []*breeze.BinaryCheck{},
					UnaryChecks: []*breeze.UnaryCheck{
						{
							Field: "foo",
							Op:    breeze.UnaryOpExists,
						},
					},
				},
			},
		},
		{
			query: "filter foo !exists",
			stages: []breeze.Stage{
				&breeze.Filter{
					BinaryChecks: []*breeze.BinaryCheck{},
					UnaryChecks: []*breeze.UnaryCheck{
						{
							Field: "foo",
							Op:    breeze.UnaryOpExistsNot,
						},
					},
				},
			},
		},
		{
			query: "filter foo exists bar = 4.2 baz !exists quux = 4.1",
			stages: []breeze.Stage{
				&breeze.Filter{
					BinaryChecks: []*breeze.BinaryCheck{
						{
							Field: "bar",
							Value: &breeze.Const{
								Kind:        breeze.ConstKindNumber,
								Stringified: "4.2",
							},
							Op: breeze.BinaryOpEquals,
						},
						{
							Field: "quux",
							Value: &breeze.Const{
								Kind:        breeze.ConstKindNumber,
								Stringified: "4.1",
							},
							Op: breeze.BinaryOpEquals,
						},
					},
					UnaryChecks: []*breeze.UnaryCheck{
						{
							Field: "foo",
							Op:    breeze.UnaryOpExists,
						},
						{
							Field: "baz",
							Op:    breeze.UnaryOpExistsNot,
						},
					},
				},
			},
		},
		{
			query: "filter foo = 4.2 foo = 4.1",
			stages: []breeze.Stage{
				&breeze.Filter{
					BinaryChecks: []*breeze.BinaryCheck{
						{
							Field: "foo",
							Value: &breeze.Const{
								Kind:        breeze.ConstKindNumber,
								Stringified: "4.2",
							},
							Op: breeze.BinaryOpEquals,
						},
						{
							Field: "foo",
							Value: &breeze.Const{
								Kind:        breeze.ConstKindNumber,
								Stringified: "4.1",
							},
							Op: breeze.BinaryOpEquals,
						},
					},
					UnaryChecks: []*breeze.UnaryCheck{},
				},
			},
		},
		{
			query:  "",
			stages: []breeze.Stage{},
		},
		{
			query:  "filter f",
			errMsg: "failed to parse: failed to parse check: failed to parse op: expected an operator, but reached end of query",
		},
		{
			query:  "filter foo = ",
			errMsg: "failed to parse: failed to parse check: failed to parse constant value: expected a constant value, but reached end of query",
		},
		{
			query:  "sort",
			errMsg: "failed to parse: failed to parse field: expected a field, but reached end of query",
		},
		{
			query:  "|",
			errMsg: "failed to parse: unrecognized stage: \"|\"",
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
