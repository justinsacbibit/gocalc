package gocalc

import "strconv"

type valueExpr struct {
	val interface{}
}

func newValueExpr(val string) *valueExpr {
	r, _ := strconv.ParseInt(val, 0, 32)
	return &valueExpr{r}
}

func (val *valueExpr) accept(v exprVisitor) {
	v.visitValueExpr(val)
}
