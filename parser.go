package gocalc

import (
	"fmt"
	"strconv"
)

func newParser(l lexer) *parser {
	return &parser{
		lexer: l,
	}
}

type parser struct {
	lexer lexer
	error string
}

func (p *parser) parseExpr() expr {
	e := p.parse(0)
	if e == nil {
		return nil
	} else if next := p.lexer.token(); next == nil || next.typ != tokenEOF {
		p.error = fmt.Sprintf("parseExpr(): Expected EOF, got %s", next)
		return nil
	}
	return e
}

func (p *parser) parsePrimary() expr {
	token := p.lexer.token()

	switch token.typ {
	case tokenMinus, tokenPlus, tokenLogicalNot, tokenBitwiseNot:
		e := p.parse(precedence(token, unary))
		if e == nil {
			return nil
		}
		return &unaryExpr{
			expr: e,
			op:   token,
		}
	case tokenLeftParen:
		e := p.parse(0)
		token = p.lexer.token()
		if token.typ != tokenRightParen {
			p.error = fmt.Sprintln("Unclosed parenth")
			return nil
		}
		return e
	case tokenInt:
		i, _ := strconv.ParseInt(token.val, 0, 32)
		return &intExpr{i}
	case tokenFloat:
		f, _ := strconv.ParseFloat(token.val, 64)
		return &floatExpr{f}
	case tokenTrue:
		return &boolExpr{true}
	case tokenFalse:
		return &boolExpr{false}
	case tokenIdentifier:
		// IDENTIFIER | IDENTIFIER '(' args ')'
		return p.parseIdentifier(token)
	case tokenError:
		p.error = fmt.Sprintf("Lexical error at %d: \"%s\"", token.pos, token.val)
		return nil
	default:
		p.error = fmt.Sprintf("Expected primary at %d, got \"%s\"", token.pos, token)
		return nil
	}
}

func (p *parser) parseFunctionArgs() []expr {
	peek := p.lexer.peekToken()
	funcArgs := []expr{}
	if peek.typ == tokenRightParen {
		// IDENTIFIER '(' ')'
		p.lexer.token()
		return funcArgs
	}

	for {
		arg := p.parse(0)
		if arg == nil {
			return nil
		}
		funcArgs = append(funcArgs, arg)
		if peek = p.lexer.peekToken(); peek.typ == tokenComma {
			p.lexer.token()
		} else if peek.typ == tokenRightParen {
			p.lexer.token()
			break
		} else {
			p.error = fmt.Sprintf("Expected a comma or right paren after function argument, got \"%s\"", peek)
			return nil
		}
	}
	return funcArgs
}

func (p *parser) parseIdentifier(token *token) expr {
	peeked := p.lexer.peekToken()
	switch peeked.typ {
	case tokenLeftParen:
		p.lexer.token()
		args := p.parseFunctionArgs()
		if args == nil {
			return nil
		}
		return &funcExpr{
			function: token.val,
			args:     args,
		}
	default:
		return &paramExpr{token.val}
	}
}

func (p *parser) parse(prec int) expr {
	e := p.parsePrimary()
	if e == nil {
		return nil
	}
	lookahead := p.lexer.peekToken()
	for binaryOp(lookahead) && precedence(lookahead, binary) >= prec {
		op := lookahead
		p.consume()
		q := 1 + precedence(lookahead, binary)
		r := p.parse(q)
		if r == nil {
			return nil
		}
		e = &binaryExpr{
			left:  e,
			right: r,
			op:    op,
		}

		lookahead = p.lexer.peekToken()
	}
	return e
}

func (p *parser) consume() {
	p.lexer.token()
}

func binaryOp(token *token) bool {
	return token.typ > tokenBinary
}

type operatorType int

const (
	unary operatorType = iota
	binary
	// ternary
)

func precedence(token *token, operatorType operatorType) int {
	switch operatorType {
	case unary:
		switch token.typ {
		case tokenMinus, tokenLogicalNot, tokenBitwiseNot:
			return 10
		}
	case binary:
		switch token.typ {
		case tokenLogicalOr:
			return 0
		case tokenLogicalAnd:
			return 1
		case tokenBitwiseOr:
			return 2
		case tokenBitwiseXor:
			return 3
		case tokenBitwiseAnd:
			return 4
		case tokenEqual, tokenNotEqual:
			return 5
		case tokenLessThan, tokenLessOrEqual, tokenGreaterThan, tokenGreaterOrEqual:
			return 6
		case tokenLeftShift, tokenRightShift:
			return 7
		case tokenPlus, tokenMinus:
			return 8
		case tokenStar, tokenSlash, tokenPercent:
			return 9
		}
	}

	return -1
}
