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
	panic("peek failed")
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
		case r >= '0' && r <= '9':
			return lexNumber
		case r == '+':
			l.emit(tokenPlus)
		case r == '-':
			l.emit(tokenMinus)
		case r == '*':
			l.emit(tokenStar)
		case r == '/':
			l.emit(tokenSlash)
		case r == ',':
			l.emit(tokenComma)
		case r == '(':
			l.emit(tokenLeftParen)
		case r == ')':
			l.emit(tokenRightParen)
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'):
			return lexIdentifier
		default:
			return l.errorf("Invalid token: %c", r)
		}
	}

	l.emit(tokenEOF)
	return nil
}

func lexNumber(l *gocalcLexer) stateFn {
	digits := "0123456789"
	l.acceptRun(digits)
	if r := l.peek(); (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(tokenNumber)
	return initialState
}

func lexIdentifier(l *gocalcLexer) stateFn {
	alpha := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	l.acceptRun(alpha)
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
