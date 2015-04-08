package gocalc

import (
	"fmt"
	"runtime"
)

type evaluator struct {
	result        interface{}
	paramResolver ParamResolver
	funcHandler   FuncHandler
}

func newEvaluator(p ParamResolver, f FuncHandler) *evaluator {
	return &evaluator{
		paramResolver: p,
		funcHandler:   f,
	}
}

func (e *evaluator) evaluate(t expr) interface{} {
	defer func() {
		if r := recover(); r != nil {
			switch err := r.(type) {
			case runtime.TypeAssertionError:
				panic(EvaluationError(err.Error()))
			default:
				panic(r)
			}
		}
	}()

	t.accept(e)
	return e.result
}

// EvaluationError is the type of an expression evaluation error.
//
type EvaluationError string

// Error is EvaluationError's implementation of the error interface.
//
func (e EvaluationError) Error() string {
	return string(e)
}

func (e *evaluator) error(format string, args ...interface{}) {
	panic(EvaluationError(fmt.Sprintf(format, args...)))
}

func (e *evaluator) visitBinaryExpr(b *binaryExpr) {
	b.left.accept(e)
	left := e.result

	b.right.accept(e)
	right := e.result

	e.result = nil

	switch b.op.typ {
	case tokenLogicalOr:
		switch l := left.(type) {
		case bool:
			switch r := right.(type) {
			case bool:
				e.result = l || r
			}
		}
	case tokenLogicalAnd:
		switch l := left.(type) {
		case bool:
			switch r := right.(type) {
			case bool:
				e.result = l && r
			}
		}
	case tokenBitwiseOr:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l | r
			}
		}
	case tokenBitwiseAnd:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l & r
			}
		}
	case tokenBitwiseXor:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l ^ r
			}
		}
	case tokenEqual:
		switch l := left.(type) {
		default:
			switch r := right.(type) {
			default:
				e.result = l == r
			}
		}
	case tokenNotEqual:
		switch l := left.(type) {
		default:
			switch r := right.(type) {
			default:
				e.result = l != r
			}
		}
	case tokenLessThan:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l < r
			case float64:
				e.result = float64(l) < r
			}
		case float64:
			switch r := right.(type) {
			case float64:
				e.result = l < r
			case int64:
				e.result = l < float64(r)
			}
		}
	case tokenLessOrEqual:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l <= r
			case float64:
				e.result = float64(l) <= r
			}
		case float64:
			switch r := right.(type) {
			case float64:
				e.result = l <= r
			case int64:
				e.result = l <= float64(r)
			}
		}
	case tokenGreaterThan:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l > r
			case float64:
				e.result = float64(l) > r
			}
		case float64:
			switch r := right.(type) {
			case float64:
				e.result = l > r
			case int64:
				e.result = l > float64(r)
			}
		}
	case tokenGreaterOrEqual:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l >= r
			case float64:
				e.result = float64(l) >= r
			}
		case float64:
			switch r := right.(type) {
			case float64:
				e.result = l >= r
			case int64:
				e.result = l >= float64(r)
			}
		}
	case tokenLeftShift:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l << uint64(r)
			}
		}
	case tokenRightShift:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l >> uint64(r)
			}
		}
	case tokenPlus:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l + r
			case float64:
				e.result = float64(l) + r
			}
		case float64:
			switch r := right.(type) {
			case float64:
				e.result = l + r
			case int64:
				e.result = l + float64(r)
			}
		}
	case tokenMinus:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l - r
			case float64:
				e.result = float64(l) - r
			}
		case float64:
			switch r := right.(type) {
			case float64:
				e.result = l - r
			case int64:
				e.result = l - float64(r)
			}
		}
	case tokenStar:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l * r
			case float64:
				e.result = float64(l) * r
			}
		case float64:
			switch r := right.(type) {
			case float64:
				e.result = l * r
			case int64:
				e.result = l * float64(r)
			}
		}
	case tokenSlash:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l / r
			case float64:
				e.result = float64(l) / r
			}
		case float64:
			switch r := right.(type) {
			case float64:
				e.result = l / r
			case int64:
				e.result = l / float64(r)
			}
		}
	case tokenPercent:
		switch l := left.(type) {
		case int64:
			switch r := right.(type) {
			case int64:
				e.result = l % r
			}
		}
	default:
		e.error("Unsupported binary operator %v", b.op)
	}

	if e.result == nil {
		e.error("Binary operation type error; left: %v (%T), right: %v (%T), op: %v",
			left, left, right, right, b.op)
	}
}

func createFunc(e *evaluator, arg expr) func() interface{} {
	return func() interface{} {
		return e.evaluate(arg)
	}
}

func (e *evaluator) mapLazy(args []expr) []func() interface{} {
	l := len(args)
	r := make([]func() interface{}, l, l)
	for i, arg := range args {
		r[i] = createFunc(e, arg)
	}
	return r
}

func (e *evaluator) visitFuncExpr(f *funcExpr) {
	if e.funcHandler != nil {
		res, err := e.funcHandler(f.function, e.mapLazy(f.args)...)
		if err != nil {
			panic(err)
		} else if res != nil {
			e.result = res
			return
		}
	}

	switch f.function {
	case "abs":
		if l := len(f.args); l != 1 {
			e.error("abs takes one param, got %d", l)
			return
		}

		r := e.evaluate(f.args[0])
		switch f := r.(type) {
		case int64:
			if f < 0 {
				e.result = -f
			}
		case float64:
			if f < 0 {
				e.result = -f
			}
		}
	default:
		e.error("Unrecognized function %s", f.function)
	}
}

func (e *evaluator) visitUnaryExpr(u *unaryExpr) {
	u.expr.accept(e)
	operand := e.result

	e.result = nil

	switch u.op.typ {
	case tokenMinus:
		switch r := operand.(type) {
		case int64:
			e.result = -r
		case float64:
			e.result = -r
		}
	case tokenLogicalNot:
		switch r := operand.(type) {
		case bool:
			e.result = !r
		}
	case tokenBitwiseNot:
		switch r := operand.(type) {
		case int64:
			e.result = ^r
		}
	default:
		e.error("Unsupported unary operator %v", u.op)
	}

	if e.result == nil {
		e.error("Unary operation type mismatch; operator: %v, operand: %v (%T)", u.op, operand, operand)
	}
}

func (e *evaluator) visitBoolExpr(b *boolExpr) {
	e.result = b.val
}

func (e *evaluator) visitFloatExpr(f *floatExpr) {
	e.result = f.val
}

func (e *evaluator) visitIntExpr(i *intExpr) {
	e.result = i.val
}

func (e *evaluator) visitParamExpr(p *paramExpr) {
	if e.paramResolver != nil {
		if res := e.paramResolver(p.identifier); res != nil {
			e.result = res
			return
		}
	}

	e.error("Identifier \"%s\" undefined", p.identifier)
}
