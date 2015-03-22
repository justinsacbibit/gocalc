package gocalc

import (
	_ "fmt"
	"testing"
)

func TestParseNumber(t *testing.T) {
	p := newParser("5 + 5 + 5")
	e := p.parseExpr()
	print(e)
}
