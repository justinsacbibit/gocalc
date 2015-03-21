package gocalc

import (
	_ "fmt"
	"strings"
	"testing"
)

func TestLexDigit(t *testing.T) {
	s := "1"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenNumber, t)
	checkVal(to, s, t)
}

func TestLexMultiDigit(t *testing.T) {
	s := "111"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenNumber, t)
	checkVal(to, s, t)
}

func TestLexMultipleNumbers(t *testing.T) {
	s := "1 2 3"
	l := newLexer(s)

	strs := strings.Split(s, " ")
	for _, str := range strs {
		to := l.token()
		checkType(to, tokenNumber, t)
		checkVal(to, str, t)
	}
}

func TestLexNumberWithLeadingWhitespace(t *testing.T) {
	s := "5"
	l := newLexer("  " + s)
	to := l.token()

	checkType(to, tokenNumber, t)
	checkVal(to, s, t)
}

func TestLexNumberWithSurroundingWhitespace(t *testing.T) {
	s := "99"
	l := newLexer("  " + s + "  ")
	to := l.token()

	checkType(to, tokenNumber, t)
	checkVal(to, s, t)
}

func TestLexPlus(t *testing.T) {
	s := "+"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenPlus, t)
}

func TestLexMinus(t *testing.T) {
	s := "-"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenMinus, t)
}

func TestLexLeftParen(t *testing.T) {
	s := "("
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenLeftParen, t)
}

func TestLexRightParen(t *testing.T) {
	s := ")"
	l := newLexer(s)
	to := l.token()

	checkType(to, tokenRightParen, t)
}

func checkVal(to token, exp string, t *testing.T) {
	if to.val != exp {
		t.Errorf("Wrong token value: %s", to.val)
	}
}

func checkType(to token, exp tokenType, t *testing.T) {
	if to.typ != exp {
		t.Errorf("Wrong token type: %d", to.typ)
	}
}

/*
func TestLexMultipleNumbers(t *testing.T) {
	str := " 1 2 3 44 123 "
	errs := testLex(t, str, NumberToken, []string{"1", "2", "3", "44", "123"})
	if errs != nil {
		for _, err := range errs {
			t.Errorf("Error: %s", err.Error())
		}
	}
}

func TestLexWhitespace(t *testing.T) {
	strs := []string{"", "  "}
	for _, str := range strs {
		err := testLexSingle(t, str, 0, "")
		if err == nil {
			t.Error("Did not error")
		}
	}
}

func TestLexPlus(t *testing.T) {
	str := "+"
	err := testLexSingle(t, str, PlusToken, str)
	if err != nil {
		t.Error(err.Error())
	}
}

func testLexSingle(t *testing.T, sequence string, expectedType int, expectedValue string) error {
	l := NewLexer([]byte(sequence))
	token, err := l.Token()

	if err != nil {
		return err
	}

	if token.Type != expectedType {
		t.Errorf("Wrong token type: %d", token.Type)
	}

	if token.Value == "" {
		t.Error("Token value is empty")
	} else if token.Value != expectedValue {
		t.Errorf("Wrong token value: %s", token.Value)
	}
	return nil
}

func testLex(t *testing.T, sequence string, expectedType int, expectedValues []string) []error {
	l := NewLexer([]byte(sequence))
	errs := []error{}
	for _, expectedValue := range expectedValues {
		token, err := l.Token()

		if err != nil {
			errs = append(errs, err)
			continue
		}

		if token.Type != expectedType {
			t.Errorf("Wrong token type: %d", token.Type)
		}

		if token.Value == "" {
			t.Error("Token value is empty")
		} else if token.Value != expectedValue {
			t.Errorf("Wrong token value: %sz", token.Value)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}*/
