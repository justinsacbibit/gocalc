package gocalc

import "fmt"

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF

	tokenIdentifier
	tokenNumber

	tokenComma

	tokenLeftParen
	tokenRightParen

	tokenPlus
	tokenMinus

	tokenStar
	tokenSlash
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
