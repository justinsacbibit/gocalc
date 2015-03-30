package gocalc

import "fmt"

type ParamResolver func(string) interface{}
type FuncResolver func(string, ...func() (interface{}, error)) (interface{}, error)

type Evaluator struct {
	result        interface{}
	stop          bool
	error         error
	paramResolver ParamResolver
	funcResolver  FuncResolver
}

func newEvaluator(p ParamResolver, f FuncResolver) *Evaluator {
	return &Evaluator{
		paramResolver: p,
		funcResolver:  f,
	}
}

func (e *Evaluator) Evaluate(t expr) (interface{}, error) {
	t.accept(e)
	return e.result, e.error
}

func (e *Evaluator) err(format string, args ...interface{}) {
	e.error = newEvaluationError(fmt.Sprintf(format, args...))
	e.stop = true
}

func (e *Evaluator) visitBinaryExpr(b *binaryExpr) {
	if e.stop {
		return
	}

	b.left.accept(e)
	left := e.result
	b.right.accept(e)
	right := e.result

	switch b.op.typ {
	case tokenPlus:
		e.result = left.(float64) + right.(float64)
	case tokenMinus:
		e.result = left.(float64) - right.(float64)
	case tokenStar:
		e.result = left.(float64) * right.(float64)
	case tokenSlash:
		e.result = left.(float64) / right.(float64)
	default:
		e.err("Unsupported binary operator %s", b.op)
	}
}

func createFunc(e *Evaluator, arg expr) func() (interface{}, error) {
	return func() (interface{}, error) {
		return e.Evaluate(arg)
	}
}

func (e *Evaluator) mapLazy(args []expr) []func() (interface{}, error) {
	l := len(args)
	r := make([]func() (interface{}, error), l, l)
	for i, arg := range args {
		r[i] = createFunc(e, arg)
	}
	return r
}

func (e *Evaluator) visitFuncExpr(f *funcExpr) {
	if e.stop {
		return
	}

	if e.funcResolver != nil {
		if res, err := e.funcResolver(f.function, e.mapLazy(f.args)...); err == nil && res != nil {
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

func (e *Evaluator) visitUnaryExpr(u *unaryExpr) {
	if e.stop {
		return
	}

	u.expr.accept(e)

	switch u.op.typ {
	case tokenMinus:
		e.result = -(e.result.(float64))
	default:
		e.err("Unsupported unary operator %s", u.op)
	}
}

func (e *Evaluator) visitValueExpr(v *valueExpr) {
	if e.stop {
		return
	}

	e.result = v.val
}

func (e *Evaluator) visitParamExpr(p *paramExpr) {
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
