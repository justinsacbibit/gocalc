package gocalc

import "testing"

type lexerSingleTokenTest struct {
	ok    bool
	input string
	typ   tokenType
	val   string
}

type lexerMultipleTokenTest struct {
	ok    bool
	input string
	types []tokenType
	vals  []string
}

var lexerSingleTokenTests = []lexerSingleTokenTest{
	{true, "1", tokenInt, "1"},
	{true, "111", tokenInt, "111"},
	{true, "  5", tokenInt, "5"},
	{true, "  99  ", tokenInt, "99"},

	{true, "1.2", tokenFloat, "1.2"},
	{true, "0.5", tokenFloat, "0.5"},
	{true, "0xF", tokenInt, "0xF"},
	{false, "0xG", tokenError, ""},
	{true, "0b10", tokenInt, "0b10"},
	{false, "0b3", tokenError, ""},
	{false, "0b", tokenError, ""},
	{true, "077", tokenInt, "077"},
	{false, "08", tokenError, ""},
	{true, "0", tokenInt, "0"},
	{false, ".", tokenError, ""},
	{false, "5..", tokenError, ""},
	{false, "3..5", tokenError, ""},
	{true, "0.", tokenFloat, "0."},

	{true, ",", tokenComma, ""},
	{true, "(", tokenLeftParen, ""},
	{true, ")", tokenRightParen, ""},

	{true, "+", tokenPlus, ""},
	{true, "-", tokenMinus, ""},
	{true, "*", tokenStar, ""},
	{true, "/", tokenSlash, ""},
	{true, "%", tokenPercent, ""},
	{true, "<<", tokenLeftShift, ""},
	{true, ">>", tokenRightShift, ""},
	{true, "<", tokenLessThan, ""},
	{true, "<=", tokenLessOrEqual, ""},
	{true, ">", tokenGreaterThan, ""},
	{true, ">=", tokenGreaterOrEqual, ""},
	{true, "=", tokenEqual, ""},
	{true, "!=", tokenNotEqual, ""},
	{true, "&", tokenBitwiseAnd, ""},
	{true, "^", tokenBitwiseXor, ""},
	{true, "|", tokenBitwiseOr, ""},
	{true, "~", tokenBitwiseNot, "~"},
	{true, "&&", tokenLogicalAnd, ""},
	{true, "||", tokenLogicalOr, ""},

	{true, "x", tokenIdentifier, "x"},
	{true, "true", tokenTrue, ""},
	{true, "false", tokenFalse, ""},

	{false, "3a", tokenError, ""},
	{false, "0x", tokenError, ""},
	{false, "`", tokenError, ""},
}

var lexerMultipleTokenTests = []lexerMultipleTokenTest{
	{true, "1 2 3", types(tokenInt, tokenInt, tokenInt), vals("1", "2", "3")},
	{true, "(5)", types(tokenLeftParen, tokenInt, tokenRightParen), vals("", "5", "")},
	{true, "5 - (10 - 5)", types(tokenInt, tokenMinus, tokenLeftParen, tokenInt, tokenMinus, tokenInt, tokenRightParen),
		vals("5", "", "", "10", "", "5", "")},
	{true, "5*1/2", types(tokenInt, tokenStar, tokenInt, tokenSlash, tokenInt),
		vals("5", "", "1", "", "2")},

	{true, "f(x)", types(tokenIdentifier, tokenLeftParen, tokenIdentifier, tokenRightParen),
		vals("f", "", "x", "")},
}

func multipleTokenTest(test lexerSingleTokenTest) lexerMultipleTokenTest {
	return lexerMultipleTokenTest{test.ok, test.input, types(test.typ), vals(test.val)}
}

func TestLexSingleToken(t *testing.T) {
	for _, test := range lexerSingleTokenTests {
		if test.ok {
			shouldLex(multipleTokenTest(test), t)
		} else {
			lexShouldFail(test.input, t)
		}
	}
}

func TestLexMultipleTokens(t *testing.T) {
	for _, test := range lexerMultipleTokenTests {
		if test.ok {
			shouldLex(test, t)
		} else {
			lexShouldFail(test.input, t)
		}
	}
}

// Helpers

func shouldLex(test lexerMultipleTokenTest, t *testing.T) {
	s := test.input
	ts := test.types
	v := test.vals
	l := newLexer(s)
	for i, e := range ts {
		to := l.token()
		if to.typ != e {
			t.Errorf("Wrong token type: expected %v, got %v", e, to.typ)
		}
		if v != nil {
			ev := v[i]
			if ev != "" && to.val != ev {
				t.Errorf("Wrong token value: expected %v, got %v", ev, to.val)
			}
		}
	}
	if typ := l.token().typ; typ != tokenEOF {
		t.Errorf("Expected EOF, got %s", typ)
	}
}

func lexShouldFail(s string, t *testing.T) {
	l := newLexer(s)
	f := false
	ts := []*token{}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Lex did not fail, tokens: %v, input: %v", ts, s)
		}
	}()

	for {
		to := l.token()
		ts = append(ts, to)
		if to.typ == tokenError {
			f = true
			break
		}
	}

	if !f {
		t.Errorf("Lex did not fail, tokens: %v", ts)
	}
}

func types(args ...tokenType) []tokenType {
	return args
}

func vals(args ...string) []string {
	return args
}
