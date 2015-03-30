package gocalc

import "strconv"

type valueExpr struct {
	val interface{}
}

func newValueExpr(val string, float bool) *valueExpr {
	var r interface{}
	if true {
		r, _ = strconv.ParseFloat(val, 64)
	} else {
		i, _ := strconv.ParseInt(val, 0, 32)
		r = int(i)
	}
	return &valueExpr{r}
}

func (val *valueExpr) accept(v exprVisitor) {
	v.visitValueExpr(val)
}
