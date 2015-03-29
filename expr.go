package gocalc

type expr interface {
	accept(exprVisitor)
}
