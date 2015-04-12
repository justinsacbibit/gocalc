package gocalc

import "fmt"

// An Expression is used to compile and evaluate a string representation of a
// mathematical expression. Multiple goroutines can use an Expression, as
// parameters and functions are not stored within the Expression, and are
// instead resolved during each evaluation.
//
type Expression struct {
	tree expr
	raw  string
}

// NewExpr initializes and returns an Expression given the string
// representation, or an error if compilation failed.
//
func NewExpr(expr string) (*Expression, error) {
	l := newLexer(expr)
	p := newParser(l)
	t := p.parseExpr()
	if t == nil {
		return nil, compileError(p.error)
	}

	return &Expression{
		tree: t,
		raw:  expr,
	}, nil
}

// ParamResolver resolves the values of any identifiers within an Expression.
//
type ParamResolver func(string) (value interface{})

// FuncHandler handles evaluates a function within an Expression, given
// parameters (which are wrapped in a function for lazy evaluation).
// FuncHandlers may make calls to panic() with an EvaluationError.
//
type FuncHandler func(string, ...func() interface{}) (result interface{}, handled bool)

// Evaluate evaluates an Expression. If any parameters or function are found,
// Evaluate will call the appropriate resolver. The evaluation result is
// returned, or an error if evaluation failed.
//
func (e *Expression) Evaluate(p ParamResolver, f FuncHandler) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch er := r.(type) {
			case error:
				err = er
			default:
				err = EvaluationError(fmt.Sprintf("An error has occurred: %s", er))
			}
		}
	}()

	v := newEvaluator(p, f)
	return v.evaluate(e.tree), nil
}

// CompileError represents a compilation error of an expression.
//
type compileError string

func (c compileError) Error() string {
	return string(c)
}
