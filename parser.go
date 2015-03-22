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
	return p.parse(p.parsePrimary(), 0)
}

func (p *parser) parsePrimary() node {
	token := p.lexer.token()
	return &literal{
		val: token.val,
	}
}

func (p *parser) parse(lhs expr, minPrecedence int) expr {
	lookahead := p.lexer.token()
	for lookahead.typ == tokenPlus || lookahead.typ == tokenMinus {
		op := lookahead
		rhs := p.parsePrimary()
		lookahead = p.lexer.token()
		lhs = &binaryExpr{
			left:  lhs,
			right: rhs,
			op:    op,
		}
	}
	return lhs
}
