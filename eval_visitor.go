package gocalc

import "fmt"

type evaluator struct {
	result        interface{}
	stop          bool
	error         error
	paramResolver ParamResolver
	funcHandler   FuncHandler
}

func newEvaluator(p ParamResolver, f FuncHandler) *evaluator {
	return &evaluator{
		paramResolver: p,
		funcHandler:   f,
	}
}

func (e *evaluator) evaluate(t expr) (interface{}, error) {
	t.accept(e)
	return e.result, e.error
}

func (e *evaluator) err(format string, args ...interface{}) {
	e.error = newEvaluationError(fmt.Sprintf(format, args...))
	// TODO: replace stop member with calls to panic?
	e.stop = true
}

func (e *evaluator) visitBinaryExpr(b *binaryExpr) {
	if e.stop {
		return
	}

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
				e.result = l & r
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
		e.err("Unsupported binary operator %s", b.op)
	}

	if e.result == nil {
		e.err("Binary operation type mismatch; left: %s, right: %s, op: %s", left, right, b.op)
	}
}

func createFunc(e *evaluator, arg expr) func() (interface{}, error) {
	return func() (interface{}, error) {
		return e.evaluate(arg)
	}
}

func (e *evaluator) mapLazy(args []expr) []func() (interface{}, error) {
	l := len(args)
	r := make([]func() (interface{}, error), l, l)
	for i, arg := range args {
		r[i] = createFunc(e, arg)
	}
	return r
}

func (e *evaluator) visitFuncExpr(f *funcExpr) {
	if e.stop {
		return
	}

	if e.funcHandler != nil {
		if res, err := e.funcHandler(f.function, e.mapLazy(f.args)...); err == nil && res != nil {
			e.result = res
			return
		} else if err != nil {
			e.err("Error occurred handling function %s: %s", f.function, err.Error())
			return
		}
	}

	switch f.function {
	case "abs":
		if l := len(f.args); l != 1 {
			e.err("abs takes one param, got %d", l)
			return
		}

		r, _ := e.evaluate(f.args[0])
		if f := r.(float64); f < 0 {
			e.result = -f
		}
	default:
		e.err("Unrecognized function %s", f.function)
	}
}

func (e *evaluator) visitUnaryExpr(u *unaryExpr) {
	if e.stop {
		return
	}

	u.expr.accept(e)

	switch u.op.typ {
	case tokenMinus:
		switch r := e.result.(type) {
		case int64:
			e.result = -r
		case float64:
			e.result = -r
		}
	case tokenLogicalNot:
		switch r := e.result.(type) {
		case bool:
			e.result = !r
		}
	case tokenBitwiseNot:
		switch r := e.result.(type) {
		case int64:
			e.result = ^r
		}
	default:
		e.err("Unsupported unary operator %s", u.op)
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
	if e.stop {
		return
	}

	if e.paramResolver != nil {
		if res := e.paramResolver(p.identifier); res != nil {
			e.result = res
			return
		}
	}

	e.err("Identifier \"%s\" undefined", p.identifier)
}

type EvaluationError struct {
	s string
}

func (e *EvaluationError) Error() string {
	return e.s
}

func newEvaluationError(s string) *EvaluationError {
	return &EvaluationError{s}
}
