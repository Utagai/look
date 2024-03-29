package breeze

import "fmt"

// TypeMismatchErr describes a type validation error in finer detail.
type TypeMismatchErr struct {
	ExpectedKind string
	Actual       Concrete
}

// NewTypeMismatchErr creates a new TypeMismatchErr.
func NewTypeMismatchErr(expected string, actual Concrete) *TypeMismatchErr {
	return &TypeMismatchErr{
		ExpectedKind: expected,
		Actual:       actual,
	}
}

// ToEmbeddedDatumErrorMessage returns a string that is suitable for embedding
// into a Datum as a way to communicate error messages about types on a
// per-datum basis.
func (t *TypeMismatchErr) ToEmbeddedDatumErrorMessage() *Scalar {
	kindStr := string(t.Actual.ConcreteKind())
	if scalar, ok := t.Actual.(*Scalar); ok {
		kindStr = string(scalar.Kind)
	}

	errMsg := fmt.Sprintf(
		"[TYPE ERR: expected %s, got '%s' (%s)]",
		t.ExpectedKind,
		t.Actual.GetStringRepr(),
		kindStr,
	)

	return &Scalar{
		Kind:        ScalarKindString,
		Stringified: errMsg,
	}
}
