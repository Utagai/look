package breeze_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/utagai/look/query/breeze"
)

// Would be good to split this gigantic test into pieces for e.g. each stage.
func TestParser(t *testing.T) {
	type testCase struct {
		query  string
		errMsg string
		stages []breeze.Stage
	}

	tcs := []testCase{
		{
			query: "filter foo = 4.2 bar >7 car=\"hello\" dar=null | sort bar",
			stages: []breeze.Stage{
				&breeze.Filter{
					BinaryChecks: []*breeze.BinaryCheck{
						{
							Field: "foo",
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindNumber,
								Stringified: "4.2",
							},
							Op: breeze.BinaryOpEquals,
						},
						{
							Field: "bar",
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindNumber,
								Stringified: "7",
							},
							Op: breeze.BinaryOpGeq,
						},
						{
							Field: "car",
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindString,
								Stringified: "hello",
							},
							Op: breeze.BinaryOpEquals,
						},
						{
							Field: "dar",
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindNull,
								Stringified: "null",
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
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindNumber,
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
			query: "filter foo = 8 * 4.2",
			stages: []breeze.Stage{
				&breeze.Filter{
					BinaryChecks: []*breeze.BinaryCheck{
						{
							Field: "foo",
							Expr: &breeze.BinaryExpr{
								Left: &breeze.Scalar{
									Kind:        breeze.ScalarKindNumber,
									Stringified: "8",
								},
								Op: breeze.BinaryOpMultiply,
								Right: &breeze.Scalar{
									Kind:        breeze.ScalarKindNumber,
									Stringified: "4.2",
								},
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
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindNumber,
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
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindNumber,
								Stringified: "4.2",
							},
							Op: breeze.BinaryOpEquals,
						},
						{
							Field: "bar",
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindNumber,
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
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindNumber,
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
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindNumber,
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
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindNumber,
								Stringified: "4.2",
							},
							Op: breeze.BinaryOpEquals,
						},
						{
							Field: "quux",
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindNumber,
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
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindNumber,
								Stringified: "4.2",
							},
							Op: breeze.BinaryOpEquals,
						},
						{
							Field: "foo",
							Expr: &breeze.Scalar{
								Kind:        breeze.ScalarKindNumber,
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
			query: "map foo = [1,2,3]",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "foo",
							Assignment: breeze.Array{
								&breeze.Scalar{
									Kind:        breeze.ScalarKindNumber,
									Stringified: "1",
								},
								&breeze.Scalar{
									Kind:        breeze.ScalarKindNumber,
									Stringified: "2",
								},
								&breeze.Scalar{
									Kind:        breeze.ScalarKindNumber,
									Stringified: "3",
								},
							},
						},
					},
				},
			},
		},
		{
			query: "map foo = [1]",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "foo",
							Assignment: breeze.Array{
								&breeze.Scalar{
									Kind:        breeze.ScalarKindNumber,
									Stringified: "1",
								},
							},
						},
					},
				},
			},
		},
		{
			query: "map foo = []",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field:      "foo",
							Assignment: breeze.Array{},
						},
					},
				},
			},
		},
		// TODO: These tests below are for maps. We should take some time at some
		// point to further flesh these out e.g. with more cases, combinations, etc.
		{
			query: "map foo = 4.2",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "foo",
							Assignment: &breeze.Scalar{
								Kind:        breeze.ScalarKindNumber,
								Stringified: "4.2",
							},
						},
					},
				},
			},
		},
		{
			query: "map foo = .jelly",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "foo",
							Assignment: &breeze.FieldRef{
								Field: "jelly",
							},
						},
					},
				},
			},
		},
		{
			query: "map foo = pow(2, 3)",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "foo",
							Assignment: &breeze.Function{
								Name: "pow",
								Args: []breeze.Expr{
									&breeze.Scalar{
										Kind:        breeze.ScalarKindNumber,
										Stringified: "2",
									},
									&breeze.Scalar{
										Kind:        breeze.ScalarKindNumber,
										Stringified: "3",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			query: "map foo = pow(.bar, pow(3, \"2\"))", // Invalid function call, but not at parse time.
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "foo",
							Assignment: &breeze.Function{
								Name: "pow",
								Args: []breeze.Expr{
									&breeze.FieldRef{
										Field: "bar",
									},
									&breeze.Function{
										Name: "pow",
										Args: []breeze.Expr{
											&breeze.Scalar{
												Kind:        breeze.ScalarKindNumber,
												Stringified: "3",
											},
											&breeze.Scalar{
												Kind:        breeze.ScalarKindString,
												Stringified: "2",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			query: "map foo = 2 + 5",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "foo",
							Assignment: &breeze.BinaryExpr{
								Left: &breeze.Scalar{
									Kind:        breeze.ScalarKindNumber,
									Stringified: "2",
								},
								Op: breeze.BinaryOpPlus,
								Right: &breeze.Scalar{
									Kind:        breeze.ScalarKindNumber,
									Stringified: "5",
								},
							},
						},
					},
				},
			},
		},
		{
			query: "map foo = 2 + 5 - 3",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "foo",
							Assignment: &breeze.BinaryExpr{
								Left: &breeze.Scalar{
									Kind:        breeze.ScalarKindNumber,
									Stringified: "2",
								},
								Op: breeze.BinaryOpPlus,
								Right: &breeze.BinaryExpr{
									Left: &breeze.Scalar{
										Kind:        breeze.ScalarKindNumber,
										Stringified: "5",
									},
									Op: breeze.BinaryOpMinus,
									Right: &breeze.Scalar{
										Kind:        breeze.ScalarKindNumber,
										Stringified: "3",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			query: "map foo = 2 + 5 * 3 - 4",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "foo",
							Assignment: &breeze.BinaryExpr{
								Left: &breeze.Scalar{Kind: "number", Stringified: "2"},
								Right: &breeze.BinaryExpr{
									Left: &breeze.BinaryExpr{
										Left:  &breeze.Scalar{Kind: "number", Stringified: "5"},
										Right: &breeze.Scalar{Kind: "number", Stringified: "3"},
										Op:    "*",
									},
									Right: &breeze.Scalar{Kind: "number", Stringified: "4"},
									Op:    "-",
								},
								Op: "+",
							},
						},
					},
				},
			},
		},
		{
			query: "map foo = 4 * 3 + 1 + 2 * 8 / 9 + 4 - 1",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "foo",
							Assignment: &breeze.BinaryExpr{
								Left: &breeze.BinaryExpr{
									Left:  &breeze.Scalar{Kind: "number", Stringified: "4"},
									Right: &breeze.Scalar{Kind: "number", Stringified: "3"},
									Op:    "*",
								},
								Right: &breeze.BinaryExpr{
									Left: &breeze.Scalar{Kind: "number", Stringified: "1"},
									Right: &breeze.BinaryExpr{
										Left: &breeze.BinaryExpr{
											Left: &breeze.Scalar{Kind: "number", Stringified: "2"},
											Right: &breeze.BinaryExpr{
												Left:  &breeze.Scalar{Kind: "number", Stringified: "8"},
												Right: &breeze.Scalar{Kind: "number", Stringified: "9"},
												Op:    "/",
											},
											Op: "*",
										},
										Right: &breeze.BinaryExpr{
											Left:  &breeze.Scalar{Kind: "number", Stringified: "4"},
											Right: &breeze.Scalar{Kind: "number", Stringified: "1"},
											Op:    "-",
										},
										Op: "+",
									},
									Op: "+",
								},
								Op: "+",
							},
						},
					},
				},
			},
		},
		{
			query: "map foo = 3 * (5 + 2)",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "foo",
							Assignment: &breeze.BinaryExpr{
								Left: &breeze.Scalar{
									Kind:        breeze.ScalarKindNumber,
									Stringified: "3",
								},
								Op: breeze.BinaryOpMultiply,
								Right: &breeze.BinaryExpr{
									Left: &breeze.Scalar{
										Kind:        breeze.ScalarKindNumber,
										Stringified: "5",
									},
									Op: breeze.BinaryOpPlus,
									Right: &breeze.Scalar{
										Kind:        breeze.ScalarKindNumber,
										Stringified: "2",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			query: "map foo = (5 + 2) * 3",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "foo",
							Assignment: &breeze.BinaryExpr{
								Left: &breeze.BinaryExpr{
									Left: &breeze.Scalar{
										Kind:        breeze.ScalarKindNumber,
										Stringified: "5",
									},
									Op: breeze.BinaryOpPlus,
									Right: &breeze.Scalar{
										Kind:        breeze.ScalarKindNumber,
										Stringified: "2",
									},
								},
								Op: breeze.BinaryOpMultiply,
								Right: &breeze.Scalar{
									Kind:        breeze.ScalarKindNumber,
									Stringified: "3",
								},
							},
						},
					},
				},
			},
		},
		{
			query: "map foo = pow(2,2) + .jelly",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "foo",
							Assignment: &breeze.BinaryExpr{
								Left: &breeze.Function{
									Name: "pow",
									Args: []breeze.Expr{
										&breeze.Scalar{
											Kind:        breeze.ScalarKindNumber,
											Stringified: "2",
										},
										&breeze.Scalar{
											Kind:        breeze.ScalarKindNumber,
											Stringified: "2",
										},
									},
								},
								Op: breeze.BinaryOpPlus,
								Right: &breeze.FieldRef{
									Field: "jelly",
								},
							},
						},
					},
				},
			},
		},
		{
			query: "map pythag = pow(pow(.a,2) + pow(.b,2), 0.5)",
			stages: []breeze.Stage{
				&breeze.Map{
					Assignments: []breeze.FieldAssignment{
						{
							Field: "pythag",
							Assignment: &breeze.Function{
								Name: "pow",
								Args: []breeze.Expr{
									&breeze.BinaryExpr{
										Left: &breeze.Function{
											Name: "pow",
											Args: []breeze.Expr{
												&breeze.FieldRef{Field: "a"},
												&breeze.Scalar{Kind: "number", Stringified: "2"},
											},
										},
										Right: &breeze.Function{
											Name: "pow",
											Args: []breeze.Expr{
												&breeze.FieldRef{Field: "b"},
												&breeze.Scalar{Kind: "number", Stringified: "2"},
											},
										},
										Op: "+",
									},
									&breeze.Scalar{Kind: "number", Stringified: "0.5"},
								},
							},
						},
					},
				},
			},
		},
		// TODO: Yea, we obviously need to be better about the error message here.
		{
			query:  "map foo = 3 * (5 + 2 LOL",
			errMsg: "failed to parse: failed to parse assignment: expected a closing paranthesis, but got \"LOL\"",
		},
		{
			query:  "map foo = 4.2 bar = jelly()",
			errMsg: "failed to parse: failed to parse assignment: failed to parse value in expr: failed to parse a value; expected a constant value (expected a constant value, but got: \"jelly\"), field reference (field references must start with '.'), function (unrecognized function: \"jelly\"), or array (expected array to start with '[', but found \"jelly\")",
		},
		{
			query:  "map foo = 4.2 bar = pow()",
			errMsg: "failed to parse: failed to parse assignment: failed to parse value in expr: failed to parse a value; expected a constant value (expected a constant value, but got: \"pow\"), field reference (field references must start with '.'), function (expected 2 args, got 0), or array (expected array to start with '[', but found \")\")",
		},
		{
			query:  "map foo = 4.2 bar = pow(",
			errMsg: "failed to parse: failed to parse assignment: failed to parse value in expr: failed to parse a value; expected a constant value (expected a constant value, but got: \"pow\"), field reference (field references must start with '.'), function (expected 2 args, got 0), or array (expected array to start with '[', but found \"\")",
		},
		{
			query:  "map foo = 4.2 bar = ishouldhaveadotatbeginning",
			errMsg: "failed to parse: failed to parse assignment: failed to parse value in expr: failed to parse a value; expected a constant value (expected a constant value, but got: \"ishouldhaveadotatbeginning\"), field reference (field references must start with '.'), function (unrecognized function: \"ishouldhaveadotatbeginning\"), or array (expected array to start with '[', but found \"ishouldhaveadotatbeginning\")",
		},
		// TODO: END OF MAP TESTS
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
			errMsg: "failed to parse: failed to parse check: failed to parse constant value: failed to parse value in expr: expected a value, but reached end of query",
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
				// pretty.Println(stages)
				assert.Equal(t, tc.stages, stages)
			} else {
				assert.EqualError(t, err, tc.errMsg)
			}
		})
	}
}
