package gocalc

import (
	"fmt"
	"testing"
)

func shouldParse(s string, t *testing.T) {
	p := newParser(s)
	if e := p.parseExpr(); e == nil {
		t.Fatalf("Parse of \"%s\" failed: %s", p.error)
	}
}

func shouldFail(s string, t *testing.T) {
	p := newParser(s)
	if e := p.parseExpr(); e != nil {
		t.Fatalf("Parse of %s passed but should have failed.", s)
	}
}

func TestParseNumber(t *testing.T) {
	shouldParse("3", t)
}

func TestParseBadFunction(t *testing.T) {
	shouldFail("f(", t)
}

func TestParseBadFunctionComma(t *testing.T) {
	shouldFail("f(,)", t)
}

func TestParseBadFunctionExprComma(t *testing.T) {
	shouldFail("f(1,)", t)
}

func TestParseSimpleFunction(t *testing.T) {
	shouldParse("f()", t)
}

func TestParseFunctionSingleArgument(t *testing.T) {
	shouldParse("f(x)", t)
}

func TestParseFunctionMultipleArguments(t *testing.T) {
	shouldParse("f(x, y)", t)
}

func TestParse(t *testing.T) {
	p := newParser("abs(-3)")
	e := p.parseExpr()
	if e != nil {
		print(e)
	} else {
		fmt.Println(p.error)
	}
}
