package gocalc

type evaluator struct {
	result interface{}
}

func newEvaluator() *evaluator {
	return &evaluator{}
}

func (e *evaluator) visitBinaryExpr(b *binaryExpr) {

}

func (e *evaluator) visitFuncExpr(f *funcExpr) {

}

func (e *evaluator) visitUnaryExpr(u *unaryExpr) {
	switch u.op.typ {
	case tokenMinus:
		e.result = -(e.result.(int))
	default:
		panic("unsupported")
	}
}

func (e *evaluator) visitValueExpr(v *valueExpr) {
	e.result = v.val
}

func (e *evaluator) visitParamExpr(p *paramExpr) {

}
