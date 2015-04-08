package gocalc

import "testing"

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
	cur    int // index of current token
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

func (m *mockLexer) reset() {
	m.cur = 0
}

func newMockLexer(expr string) *mockLexer {
	l := &mockLexer{}

	// Use gocalcLexer to fill mockLexer's token buffer
	rl := newLexer(expr)
	for {
		t := rl.token()
		l.tokens = append(l.tokens, t)
		if t.typ == tokenEOF {
			break
		}
	}
	return l
}

func BenchmarkParseConstantExpression(b *testing.B) {
	l := newMockLexer("((((1) + (2) - (3) & (4)) * (5) / (1.)) >= (2)) && ((((5) - (4) * (3)) / (2)) <= (1))")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := newParser(l)
		p.parseExpr()
		l.reset()
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
