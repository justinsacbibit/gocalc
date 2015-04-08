package gocalc

import (
	"fmt"
	"testing"
)

var parserTests = []struct {
	ok   bool
	expr string
}{
	// Values
	{true, "3"},
	{true, "5.5"},
	{true, "true"},
	{true, "false"},
	{false, "0x"},
	{false, "0a"},

	// Functions
	{true, "f()"},
	{true, "f(x)"},
	{true, "f(x, y)"},
	{false, "f("},
	{false, "f(,)"},
	{false, "f(1,)"},
	{false, "g(a,)"},
	{false, "f("},
	{false, "f(1,"},

	// Unary
	{true, "-9"},
	{true, "++1"},
	{true, "--1"},
	{true, "!true"},
	{true, "~3"},
	{false, "1++"},

	// Binary
	{true, "a || b"},
	{true, "a && b"},
	{true, "2 | 4"},
	{true, "9 ^ 10"},
	{true, "4 & 4"},
	{true, "4 = 5"},
	{true, "9 != 2"},
	{true, "2 < 4"},
	{true, "1 <= 0"},
	{true, "5 > 8"},
	{true, "53>=1"},
	{true, "a>>b"},
	{true, "b<<a"},
	{true, "9 + 10"},
	{true, "4 - 7"},
	{true, "4 * 7 "},
	{true, "1 / 0"},
	{true, "3 % 2"},
	{false, "1 + 1 1"},
	{false, "4 5 + 1"},

	// Parenthesized
	{true, "(7)"},
	{false, "(1 + 2"},
}

func TestParse(t *testing.T) {
	for _, test := range parserTests {
		if test.ok {
			shouldParse(test.expr, t)
		} else {
			shouldFail(test.expr, t)
		}
	}
}

type mockLexer struct {
	cur    int
	tokens []*token
}

func (m *mockLexer) token() *token {
	t := m.tokens[m.cur]
	m.cur++
	return t
}

func (m *mockLexer) peekToken() *token {
	return m.tokens[m.cur]
}

func BenchmarkParseConstantExpression(b *testing.B) {
	lp := &token{tokenLeftParen, "("}
	rp := &token{tokenRightParen, ")"}
	in := func(i int) *token {
		return &token{tokenInt, fmt.Sprintf("%d", i)}
	}
	pl := &token{tokenPlus, "+"}
	mi := &token{tokenMinus, "-"}
	st := &token{tokenStar, "*"}
	sl := &token{tokenSlash, "/"}
	fl := func(f float64) *token {
		return &token{tokenFloat, fmt.Sprintf("%f", f)}
	}
	ba := &token{tokenBitwiseAnd, "&"}
	ge := &token{tokenGreaterOrEqual, ">="}
	la := &token{tokenLogicalAnd, "&&"}
	le := &token{tokenLessOrEqual, "<="}
	eo := &token{tokenEOF, ""}
	// s := "((((1) + (2) - (3) & (4)) * (5) / (1.)) >= (2)) && ((((5) - (4) * (3)) / (2)) <= (1))"
	l := &mockLexer{tokens: []*token{
		lp, lp, lp, lp, in(1), rp, pl, lp, in(2), rp, mi, lp, in(3), rp, ba, lp, in(4), rp, rp, st, lp, in(5),
		rp, sl, lp, fl(1.0), rp, rp, ge, lp, in(2), rp, rp, la, lp, lp, lp, lp, in(5), rp, mi, lp, in(4), rp,
		st, lp, in(3), rp, rp, sl, lp, in(2), rp, rp, le, lp, in(1), rp, rp, eo,
	}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := newParser(l)
		p.parseExpr()
		l.cur = 0
	}
}

func shouldParse(s string, t *testing.T) {
	p := newParser(newLexer(s))
	if e := p.parseExpr(); e == nil {
		t.Fatalf("Parse of \"%s\" failed: %s", s, p.error)
	}
}

func shouldFail(s string, t *testing.T) {
	p := newParser(newLexer(s))
	if e := p.parseExpr(); e != nil {
		t.Fatalf("Parse of %s passed but should have failed.", s)
	}
}

// func TestParse(t *testing.T) {
// 	p := newParser("-----------1")
// 	e := p.parseExpr()
// 	if e != nil {
// 		s := newSerializer()
// 		e.accept(s)
// 		s.serialize(os.Stdout)
// 	} else {
// 		fmt.Println(p.error)
// 	}
// }
