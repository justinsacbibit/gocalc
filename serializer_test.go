package gocalc

import (
	"os"
	"testing"
)

func TestSerializer(t *testing.T) {
	p := newParser("(1 + abs(-5)) > 1.0 + a")
	e := p.parseExpr()
	s := newSerializer()
	e.accept(s)
	s.serialize(os.Stdout)
}
