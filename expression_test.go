package gocalc

import "testing"

type testFixture struct {
	ok     bool
	expr   string
	expect interface{}
}

var logicalOrTests = []testFixture{
	{true, "1 < 0 || 2 < 1", false},
	{true, "1 > 0 || 2 < 1", true},
	{true, "1 > 0 || 2 > 1", true},
	{false, "9 || 10", nil},
	{false, "3.5 || -1", nil},
}

var logicalAndTests = []testFixture{}

var tests = [][]testFixture{
	logicalOrTests,
	logicalAndTests,
}

func TestExpression(t *testing.T) {
	for _, testGroup := range tests {
		for _, fixture := range testGroup {
			e, err := NewExpr(fixture.expr)
			if err != nil {
				t.Errorf("Cannot test expression; lexer or parser error: %v", err)
				continue
			}

			res, err := e.Evaluate(nil, nil)
			if fixture.ok && err != nil {
				t.Errorf("Expression evaluation failed with error: %v, expected result: %v",
					err, fixture.expect)
			} else if fixture.ok && fixture.expect != res {
				t.Errorf("Expression evaluation returned %v, expected %v",
					res, fixture.expect)
			} else if !fixture.ok && err == nil {
				t.Errorf("Expression evaluation passed but should have failed")
			}
		}
	}
}
