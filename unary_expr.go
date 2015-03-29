package gocalc

type unaryExpr struct {
	expr expr
	op   *token
}

func (u *unaryExpr) accept(v exprVisitor) {
	v.visitUnaryExpr(u)
}
