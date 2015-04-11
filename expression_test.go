package gocalc

import (
	"fmt"
	"testing"
)

type expressionTest struct {
	ok     bool
	expr   string
	expect interface{}
}

type resolverExpressionTest struct {
	expressionTest
	p ParamResolver
	f FuncHandler
}

var tests = []expressionTest{
	// Compilation errors
	{false, "@#$", nil},
	{false, "1<", nil},
	{false, "f(a a", nil},
	{false, "1.0a", nil},
	{false, "0d", nil},
	{false, "07c", nil},
	{false, "0bb", nil},
	{false, "0xabcg", nil},

	// Binary expressions

	// Logical or
	{true, "1 < 0 || 2 < 1", false},
	{true, "1 > 0 || 2 < 1", true},
	{true, "1 > 0 || 2 > 1", true},
	{false, "9 || 10", nil},
	{false, "3.5 || -1", nil},

	// Logical and
	{true, "true && false", false},
	{true, "true && true", true},
	{true, "1 > 2 && true", false},

	// Bitwise or
	{true, "1 | 2", 3},

	// Bitwise and
	{true, "1 & 2", 0},
	{true, "4 & 4", 4},

	// Bitwise xor
	{true, "1 ^ 2", 3},
	{true, "2 ^ 2", 0},
	{true, "3 ^ 2", 1},

	// Equal
	{true, "1=1", true},
	{true, "3=5", false},
	{true, "true=true", true},
	{true, "true=false", false},
	{true, "1.0=1.0", true},
	{true, "3.5=4.0", false},

	// Not equal
	{true, "1!=1", false},
	{true, "4 != 40", true},
	{true, "true != true", false},
	{true, "false != true", true},
	{true, "1.0 != 3.5", true},
	{true, "8.5 != 8.5", false},

	// Less than
	{true, "1 < 2", true},
	{true, "2 < 1", false},
	{true, "2.0 < 3", true},
	{true, "2 < 3.0", true},
	{true, "5.0 < 10.0", true},
	{true, "9.0 < 4.0", false},

	// Less or equal
	{true, "1 <= 1", true},
	{true, "1 <= 2", true},
	{true, "8 <= 3", false},
	{true, "8.0 <= 8", true},
	{true, "4 <= 4.0", true},
	{true, "3.0 <= 4.0", true},
	{true, "3.0 <= 3.0", true},
	{true, "3.0 <= 2.99999", false},

	// Greater than
	{true, "50 > 25", true},
	{true, "100 > 1000", false},
	{true, "50 > 25.0", true},
	{true, "25.0 > 50", false},
	{true, "3.0 > 2.0", true},

	// Greater or equal
	{true, "5 >= 3", true},
	{true, "5 >= 5", true},
	{true, "5 >= 6", false},
	{true, "5.0 >= 5", true},
	{true, "6.0 >= 1", true},
	{true, "1 >= -1.0", true},
	{true, "1 >= 1.0", true},
	{true, "1 >= 2.0", false},
	{true, "4.0 >= 5.0", false},
	{true, "8.0 >= 1.0", true},

	// Left shift
	{true, "1 << 1", 2},

	// Right shift
	{true, "1 >> 1", 0},

	// Plus
	{true, "9 + 10", 19},
	{true, "5.0 + 10", 15.0},
	{true, "1 + 4.0", 5.0},
	{true, "5.0 + 15.0", 20.0},

	// Minus
	{true, "9 - 10", -1},
	{true, "5.0 - 10", -5.0},
	{true, "1 - 4.0", -3.0},
	{true, "5.0 - 15.0", -10.0},

	// Star
	{true, "9 * 10", 90},
	{true, "5.0 * 10", 50.0},
	{true, "1 * 4.0", 4.0},
	{true, "5.0 * 15.0", 75.0},

	// Slash
	{true, "3 / 2", 1},
	{true, "5.0 / 10", 0.5},
	{true, "1 / 4.0", 0.25},
	{true, "4.0 / 2.0", 2.0},

	// Percent
	{true, "5 % 3", 2},

	// Unary
	{true, "-1", -1},
	{true, "!false", true},
	{true, "~2", -3},
	{false, "~2.0", nil},

	// Function
	{true, "abs(-2)", 2},
	{true, "abs(-1.0)", 1.0},
	{true, "abs(4)", 4},
	{true, "abs(9.0)", 9.0},
	{false, "g()", nil},

	// Identifiers
	{false, "d", nil},
}

var resolverTests = []resolverExpressionTest{
	{expressionTest{true, "a", 5}, func(s string) interface{} {
		if s == "a" {
			return int64(5)
		}

		return nil
	}, nil},

	{expressionTest{true, "abs(-3)", 3}, nil, func(f string, args ...func() interface{}) (interface{}, bool) {
		return nil, false
	}},

	{expressionTest{false, "abs()", nil}, nil, nil},

	{expressionTest{false, "f(2.0)", 3}, nil, func(f string, args ...func() interface{}) (interface{}, bool) {
		return args[0]().(int), true
	}},

	{expressionTest{true, "add(1, 2)", 3}, nil, func(f string, args ...func() interface{}) (interface{}, bool) {
		if f == "add" {
			if l := len(args); l != 2 {
				panic(EvaluationError(fmt.Sprintf("add takes two params, got %d", l)))
			}

			return args[0]().(int64) + args[1]().(int64), true
		}

		return nil, false
	}},
}

func allTests() []resolverExpressionTest {
	r := resolverTests
	for _, test := range tests {
		r = append(r, resolverExpressionTest{test, nil, nil})
	}
	return r
}

func TestExpressionEvaluation(t *testing.T) {
	for _, test := range allTests() {
		e, err := NewExpr(test.expr)
		if err != nil {
			if test.ok {
				t.Errorf("Expression \"%v\": Cannot test; lexer or parser error: %v", test.expr, err)
			} else {
				t.Logf("Logging expected compilation error: \"%v\"", err.Error())
			}

			continue
		}

		switch r := test.expect.(type) {
		case int:
			test.expect = int64(r)
		}

		res, err := e.Evaluate(test.p, test.f)
		if test.ok && err != nil {
			t.Errorf("Expression \"%v\": Evaluation failed with error: %v, expected result: %v",
				test.expr, err, test.expect)
		} else if test.ok && test.expect != res {
			t.Errorf("Expression \"%v\": Evaluation returned %v (%T), expected %v (%T)",
				test.expr, res, res, test.expect, test.expect)
		} else if !test.ok && err == nil {
			t.Errorf("Expression \"%v\": Evaluation passed but should have failed", test.expr)
		} else if !test.ok {
			t.Logf("Logging expected evaluation error: \"%v\"", err.Error())
		}
	}
}

func TestEvaluatorCountMallocs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping malloc count in short mode")
	}

	s := "((((1) + (2) - (3) & (4)) * (5) / (1.)) >= (2)) && ((((5) - (4) * (3)) / (2)) <= (1))"
	e, _ := NewExpr(s)
	mallocs := testing.AllocsPerRun(100, func() {
		e.Evaluate(nil, nil)
	})

	t.Logf("Expression \"%v\": got %v mallocs", s, mallocs)
}

func BenchmarkConstantExpressionCompilation(b *testing.B) {
	s := "((((1) + (2) - (3) & (4)) * (5) / (1.)) >= (2)) && ((((5) - (4) * (3)) / (2)) <= (1))"
	for i := 0; i < b.N; i++ {
		NewExpr(s)
	}
}

func BenchmarkConstantExpressionEvaluation(b *testing.B) {
	s := "((((1) + (2) - (3) & (4)) * (5) / (1.)) >= (2)) && ((((5) - (4) * (3)) / (2)) <= (1))"
	e, _ := NewExpr(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Evaluate(nil, nil)
	}
}

func BenchmarkParamExpressionEvaluation(b *testing.B) {
	s := "((((a) + (b) - (c) & (d)) * (e) / (1.)) >= (2)) && ((((5) - (4) * (3)) / (2)) <= (1))"
	m := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	r := func(p string) interface{} {
		return m[p]
	}
	e, _ := NewExpr(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Evaluate(r, nil)
	}
}

func Example_simple() {
	expression, _ := NewExpr("1 + 2")
	// if err != nil {
	// 	// Compilation of the expression failed
	// 	fmt.Println(err)
	// 	return
	// }

	result, _ := expression.Evaluate(nil, nil)
	// if err != nil {
	// 	// Evaluation of the expression failed
	// 	fmt.Println(err)
	// 	return
	// }

	fmt.Println(result)

	// Output:
	// 3
}

func Example_identifier() {
	expression, _ := NewExpr("a * 2")

	result, _ := expression.Evaluate(func(id string) interface{} {
		switch id {
		case "a":
			return 2
		}

		return nil
	}, nil)

	fmt.Println(result)

	// Output:
	// 4
}

func Example_function() {
	expression, _ := NewExpr("sub2(5)")

	result, _ := expression.Evaluate(nil, func(fn string, args ...func() interface{}) (result interface{}, handled bool) {
		switch fn {
		case "sub2":
			if len(args) != 1 {
				panic(EvaluationError("sub2 requires 1 argument"))
			}

			arg := args[0]().(int64)
			return arg - 2, true
		}

		return nil, false
	})

	fmt.Println(result)

	// Output:
	// 3
}

func Example_compileError() {
	_, err := NewExpr("1 +")
	if err != nil {
		// A compilation error has occurred
	}
}

func Example_evaluationError() {
	expression, _ := NewExpr("a")

	_, err := expression.Evaluate(nil, nil)
	if err != nil {
		// An evaluation error has occurred
	}
}
