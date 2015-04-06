package gocalc

// All expressions implement the expr interface.
type expr interface {
	accept(exprVisitor)
}

type (
	// A binaryExpr represents a binary expression.
	binaryExpr struct {
		left  expr   // left operand
		op    *token // binary operator
		right expr   // right operand
	}

	// A funcExpr represents a function call.
	funcExpr struct {
		function string // function name
		args     []expr // argument list
	}

	// A paramExpr represents a parameter.
	paramExpr struct {
		identifier string // parameter name
	}

	// A unaryExpr represents a unary expression.
	unaryExpr struct {
		expr expr   // operand
		op   *token // unary operator
	}

	// A boolExpr represents a boolean literal.
	boolExpr struct {
		val bool
	}

	// A floatExpr represents a float literal.
	floatExpr struct {
		val float64
	}

	// An intExpr represents a integer literal.
	intExpr struct {
		val int64
	}
)

func (b *binaryExpr) accept(v exprVisitor) {
	v.visitBinaryExpr(b)
}

func (f *funcExpr) accept(v exprVisitor) {
	v.visitFuncExpr(f)
}

func (p *paramExpr) accept(v exprVisitor) {
	v.visitParamExpr(p)
}

func (u *unaryExpr) accept(v exprVisitor) {
	v.visitUnaryExpr(u)
}

func (b *boolExpr) accept(v exprVisitor) {
	v.visitBoolExpr(b)
}

func (f *floatExpr) accept(v exprVisitor) {
	v.visitFloatExpr(f)
}

func (i *intExpr) accept(v exprVisitor) {
	v.visitIntExpr(i)
}
