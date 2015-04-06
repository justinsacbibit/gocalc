package gocalc

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
	p := newParser(expr)
	t := p.parseExpr()
	if t == nil {
		return nil, newCompileError(p.error)
	}

	return &Expression{
		tree: t,
		raw:  expr,
	}, nil
}

// ParamResolver resolves the values of any identifiers within an Expression.
//
type ParamResolver func(string) interface{}

// FuncHandler handles evaluates a function within an Expression, given
// parameters (which are wrapped in a function for lazy evaluation).
//
type FuncHandler func(string, ...func() (interface{}, error)) (interface{}, error)

// Evaluate evaluates an Expression. If any parameters or function are found,
// Evaluate will call the appropriate resolver. The evaluation result is
// returned, or an error if evaluation failed.
//
func (e *Expression) Evaluate(p ParamResolver, f FuncHandler) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = r.(EvaluationError)
		}
	}()

	v := newEvaluator(p, f)
	return v.evaluate(e.tree), nil
}

// CompileError represents a compilation error of an expression.
//
type CompileError struct {
	s string
}

func (c *CompileError) Error() string {
	return c.s
}

func newCompileError(s string) *CompileError {
	return &CompileError{s}
}
