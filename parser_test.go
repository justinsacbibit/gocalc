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

func shouldParse(s string, t *testing.T) {
	p := newParser(s)
	if e := p.parseExpr(); e == nil {
		t.Fatalf("Parse of \"%s\" failed: %s", s, p.error)
	}
}

func shouldFail(s string, t *testing.T) {
	p := newParser(s)
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
