package execution

import (
	"fmt"
	"math"

	"github.com/utagai/look/query/breeze"
)

func executeFunction(function *breeze.Function, args []breeze.Concrete) (breeze.Concrete, error) {
	// This should always exist, because if it did not, parsing would have failed
	// earlier.
	funcValidator, _ := breeze.LookupFuncValidator(function.Name)

	if err := funcValidator.ValidateTypes(args); err != nil {
		return nil, err
	}

	switch function.Name {
	case "pow":
		untypedBase, err := args[0].Interface()
		if err != nil {
			return nil, err
		}
		base := untypedBase.(float64)
		untypedExp, err := args[1].Interface()
		if err != nil {
			return nil, err
		}
		exp := untypedExp.(float64)

		return pow(base, exp), nil
	case "hello":
		return hello(), nil
	}

	return nil, fmt.Errorf("unrecognized function: %q", function.Name)
}

func pow(base float64, exp float64) *breeze.Scalar {
	return &breeze.Scalar{
		Kind:        breeze.ScalarKindNumber,
		Stringified: fmt.Sprintf("%f", math.Pow(base, exp)),
	}
}

func hello() *breeze.Scalar {
	return &breeze.Scalar{
		Kind:        breeze.ScalarKindString,
		Stringified: "hello world!",
	}
}
