package gocalc

// An Expression does
//
type Expression struct {
	parser *parser
	tree   expr
	raw    string
}

// NewExpr compiles a string expression
//
func NewExpr(expr string) (*Expression, error) {
	p := newParser(expr)
	t := p.parseExpr()
	if t == nil {
		return nil, newCompileError(p.error)
	}
	return &Expression{
		parser: newParser(expr),
		tree:   t,
		raw:    expr,
	}, nil
}

// Evaluate evaluates expression
//
func (e *Expression) Evaluate() (interface{}, error) {
	v := newEvaluator()
	e.tree.accept(v)
	return v.result, nil
}

// CompileError is
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
