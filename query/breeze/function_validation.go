package breeze

// FunctionValidator validates the use of a function.
type FunctionValidator interface {
	// ExpectedNumArgs returns the expected number of arguments.
	ExpectedNumArgs() int
	// ValidateTypes determines if the given arguments are of valid types.
	// It assumes that the members of the given slice of *Const are non-nil.
	ValidateTypes(args []Const) *TypeMismatchErr
}

type helloValidator struct{}

func (p *helloValidator) ExpectedNumArgs() int {
	return 0
}

func (p *helloValidator) ValidateTypes(args []Const) *TypeMismatchErr {
	return nil
}

type powValidator struct{}

func (p *powValidator) ExpectedNumArgs() int {
	return 2
}

func (p *powValidator) ValidateTypes(args []Const) *TypeMismatchErr {
	for _, arg := range args {
		if arg.Kind != ConstKindNumber {
			return &TypeMismatchErr{
				ExpectedKind: ConstKindNumber,
				Actual:       arg,
			}
		}
	}

	return nil
}

var functionValidators = map[string]FunctionValidator{
	"pow":   &powValidator{},
	"hello": &helloValidator{},
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
