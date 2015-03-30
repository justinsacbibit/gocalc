package gocalc

import "fmt"

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF

	tokenIdentifier
	tokenInt
	tokenFloat

	tokenLeftParen
	tokenRightParen

	tokenComma

	tokenLogicalNot
	tokenBitwiseNot

	tokenBinary

	tokenStar
	tokenSlash
	tokenPercent

	tokenPlus
	tokenMinus

	tokenLeftShift
	tokenRightShift

	tokenLessThan
	tokenLessOrEqual
	tokenGreaterThan
	tokenGreaterOrEqual

	tokenEqual
	tokenNotEqual

	tokenBitwiseAnd

	tokenBitwiseXor

	tokenBitwiseOr

	tokenLogicalAnd

	tokenLogicalOr
)

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	switch t.typ {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return t.val
	}
	return fmt.Sprintf("%s", t.val)
}
