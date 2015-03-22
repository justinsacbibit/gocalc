package gocalc

import (
	"fmt"
	"testing"
)

func TestParseNumber(t *testing.T) {
	p := newParser("+5-(5+5)")
	e := p.parseExpr()
	if e != nil {
		print(e)
	} else {
		fmt.Println(p.error)
	}
}
