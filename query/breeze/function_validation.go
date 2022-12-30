package breeze

import (
	"fmt"
	"regexp"
)

// FunctionValidator validates the use of a function.
type FunctionValidator interface {
	// ExpectedNumArgs returns the expected number of arguments.
	ExpectedNumArgs() int
	// ValidateTypes determines if the given arguments are of valid types.
	// It assumes that the members of the given slice of *Const are non-nil.
	// This function is called _after_ ExpectedNumArgs(), so it can rely
	// on the number of args passed to it to conform to the returned
	// result of ExpectedNumArgs().
	ValidateTypes(args []Concrete) *TypeMismatchErr
	// ValidateValues runs validations after types have been confirmed
	// to be correct. This means that this function runs after both
	// ExpectedNumArgs() and ValidateTypes(). The validation this
	// function performs can be literally anything, and is dependent
	// entirely on the function implementation. For some functions, this
	// may do nothing.
	ValidateValues(args []Concrete) error
}

type helloValidator struct{}

func (p *helloValidator) ExpectedNumArgs() int {
	return 0
}

func (p *helloValidator) ValidateTypes(args []Concrete) *TypeMismatchErr {
	return nil
}

func (p *helloValidator) ValidateValues(args []Concrete) error {
	return nil
}

// pow() is expected to be used as follows:
//   pow(base: <number>, exp: <number>)
type powValidator struct{}

func (p *powValidator) ExpectedNumArgs() int {
	return 2
}

func (p *powValidator) ValidateTypes(args []Concrete) *TypeMismatchErr {
	for _, arg := range args {
		if !(arg.ConcreteKind() == ConcreteKindScalar && arg.(*Scalar).Kind == ScalarKindNumber) {
			return &TypeMismatchErr{
				ExpectedKind: string(ScalarKindNumber),
				Actual:       arg,
			}
		}
	}

	return nil
}

func (p *powValidator) ValidateValues(args []Concrete) error {
	return nil
}

// regex() is expected to be used as follows:
//	regex(str: <string>, pattern: <string>)
type regexValidator struct{}

func (r *regexValidator) ExpectedNumArgs() int {
	return 2
}

func (r *regexValidator) ValidateTypes(args []Concrete) *TypeMismatchErr {
	for _, arg := range args {
		if arg.ConcreteKind() != ConcreteKindScalar || arg.(*Scalar).Kind != ScalarKindString {
			return &TypeMismatchErr{
				ExpectedKind: string(ScalarKindString),
				Actual:       arg,
			}
		}
	}

	return nil
}

func (r *regexValidator) ValidateValues(args []Concrete) error {
	// The second argument is the regex pattern.
	regex, err := args[1].Interface()
	if err != nil {
		return err
	}

	if _, err := regexp.Compile(regex.(string)); err != nil {
		return fmt.Errorf("invalid regex: %w", err)
	}

	return nil
}

type existsValidator struct{}

func (e *existsValidator) ExpectedNumArgs() int {
	return 1
}

func (e *existsValidator) ValidateTypes(args []Concrete) *TypeMismatchErr {
	return nil
}

func (e *existsValidator) ValidateValues(args []Concrete) error {
	return nil
}

var functionValidators = map[string]FunctionValidator{
	"pow":       &powValidator{},
	"hello":     &helloValidator{},
	"regex":     &regexValidator{},
	"exists":    &existsValidator{},
	"notexists": &existsValidator{},
}

// LookupFuncValidator looks up a function by its name and returns its validator
// if it exists. If the function does not exist, then it is unrecognized and
// this function returns (nil, false).
// Note that the lookup of the function itself is a form of validation (that the
// function is even recognizable).
func LookupFuncValidator(funcName string) (FunctionValidator, bool) {
	funcValidator, ok := functionValidators[funcName]
	return funcValidator, ok
}
