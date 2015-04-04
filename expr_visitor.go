package gocalc

type exprVisitor interface {
	visitBinaryExpr(*binaryExpr)
	visitFuncExpr(*funcExpr)
	visitUnaryExpr(*unaryExpr)
	visitParamExpr(*paramExpr)

	visitBoolExpr(*boolExpr)
	visitFloatExpr(*floatExpr)
	visitIntExpr(*intExpr)
}
