package gocalc

type exprVisitor interface {
	visitBinaryExpr(*binaryExpr)
	visitFuncExpr(*funcExpr)
	visitUnaryExpr(*unaryExpr)
	visitValueExpr(*valueExpr)
	visitParamExpr(*paramExpr)
}
