package gocalc

import (
	"fmt"
)

func newParser(expr string) *parser {
	return &parser{
		lexer: newLexer(expr),
	}
}

type parser struct {
	lexer lexer
	error string
}

func (p *parser) parseExpr() expr {
	return p.parse(p.parsePrimary(), 0)
}

func (p *parser) parsePrimary() expr {
	token := p.lexer.token()

	switch token.typ {
	case tokenNumber:
		return &literal{token.val}
	case tokenLeftParen:
		e := p.parseExpr()
		token = p.lexer.token()
		if token.typ != tokenRightParen {
			p.error = fmt.Sprintln("Unclosed parenth")
		} else {
			return e
		}
	case tokenMinus, tokenPlus:
		return &unaryExpr{
			expr: p.parseExpr(),
			op:   token,
		}
	default:
		p.error = fmt.Sprintln("Invalid token:", token)
	}

	return nil
}

func (p *parser) parse(lhs expr, minPrecedence int) expr {
	lookahead := p.lexer.peekToken()
	for lookahead.typ == tokenPlus || lookahead.typ == tokenMinus {
		op := lookahead
		p.lexer.token()
		rhs := p.parsePrimary()
		if rhs == nil {
			return nil
		}
		lookahead = p.lexer.peekToken()
		lhs = &binaryExpr{
			left:  lhs,
			right: rhs,
			op:    op,
		}
	}
	return lhs
}
