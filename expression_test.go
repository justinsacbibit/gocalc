package gocalc

import "testing"

type expressionTest struct {
	ok     bool
	expr   string
	expect interface{}
}

var tests = []expressionTest{
	// Compilation errors
	{false, "@#$", nil},
	{false, "1<", nil},

	// Logical or
	{true, "1 < 0 || 2 < 1", false},
	{true, "1 > 0 || 2 < 1", true},
	{true, "1 > 0 || 2 > 1", true},
	{false, "9 || 10", nil},
	{false, "3.5 || -1", nil},

	// Unary
	{true, "-1", -1},
}

func TestExpression(t *testing.T) {
	for _, test := range tests {
		e, err := NewExpr(test.expr)
		if err != nil {
			if test.ok {
				t.Errorf("Expression \"%v\": Cannot test; lexer or parser error: %v", test.expr, err)
			} else {
				t.Log(err.Error())
			}

			continue
		}

		switch r := test.expect.(type) {
		case int:
			test.expect = int64(r)
		}

		res, err := e.Evaluate(nil, nil)
		if test.ok && err != nil {
			t.Errorf("Expression \"%v\": Evaluation failed with error: %v, expected result: %v",
				test.expr, err, test.expect)
		} else if test.ok && test.expect != res {
			t.Errorf("Expression \"%v\": Evaluation returned %v (%T), expected %v (%T)",
				test.expr, res, res, test.expect, test.expect)
		} else if !test.ok && err == nil {
			t.Errorf("Expression \"%v\": Evaluation passed but should have failed", test.expr)
		}
	}
}
