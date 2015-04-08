package gocalc

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type lexer interface {
	token() *token
	peekToken() *token
}

func newLexer(input string) lexer {
	return &gocalcLexer{
		input:  input,
		state:  initialState,
		tokens: queue{},
		alpha:  "0123456789abcdefABCDEFghijklmnopqrstuvwxyzGHIJKLMNOPQRSTUVWXYZ",
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
F:
	for {
		r := l.next()
		switch {
		case r == eof:
			break F
		case r == ' ':
			l.ignore()
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'):
			l.backup()
			return lexIdentifier
		case r >= '0' && r <= '9':
			l.backup()
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
	if l.accept("0") {
		return lexZero
	}

	l.acceptRun(l.alpha[0:10])
	if l.accept(".") {
		return lexFloat
	}
	if isAlpha(l.peek()) {
		return lexNumberError
	}
	l.emit(tokenInt)
	return initialState
}

func lexZero(l *gocalcLexer) stateFn {
	if l.accept("x") {
		return lexHex
	} else if l.accept("b") {
		return lexBinary
	} else if l.accept(".") {
		return lexFloat
	} else if l.accept(l.alpha[0:8]) {
		l.backup()
		return lexOctal
	}

	if r := l.peek(); isAlpha(r) || r >= '0' && r <= '9' {
		return lexNumberError
	}

	l.emit(tokenInt)
	return initialState
}

func lexFloat(l *gocalcLexer) stateFn {
	l.acceptRun(l.alpha[0:10])

	if isAlpha(l.peek()) {
		return lexNumberError
	}

	l.emit(tokenFloat)
	return initialState
}

func lexOctal(l *gocalcLexer) stateFn {
	l.acceptRun(l.alpha[0:8])

	if isAlpha(l.peek()) {
		return lexNumberError
	}

	l.emit(tokenInt)
	return initialState
}

func lexBinary(l *gocalcLexer) stateFn {
	if !l.acceptRun(l.alpha[0:2]) {
		return lexNumberError
	}

	if isAlpha(l.peek()) {
		return lexNumberError
	}

	l.emit(tokenInt)
	return initialState
}

func lexHex(l *gocalcLexer) stateFn {
	if !l.acceptRun(l.alpha[0:22]) {
		return lexNumberError
	}

	if isAlpha(l.peek()) {
		return lexNumberError
	}

	l.emit(tokenInt)
	return initialState
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func lexNumberError(l *gocalcLexer) stateFn {
	return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
}

func lexIdentifier(l *gocalcLexer) stateFn {
	t := l.accept("true")
	f := l.accept("false")
	l.acceptRun(l.alpha)
	if length := l.pos - l.start; length == 4 && t {
		l.emit(tokenTrue)
	} else if length == 5 && f {
		l.emit(tokenFalse)
	} else {
		l.emit(tokenIdentifier)
	}
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

func (l *gocalcLexer) acceptRun(valid string) bool {
	r := false
	for strings.IndexRune(valid, l.next()) >= 0 {
		r = true
	}
	l.backup()
	return r
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
