package gocalc

type paramExpr struct {
	identifier string
}

func (p *paramExpr) accept(v exprVisitor) {
	v.visitParamExpr(p)
}
