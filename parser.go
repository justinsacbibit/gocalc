package gocalc

import (
	_ "fmt"
)

func newParser(expr string) *parser {
	return &parser{
		lexer: newLexer(expr),
	}
}

type parser struct {
	lexer lexer
}

func (p *parser) parseExpr() expr {
	return p.parse(p.lexer.token(), 0)
}

func (p *parser) parse(lhs token, minPrecedence int) expr {
	var root expr
	lookahead := p.lexer.token()
	for lookahead.typ == tokenPlus || lookahead.typ == tokenMinus {
		op := lookahead
		rhs := p.lexer.token()
		lookahead = p.lexer.token()
		binExpr := &binaryExpr{
			left: &literal{
				lhs.val,
			},
			right: &literal{
				rhs.val,
			},
			op: op,
		}
		if root == nil {
			root = binExpr
		}
	}
	if root == nil && lhs.typ != tokenEOF {
		root = &literal{
			val: lhs.val,
		}
	}
	return root
}
