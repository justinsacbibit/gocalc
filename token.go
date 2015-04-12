package gocalc

//go:generate stringer -type=tokenType

type tokenType int

const (
	tokenError tokenType = iota
	tokenWhitespace
	tokenEOF

	tokenIdentifier
	tokenTrue
	tokenFalse

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
	pos int // starting position of token
	end int // ending position of token
}

func (t token) String() string {
	switch t.typ {
	case tokenEOF:
		return "EOF"
	}
	return t.val
}
