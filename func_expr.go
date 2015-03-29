package gocalc

type funcExpr struct {
	function string
	args     []expr
}

func (f *funcExpr) accept(v exprVisitor) {
	v.visitFuncExpr(f)
}
