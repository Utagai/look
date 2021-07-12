package execution

import "fmt"

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
	if cmp := compareInterfaceToInterface(a, b); cmp != nil {
		return *cmp
	}
	return Equal
}

func compareInterfaceToInterface(a, b interface{}) *Comparison {
	if maybeNumber := convertPotentialNumber(a); maybeNumber != nil {
		return compareNumberToInterface(*maybeNumber, b)
	} else if maybeStr := convertPotentialString(a); maybeStr != nil {
		return compareStringToInterface(*maybeStr, b)
	} else if maybeBool := convertPotentialBool(a); maybeBool != nil {
		return compareBoolToInterface(*maybeBool, b)
	}

	return nil
}

/*
   ===========
   Comparisons to interface.
   ===========
*/
func compareNumberToInterface(a float64, b interface{}) *Comparison {
	var cmp Comparison
	if maybeNumber := convertPotentialNumber(b); maybeNumber != nil {
		cmp = compareNumbers(a, *maybeNumber)
	} else if maybeStr := convertPotentialString(b); maybeStr != nil {
		cmp = compareNumberAndString(a, *maybeStr)
	} else if maybeBool := convertPotentialBool(b); maybeBool != nil {
		cmp = compareNumberAndBool(a, *maybeBool)
	} else {
		return nil
	}

	return &cmp
}

func compareStringToInterface(a string, b interface{}) *Comparison {
	var cmp Comparison
	if maybeNumber := convertPotentialNumber(b); maybeNumber != nil {
		cmp = compareStringAndNumber(a, *maybeNumber)
	} else if maybeStr := convertPotentialString(b); maybeStr != nil {
		cmp = compareStrings(a, *maybeStr)
	} else if maybeBool := convertPotentialBool(b); maybeBool != nil {
		cmp = compareStringAndBool(a, *maybeBool)
	} else {
		return nil
	}

	return &cmp
}

func compareBoolToInterface(a bool, b interface{}) *Comparison {
	var cmp Comparison
	if maybeBool := convertPotentialBool(b); maybeBool != nil {
		cmp = compareBools(a, *maybeBool)
	} else if maybeStr := convertPotentialString(b); maybeStr != nil {
		cmp = compareBoolAndString(a, *maybeStr)
	} else if maybeNumber := convertPotentialNumber(b); maybeNumber != nil {
		cmp = compareBoolAndNumber(a, *maybeNumber)
	} else {
		return nil
	}

	return &cmp
}

/*
   ===========
   Conversions.
   ===========
*/
func convertPotentialNumber(a interface{}) *float64 {
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
		return nil
	}

	return &af64
}

func convertPotentialString(a interface{}) *string {
	if s, ok := a.(string); ok {
		return &s
	}

	return nil
}

func convertPotentialBool(a interface{}) *bool {
	if b, ok := a.(bool); ok {
		return &b
	}

	return nil
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
