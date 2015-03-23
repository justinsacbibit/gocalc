package gocalc

import (
	_ "fmt"
	"strings"
	"testing"
)

func TestLexDigit(t *testing.T) {
	s := "1"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenNumber, t)
	checkVal(to, s, t)

	to = l.token()
	checkType(to, tokenEOF, t)
}

func TestLexMultiDigit(t *testing.T) {
	s := "111"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenNumber, t)
	checkVal(to, s, t)
}

func TestLexMultipleNumbers(t *testing.T) {
	s := "1 2 3"
	l := newLexer(s)

	strs := strings.Split(s, " ")
	for _, str := range strs {
		to := l.token()
		checkType(to, tokenNumber, t)
		checkVal(to, str, t)
	}
}

func TestLexNumberWithLeadingWhitespace(t *testing.T) {
	s := "5"
	l := newLexer("  " + s)
	to := l.token()

	checkType(to, tokenNumber, t)
	checkVal(to, s, t)
}

func TestLexNumberWithSurroundingWhitespace(t *testing.T) {
	s := "99"
	l := newLexer("  " + s + "  ")
	to := l.token()

	checkType(to, tokenNumber, t)
	checkVal(to, s, t)
}

func TestLexPlus(t *testing.T) {
	s := "+"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenPlus, t)
}

func TestLexMinus(t *testing.T) {
	s := "-"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenMinus, t)
}

func TestLexLeftParen(t *testing.T) {
	s := "("
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenLeftParen, t)
}

func TestLexRightParen(t *testing.T) {
	s := ")"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenRightParen, t)
}

func TestLexNumberWithSurroundingParentheses(t *testing.T) {
	s := "(5)"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenLeftParen, t)

	to = l.token()
	checkType(to, tokenNumber, t)
	checkVal(to, "5", t)

	to = l.token()
	checkType(to, tokenRightParen, t)
}

func TestLexComplex(t *testing.T) {
	s := "5 - (10 - 5)"
	l := newLexer(s)
	eTypes := []tokenType{tokenNumber, tokenMinus, tokenLeftParen, tokenNumber, tokenMinus, tokenNumber, tokenRightParen}
	eVals := []string{"5", "-", "(", "10", "-", "5", ")"}

	for i, eType := range eTypes {
		eVal := eVals[i]
		to := l.token()

		checkType(to, eType, t)
		checkVal(to, eVal, t)
	}
}

func TestLexStar(t *testing.T) {
	s := "*"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenStar, t)
}

func TestLexSlash(t *testing.T) {
	s := "/"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenSlash, t)
}

func TestLexStarAndSlash(t *testing.T) {
	s := "5*1/2"
	l := newLexer(s)

	eTypes := []tokenType{tokenNumber, tokenStar, tokenNumber, tokenSlash, tokenNumber}
	for _, eType := range eTypes {
		to := l.token()
		checkType(to, eType, t)
	}
}

func TestLexIdentifier(t *testing.T) {
	s := "x"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenIdentifier, t)
	checkVal(to, s, t)

	checkEof(l, t)
}

func checkVal(to token, exp string, t *testing.T) {
	if to.val != exp {
		t.Errorf("Wrong token value: %s", to.val)
	}
}

func checkType(to token, exp tokenType, t *testing.T) {
	if to.typ != exp {
		t.Errorf("Wrong token type: %d, expected: %d", to.typ, exp)
	}
}

func checkEof(l lexer, t *testing.T) {
	to := l.token()
	if to.typ != tokenEOF {
		t.Errorf("Token is not EOF: %d", to.typ)
	}
}
