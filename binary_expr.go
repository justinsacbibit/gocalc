package gocalc

type binaryExpr struct {
	left  expr
	right expr
	//opPos int
	op *token
}

func (b *binaryExpr) accept(v exprVisitor) {
	v.visitBinaryExpr(b)
}
