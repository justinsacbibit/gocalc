package gocalc

import (
	"fmt"
	"io"
)

type serializer struct {
	buffer buffer
	indent int
	ignore bool
}

func newSerializer() *serializer {
	return &serializer{}
}

func (s *serializer) serialize(w io.Writer) {
	w.Write(s.buffer)
}

type buffer []byte

func (b *buffer) Write(p []byte) (n int, err error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (s *serializer) printf(format string, args ...interface{}) {
	s.printIndent()
	fmt.Fprintf(&s.buffer, format, args...)
}

func (s *serializer) println(str string) {
	s.printIndent()
	fmt.Fprintln(&s.buffer, str)
}

func (s *serializer) visitBinaryExpr(b *binaryExpr) {
	s.println("*binaryExpr {")
	s.indent++
	s.printf("lhs: ")
	s.ignore = true
	b.left.accept(s)
	s.printf("op: %s\n", b.op.val)
	s.printf("rhs: ")
	s.ignore = true
	b.right.accept(s)
	s.indent--
	s.println("}")
}

func (s *serializer) visitFuncExpr(f *funcExpr) {
	s.println("*funcExpr {")
	s.indent++
	s.printf("func: %s\n", f.function)
	s.printf("args (len: %d) {\n", len(f.args))
	s.indent++
	for _, arg := range f.args {
		arg.accept(s)
	}
	s.indent--
	s.println("}")
	s.indent--
	s.println("}")
}

func (s *serializer) visitUnaryExpr(u *unaryExpr) {
	s.println("*unaryExpr {")
	s.indent++
	s.printf("op: %v\n", u.op.val)
	s.printf("expr: ")
	s.ignore = true
	u.expr.accept(s)
	s.indent--
	s.println("}")
}

func (s *serializer) visitBoolExpr(b *boolExpr) {
	s.println("*boolExpr {")
	s.indent++
	s.printf("val: %t\n", b.val)
	s.indent--
	s.println("}")
}

func (s *serializer) visitFloatExpr(f *floatExpr) {
	s.println("*floatExpr {")
	s.indent++
	s.printf("val: %f\n", f.val)
	s.indent--
	s.println("}")
}

func (s *serializer) visitIntExpr(i *intExpr) {
	s.println("*intExpr {")
	s.indent++
	s.printf("val: %d\n", i.val)
	s.indent--
	s.println("}")
}

func (s *serializer) visitParamExpr(e *paramExpr) {
	s.println("*identifier {")
	s.indent++
	s.printf("val: \"%s\"\n", e.identifier)
	s.indent--
	s.println("}")
}

var indent = []byte(".   ")

func (s *serializer) printIndent() {
	if s.ignore {
		s.ignore = false
		return
	}
	for i := 0; i < s.indent; i++ {
		s.buffer.Write(indent)
	}
}
