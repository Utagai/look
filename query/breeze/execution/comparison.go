package execution

import (
	"fmt"
)

// This module implements comparison logic for untyped values which are the
// kinds of inputs that pass through the language implementation.

// Comparison represents the result of a comparison between two const types.
type Comparison string

// The three cases of comparison.
const (
	Lesser  = "less than"
	Greater = "greater than"
	Equal   = "equal"
)

// Compare takes two interfaces that represent some const (see breeze.Const)
// value, compares them, and returns a Comparison result.
// For non-const types, if both types are non-const, they are Equal. If only one
// is non-Const, then the non-Const is lesser.
// There is only a single numeric 'type', Number, and it is a float64.
// Pointers to const types are not const.
// When differing const kinds are compared, they are casted to allow comparison.
// The type hierarchy is string <- number <- bool.
func Compare(a, b interface{}) Comparison {
	return compareInterfaceToInterface(a, b)
}

func compareInterfaceToInterface(a, b interface{}) Comparison {
	if num, ok := convertPotentialNumber(a); ok {
		return compareNumberToInterface(num, b)
	} else if str, ok := convertPotentialString(a); ok {
		return compareStringToInterface(str, b)
	} else if bool, ok := convertPotentialBool(a); ok {
		return compareBoolToInterface(bool, b)
	} else if null, ok := convertPotentialNull(a); ok {
		return compareNullToInterface(null, b)
	}

	panic(fmt.Sprintf("failed to cast to known type (was %T)", a))
}

/*
   ===========
   Comparisons to interface.
   ===========
*/
func compareNumberToInterface(a float64, b interface{}) Comparison {
	if num, ok := convertPotentialNumber(b); ok {
		return compareNumbers(a, num)
	} else if str, ok := convertPotentialString(b); ok {
		return compareNumberAndString(a, str)
	} else if bool, ok := convertPotentialBool(b); ok {
		return compareNumberAndBool(a, bool)
	} else if null, ok := convertPotentialNull(b); ok {
		return compareInterfaceToNull(a, null)
	}

	panic(fmt.Sprintf("failed to cast to known type (was %T)", b))
}

func compareStringToInterface(a string, b interface{}) Comparison {
	if num, ok := convertPotentialNumber(b); ok {
		return compareStringAndNumber(a, num)
	} else if str, ok := convertPotentialString(b); ok {
		return compareStrings(a, str)
	} else if bool, ok := convertPotentialBool(b); ok {
		return compareStringAndBool(a, bool)
	} else if null, ok := convertPotentialNull(b); ok {
		return compareInterfaceToNull(a, null)
	}

	panic(fmt.Sprintf("failed to cast to known type (was %T)", b))
}

func compareBoolToInterface(a bool, b interface{}) Comparison {
	if bool, ok := convertPotentialBool(b); ok {
		return compareBools(a, bool)
	} else if str, ok := convertPotentialString(b); ok {
		return compareBoolAndString(a, str)
	} else if num, ok := convertPotentialNumber(b); ok {
		return compareBoolAndNumber(a, num)
	} else if null, ok := convertPotentialNull(b); ok {
		return compareInterfaceToNull(a, null)
	}

	panic(fmt.Sprintf("failed to cast to known type (was %T)", b))
}

func compareNullToInterface(a interface{}, b interface{}) Comparison {
	// Nulls are simple. They are only ever equal to b if b is itself null,
	// otherwise, they are always lower in ordering.
	_, ok := convertPotentialNull(b)
	var cmp Comparison = Lesser
	if ok {
		cmp = Equal
	}

	return cmp
}

func compareInterfaceToNull(a interface{}, b interface{}) Comparison {
	cmp := compareNullToInterface(b, a)
	if cmp == Lesser {
		cmp = Greater
	}

	return cmp
}

/*
   ===========
   Conversions.
   ===========
*/
func convertPotentialNumber(a interface{}) (float64, bool) {
	// If integer, convert to float64 and compare.
	var af64 float64
	switch ta := a.(type) {
	case int:
		af64 = float64(ta)
	case int16:
		af64 = float64(ta)
	case int32:
		af64 = float64(ta)
	case int64:
		af64 = float64(ta)
	case uint:
		af64 = float64(ta)
	case uint16:
		af64 = float64(ta)
	case uint32:
		af64 = float64(ta)
	case uint64:
		af64 = float64(ta)
	case float32:
		af64 = float64(ta)
	case float64:
		af64 = ta
	default:
		return 0, false
	}

	return af64, true
}

func convertPotentialString(a interface{}) (string, bool) {
	if s, ok := a.(string); ok {
		return s, true
	}

	return "", false
}

func convertPotentialBool(a interface{}) (bool, bool) {
	if b, ok := a.(bool); ok {
		return b, true
	}

	return false, false
}

func convertPotentialNull(a interface{}) (interface{}, bool) {
	if a == nil {
		return a, true
	}

	return nil, false
}

/*
   ===========
   Homogeneous type comparisons.
   ===========
*/
func compareNumbers(a, b float64) Comparison {
	switch {
	case a == b:
		return Equal
	case a < b:
		return Lesser
	default:
		return Greater
	}
}

func compareStrings(a, b string) Comparison {
	switch {
	case a == b:
		return Equal
	case a < b:
		return Lesser
	default:
		return Greater
	}
}

func compareBools(a, b bool) Comparison {
	switch {
	case a == b:
		return Equal
	case a: // Implies that a = True and b = False.
		return Greater
	default: // Implies that a = False and b = True.
		return Lesser
	}
}

/*
   ===========
   Heterogenous type comparisons.
   ===========
*/
func compareStringAndNumber(a string, b float64) Comparison {
	return compareStrings(a, fmt.Sprintf("%v", b))
}

func compareStringAndBool(a string, b bool) Comparison {
	return compareStringAndNumber(a, boolToNumber(b))
}

func compareNumberAndString(a float64, b string) Comparison {
	// I could call compareStringAndNumber() here, but it's literally more work.
	return compareStrings(fmt.Sprintf("%v", a), b)
}

func compareNumberAndBool(a float64, b bool) Comparison {
	return compareNumbers(a, boolToNumber(b))
}

func compareBoolAndNumber(a bool, b float64) Comparison {
	return compareNumbers(boolToNumber(a), b)
}

func compareBoolAndString(a bool, b string) Comparison {
	return compareNumberAndString(boolToNumber(a), b)
}

/*
   ===========
   Helpers.
   ===========
*/
func boolToNumber(a bool) float64 {
	if a {
		return 1
	}

	return 0
}
