package execution

import (
	"fmt"
	"math"
	"regexp"

	"github.com/utagai/look/query/breeze"
)

func executeFunction(function *breeze.Function, args []breeze.Concrete) (breeze.Concrete, error) {
	// This should always exist, because if it did not, parsing would have failed
	// earlier.
	funcValidator, _ := breeze.LookupFuncValidator(function.Name)

	if err := funcValidator.ValidateTypes(args); err != nil {
		return err.ToEmbeddedDatumErrorMessage(), nil
	}

	if err := funcValidator.ValidateValues(args); err != nil {
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
	case "regex":
		untypedMatchee, err := args[0].Interface()
		if err != nil {
			return nil, err
		}
		matchee := untypedMatchee.(string)

		untypedPattern, err := args[1].Interface()
		if err != nil {
			return nil, err
		}
		// This will always compile because we've already validated it in
		// the validator for regex.
		pattern := regexp.MustCompile(untypedPattern.(string))

		return regex(matchee, pattern), nil
	case "exists":
		return boolToConcrete(args[0].ConcreteKind() != breeze.ConcreteKindMissing), nil
	case "notexists":
		// HACK: This is likely not how we should be doing
		// this... instead, we should be parsing the unary operator as an
		// expr, and then use ! as its own operator that applies to that
		// expr, inverting its final value.
		return boolToConcrete(args[0].ConcreteKind() == breeze.ConcreteKindMissing), nil
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

func regex(matchee string, pattern *regexp.Regexp) *breeze.Scalar {
	return &breeze.Scalar{
		Kind:        breeze.ScalarKindBool,
		Stringified: fmt.Sprintf("%t", len(pattern.FindString(matchee)) > 0),
	}
}

func hello() *breeze.Scalar {
	return &breeze.Scalar{
		Kind:        breeze.ScalarKindString,
		Stringified: "hello world!",
	}
}
