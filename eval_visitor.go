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

func (e *evaluator) Evaluate(t expr) (interface{}, error) {
	t.accept(e)
	return e.result, e.error
}

func (e *evaluator) err(format string, args ...interface{}) {
	e.error = newEvaluationError(fmt.Sprintf(format, args...))
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

	// TODO: consider support for integers as underlying type
	// consider eliminating type assertions for performance
	switch b.op.typ {
	case tokenLogicalOr:
		e.result = left.(bool) || right.(bool)
	case tokenLogicalAnd:
		e.result = left.(bool) && right.(bool)
	case tokenBitwiseOr:
		e.result = float64(int(left.(float64)) | int(right.(float64)))
	case tokenBitwiseAnd:
		e.result = float64(int(left.(float64)) & int(right.(float64)))
	case tokenBitwiseXor:
		e.result = float64(int(left.(float64)) ^ int(right.(float64)))
	case tokenEqual:
		e.result = left.(float64) == right.(float64)
	case tokenNotEqual:
		e.result = left.(float64) != right.(float64)
	case tokenLessThan:
		e.result = left.(float64) < right.(float64)
	case tokenLessOrEqual:
		e.result = left.(float64) <= right.(float64)
	case tokenGreaterThan:
		e.result = left.(float64) > right.(float64)
	case tokenGreaterOrEqual:
		e.result = left.(float64) >= right.(float64)
	case tokenLeftShift:
		e.result = float64(int(left.(float64)) << uint(right.(float64)))
	case tokenRightShift:
		e.result = float64(int(left.(float64)) >> uint(right.(float64)))
	case tokenPlus:
		e.result = left.(float64) + right.(float64)
	case tokenMinus:
		e.result = left.(float64) - right.(float64)
	case tokenStar:
		e.result = left.(float64) * right.(float64)
	case tokenSlash:
		e.result = left.(float64) / right.(float64)
	case tokenPercent:
		e.result = float64(int(left.(float64)) % int(right.(float64)))
	default:
		e.err("Unsupported binary operator %s", b.op)
	}
}

func createFunc(e *evaluator, arg expr) func() (interface{}, error) {
	return func() (interface{}, error) {
		return e.Evaluate(arg)
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

		r, _ := e.Evaluate(f.args[0])
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
		e.result = -(e.result.(float64))
	case tokenLogicalNot:
		e.result = !e.result.(bool)
	case tokenBitwiseNot:
		e.result = float64(^int(e.result.(float64)))
	default:
		e.err("Unsupported unary operator %s", u.op)
	}
}

func (e *evaluator) visitBoolExpr(b *boolExpr) {
	panic("not implemented")
}

func (e *evaluator) visitFloatExpr(f *floatExpr) {
	e.result = f.val
}

func (e *evaluator) visitIntExpr(i *intExpr) {
	e.result = float64(i.val)
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
