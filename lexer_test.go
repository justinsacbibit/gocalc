package gocalc

import "testing"

func TestLexDigit(t *testing.T) {
	s := "1"
	shouldLex(s, types(tokenNumber), vals(s), t)
}

func TestLexMultiDigit(t *testing.T) {
	s := "111"
	shouldLex(s, types(tokenNumber), vals(s), t)
}

func TestLexMultipleNumbers(t *testing.T) {
	s := "1 2 3"
	shouldLex(s,
		types(tokenNumber, tokenNumber, tokenNumber),
		vals("1", "2", "3"),
		t)
}

func TestLexNumberWithLeadingWhitespace(t *testing.T) {
	shouldLex("  5", types(tokenNumber), vals("5"), t)
}

func TestLexNumberWithSurroundingWhitespace(t *testing.T) {
	shouldLex("99", types(tokenNumber), vals("99"), t)
}

func TestLexPlus(t *testing.T) {
	shouldLex("+", types(tokenPlus), nil, t)
}

func TestLexMinus(t *testing.T) {
	shouldLex("-", types(tokenMinus), nil, t)
}

func TextLexComma(t *testing.T) {
	shouldLex(",", types(tokenComma), nil, t)
}

func TestLexLeftParen(t *testing.T) {
	shouldLex("(", types(tokenLeftParen), nil, t)
}

func TestLexRightParen(t *testing.T) {
	shouldLex(")", types(tokenRightParen), nil, t)
}

func TestLexNumberWithSurroundingParentheses(t *testing.T) {
	shouldLex("(5)",
		types(tokenLeftParen, tokenNumber, tokenRightParen),
		vals("", "5", ""),
		t)
}

func TestLexComplex(t *testing.T) {
	shouldLex("5 - (10 - 5)",
		types(tokenNumber, tokenMinus, tokenLeftParen, tokenNumber, tokenMinus, tokenNumber, tokenRightParen),
		vals("5", "", "", "10", "", "5", ""),
		t)
}

func TestLexStar(t *testing.T) {
	shouldLex("*", types(tokenStar), nil, t)
}

func TestLexSlash(t *testing.T) {
	shouldLex("/", types(tokenSlash), nil, t)
}

func TestLexStarAndSlash(t *testing.T) {
	shouldLex("5*1/2",
		types(tokenNumber, tokenStar, tokenNumber, tokenSlash, tokenNumber),
		vals("5", "", "1", "", "2"),
		t)
}

func TestLexIdentifier(t *testing.T) {
	s := "x"
	shouldLex(s, types(tokenIdentifier), vals(s), t)
}

func TestLexFunc(t *testing.T) {
	shouldLex("f(x)",
		types(tokenIdentifier, tokenLeftParen, tokenIdentifier, tokenRightParen),
		vals("f", "", "x", ""),
		t)
}

func TestLexPercent(t *testing.T) {
	shouldLex("%", types(tokenPercent), nil, t)
}

func TestLexLeftShift(t *testing.T) {
	shouldLex("<<", types(tokenLeftShift), nil, t)
}

func TestLexRightShift(t *testing.T) {
	shouldLex(">>", types(tokenRightShift), nil, t)
}

func TestLexLessThan(t *testing.T) {
	shouldLex("<", types(tokenLessThan), nil, t)
}

func TestLexLessOrEqual(t *testing.T) {
	shouldLex("<=", types(tokenLessOrEqual), nil, t)
}

func TestLexGreaterThan(t *testing.T) {
	shouldLex(">", types(tokenGreaterThan), nil, t)
}

func TestLexGreaterOrEqual(t *testing.T) {
	shouldLex(">=", types(tokenGreaterOrEqual), nil, t)
}

func TestLexEqual(t *testing.T) {
	shouldLex("=", types(tokenEqual), nil, t)
}

func TestLexNotEqual(t *testing.T) {
	shouldLex("!=", types(tokenNotEqual), nil, t)
}

func TestLexBitwiseAnd(t *testing.T) {
	shouldLex("&", types(tokenBitwiseAnd), nil, t)
}

func TestLexBitwiseXor(t *testing.T) {
	shouldLex("^", types(tokenBitwiseXor), nil, t)
}

func TestLexBitwiseOr(t *testing.T) {
	shouldLex("|", types(tokenBitwiseOr), nil, t)
}

func TestLexLogicalAnd(t *testing.T) {
	shouldLex("&&", types(tokenLogicalAnd), nil, t)
}

func TestLexLogicalOr(t *testing.T) {
	shouldLex("||", types(tokenLogicalOr), nil, t)
}

// Helpers

func shouldLex(s string, ts []tokenType, v *[]string, t *testing.T) {
	l := newLexer(s)
	for i, e := range ts {
		to := l.token()
		if to.typ != e {
			t.Errorf("Wrong token type: expected %s, got %s", e, to.typ)
		}
		if v != nil {
			ev := (*v)[i]
			if ev != "" && to.val != ev {
				t.Errorf("Wrong token value: expected %s, got %s", ev, to.val)
			}
		}
	}
	if typ := l.token().typ; typ != tokenEOF {
		t.Errorf("Expected EOF, got %s", typ)
	}
}

func types(args ...tokenType) []tokenType {
	return args
}

func vals(args ...string) *[]string {
	return &args
}
