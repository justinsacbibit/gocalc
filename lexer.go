package gocalc

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func newLexer(input string) *gocalcLexer {
	return &gocalcLexer{
		input:  input,
		state:  initialState,
		tokens: queue{},
		alpha:  "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		digits: "0123456789",
	}
}

func (l *gocalcLexer) token() *token {
	for {
		switch len(l.tokens) {
		case 0:
			l.state = l.state(l)
		default:
			return l.tokens.pop()
		}
	}
	panic("not reached")
}

func (l *gocalcLexer) peekToken() *token {
	for {
		switch len(l.tokens) {
		case 0:
			l.state = l.state(l)
		default:
			return l.tokens.first()
		}
	}
	panic("not reached")
}

// mark: Internal use

const eof = -1

type stateFn func(*gocalcLexer) stateFn

type gocalcLexer struct {
	input  string
	start  int
	pos    int
	width  int
	tokens queue
	state  stateFn
	digits string
	alpha  string
}

func (l *gocalcLexer) emit(t tokenType) {
	l.tokens.push(&token{
		typ: t,
		val: l.input[l.start:l.pos],
	})
	l.start = l.pos
}

func initialState(l *gocalcLexer) stateFn {
	for {
		r := l.next()
		switch {
		case r == eof:
			l.emit(tokenEOF)
			return nil
		case r == ' ':
			l.ignore()
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'):
			return lexIdentifier
		case r >= '0' && r <= '9':
			return lexNumber
		case r == '(':
			l.emit(tokenLeftParen)
		case r == ')':
			l.emit(tokenRightParen)
		case r == ',':
			l.emit(tokenComma)
		case r == '!':
			return lexExclamation
		case r == '~':
			l.emit(tokenBitwiseNot)
		case r == '*':
			l.emit(tokenStar)
		case r == '/':
			l.emit(tokenSlash)
		case r == '%':
			l.emit(tokenPercent)
		case r == '+':
			l.emit(tokenPlus)
		case r == '-':
			l.emit(tokenMinus)
		case r == '<':
			return lexLessThan
		case r == '>':
			return lexGreaterThan
		case r == '=':
			l.emit(tokenEqual)
		case r == '&':
			return lexAnd
		case r == '^':
			l.emit(tokenBitwiseXor)
		case r == '|':
			return lexOr
		default:
			return l.errorf("Invalid token: %c", r)
		}
	}

	l.emit(tokenEOF)
	return nil
}

func (l *gocalcLexer) checkNext(next rune, match tokenType, otherwise tokenType) stateFn {
	if l.next() == next {
		l.emit(match)
	} else {
		l.backup()
		l.emit(otherwise)
	}

	return initialState
}

func lexExclamation(l *gocalcLexer) stateFn {
	return l.checkNext('=', tokenNotEqual, tokenLogicalNot)
}

func lexLessThan(l *gocalcLexer) stateFn {
	switch l.next() {
	case '=':
		l.emit(tokenLessOrEqual)
	case '<':
		l.emit(tokenLeftShift)
	default:
		l.backup()
		l.emit(tokenLessThan)
	}

	return initialState
}

func lexGreaterThan(l *gocalcLexer) stateFn {
	switch l.next() {
	case '=':
		l.emit(tokenGreaterOrEqual)
	case '>':
		l.emit(tokenRightShift)
	default:
		l.backup()
		l.emit(tokenGreaterThan)
	}

	return initialState
}

func lexAnd(l *gocalcLexer) stateFn {
	return l.checkNext('&', tokenLogicalAnd, tokenBitwiseAnd)
}

func lexOr(l *gocalcLexer) stateFn {
	return l.checkNext('|', tokenLogicalOr, tokenBitwiseOr)
}

func lexNumber(l *gocalcLexer) stateFn {
	l.acceptRun(l.digits)
	if r := l.peek(); (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(tokenNumber)
	return initialState
}

func lexIdentifier(l *gocalcLexer) stateFn {
	l.acceptRun(l.alpha)
	l.emit(tokenIdentifier)
	return initialState
}

func (l *gocalcLexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens.push(&token{
		tokenError,
		fmt.Sprintf(format, args...),
	})
	return nil
}

func (l *gocalcLexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *gocalcLexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *gocalcLexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *gocalcLexer) backup() {
	l.pos -= l.width
}

func (l *gocalcLexer) ignore() {
	l.start = l.pos
}

func (l *gocalcLexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}
