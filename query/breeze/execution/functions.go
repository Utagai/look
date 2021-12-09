package execution

import (
	"fmt"
	"math"

	"github.com/utagai/look/query/breeze"
)

func executeFunction(function *breeze.Function, args []*breeze.Const) (*breeze.Const, error) {
	// This should always exist, because if it did not, parsing would have failed
	// earlier.
	funcValidator, _ := breeze.LookupFuncValidator(function.Name)

	if err := funcValidator.ValidateTypes(args); err != nil {
		return err.ToEmbeddedDatumErrorMessage(), nil
	}

	switch function.Name {
	case "pow":
		base := args[0].Interface().(float64)
		exp := args[1].Interface().(float64)

		return pow(base, exp), nil
	case "hello":
		return hello(), nil
	}

	return nil, fmt.Errorf("unrecognized function: %q", function.Name)
}

func pow(base float64, exp float64) *breeze.Const {
	return &breeze.Const{
		Kind:        breeze.ConstKindNumber,
		Stringified: fmt.Sprintf("%f", math.Pow(base, exp)),
	}
}

func hello() *breeze.Const {
	return &breeze.Const{
		Kind:        breeze.ConstKindString,
		Stringified: "hello world!",
	}
}
