package gocalc

type valueExpr struct {
	val string
}

func (val *valueExpr) accept(v exprVisitor) {
	v.visitValueExpr(val)
}
