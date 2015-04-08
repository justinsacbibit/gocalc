package gocalc

import (
	_ "os"
	"testing"
)

type devNull struct{}

func (d devNull) Write([]byte) (int, error) {
	return 0, nil
}

func TestSerializer(t *testing.T) {
	p := newParser(newLexer("((1 + abs(-5)) > 1.0 + a) || (a > 2 && false)"))
	e := p.parseExpr()
	s := newSerializer()
	e.accept(s)
	s.serialize(devNull{})
}
