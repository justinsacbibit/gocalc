package gocalc

import (
	"fmt"
	"testing"
)

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

func TestParseDoubleBinary(t *testing.T) {
	shouldFail("1++", t)
}

func TestParseDoubleUnaryPlus(t *testing.T) {
	shouldParse("++1", t)
}

func TestParseDoubleUnaryMinus(t *testing.T) {
	shouldParse("--1", t)
}

func TestParseTrue(t *testing.T) {
	shouldParse("true", t)
}

func TestParseFalse(t *testing.T) {
	shouldParse("false", t)
}

func TestParse(t *testing.T) {
	p := newParser("-----------1")
	e := p.parseExpr()
	if e != nil {
		// s := newSerializer()
		// e.accept(s)
		// s.serialize(os.Stdout)
	} else {
		fmt.Println(p.error)
	}
}
