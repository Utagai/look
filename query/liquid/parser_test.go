package liquid_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utagai/look/query/liquid"
)

func TestParser(t *testing.T) {
	parser := liquid.NewParser("find foo = 4.2 bar >7 car=\"hello\" | sort bar")
	stages, err := parser.Parse()
	assert.NoError(t, err)
	assert.Equal(t,
		[]liquid.Stage{
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
		stages,
	)
}
