package execution

import (
	"fmt"
	"log"
	"strconv"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/breeze"
)

func evaluateValueOrExpr(valOrExpr breeze.ValueOrExpr, datum datum.Datum) (interface{}, error) {
	breezeConst := evaluateValueOrExprToConst(valOrExpr, datum)

	switch breezeConst.Kind {
	case breeze.ConstKindBool:
		return breezeConst.Stringified == "true", nil
	case breeze.ConstKindNull:
		return nil, nil
	case breeze.ConstKindString:
		return breezeConst.Stringified, nil
	case breeze.ConstKindNumber:
		num, err := strconv.Atoi(breezeConst.Stringified)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q as a number: %w", breezeConst.Stringified, err)
		}
		return num, nil
	default:
		panic(fmt.Errorf("unrecognized kind: %q", breezeConst.Kind))
	}
}

func evaluateValueOrExprToConst(valOrExpr breeze.ValueOrExpr, datum datum.Datum) breeze.Const {
	if valOrExpr.Value != nil {
		// The simple case where this is already just a single value.
		return evaluateValue(valOrExpr.Value, datum)
	}

	// Otherwise, we must evaluate this expr:
	return evaluateExpr(valOrExpr.Expr, datum)
}

func evaluateValue(value breeze.Value, datum datum.Datum) breeze.Const {
	// TODO: Later, this can return things like function or field reference. In
	// these cases, we can't just cast into a native go type.
	switch value.GetKind() {
	case breeze.ValueKindConst:
		log.Println("CONST VALUE EVAL")
		constValue := value.(*breeze.Const)
		return *constValue
	case breeze.ValueKindFieldRef:
		log.Println("FIELD REF VALUE EVAL")
		return evaluateFieldRef(value.(*breeze.FieldRef), datum)
	default:
		panic(fmt.Sprintf("unrecognized value kind: %v", value.GetKind()))
	}
}

func evaluateFieldRef(fieldRef *breeze.FieldRef, datum datum.Datum) breeze.Const {
	panic("TODO")
}

func evaluateExpr(expr breeze.Expr, datum datum.Datum) breeze.Const {
	panic("TODO")
}
