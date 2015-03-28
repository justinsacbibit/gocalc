package gocalc

import "fmt"

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
	// all expressions start with a primary
	primary := p.parsePrimary()
	if primary == nil {
		return nil
	}
	return p.parse(primary, 0)
}

func (p *parser) parsePrimary() expr {
	token := p.lexer.token()

	switch token.typ {
	case tokenNumber:
		// NUMBER
		return &literal{token.val}
	case tokenIdentifier:
		// IDENTIFIER | IDENTIFIER '(' args ')'
		return p.parseIdentifier(token)
	case tokenLeftParen:
		// '(' expression ')'
		e := p.parseExpr()
		token = p.lexer.token()
		if token.typ != tokenRightParen {
			p.error = fmt.Sprintln("Unclosed parenth")
			return nil
		}

		return e
	case tokenMinus, tokenPlus:
		// ( '+' | '-' ) * primary
		return &unaryExpr{
			expr: p.parseExpr(),
			op:   token,
		}
	default:
		p.error = fmt.Sprintf("Expected primary, got \"%s\"", token)
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
		arg := p.parseExpr()
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
		return &callExpr{
			function: token.val,
			args:     args,
		}
	default:
		return &identifier{token.val}
	}
}

func (p *parser) parse(lhs expr, minPrecedence int) expr {
	lookahead := p.lexer.peekToken()
	opPrecedence := precedence(lookahead)
	for binaryOp(lookahead) && opPrecedence >= minPrecedence {
		op := lookahead
		p.lexer.token()
		rhs := p.parsePrimary()
		if rhs == nil {
			return nil
		}
		lookahead = p.lexer.peekToken()
		lookaheadPrecedence := precedence(lookahead)
		for binaryOp(lookahead) && lookaheadPrecedence > opPrecedence {
			rhs = p.parse(rhs, lookaheadPrecedence)
			lookahead = p.lexer.peekToken()
		}
		lhs = &binaryExpr{
			left:  lhs,
			right: rhs,
			op:    op,
		}
	}
	return lhs
}

func binaryOp(lookahead *token) bool {
	return lookahead.typ > tokenBinaryOp
}

func precedence(lookahead *token) int {
	switch {
	case lookahead.typ > tokenMultiplicative:
		return 2
	case lookahead.typ > tokenAdditive:
		return 1
	default:
		return -1
	}
}
