package gocalc

import "fmt"

type printerVisitor struct {
	indent int
	ignore bool
}

var indentStr = ".   "

func newPrinter() *printerVisitor {
	return &printerVisitor{}
}

func (p *printerVisitor) printf(format string, args ...interface{}) {
	p.printIndent()
	fmt.Printf(format, args...)
}

func (p *printerVisitor) println(s string) {
	p.printIndent()
	fmt.Println(s)
}

func (p *printerVisitor) visitBinaryExpr(b *binaryExpr) {
	p.println("*binaryExpr {")
	p.indent++
	p.printf("lhs: ")
	p.ignore = true
	b.left.accept(p)
	p.printf("op: %s\n", b.op.val)
	p.printf("rhs: ")
	p.ignore = true
	b.right.accept(p)
	p.indent--
	p.println("}")
}

func (p *printerVisitor) visitFuncExpr(f *funcExpr) {
	p.println("*funcExpr {")
	p.indent++
	p.printf("func: %s\n", f.function)
	p.printf("args (len: %d) {\n", len(f.args))
	p.indent++
	for _, arg := range f.args {
		arg.accept(p)
	}
	p.indent--
	p.println("}")
	p.indent--
	p.println("}")
}

func (p *printerVisitor) visitUnaryExpr(u *unaryExpr) {
	p.println("*unaryExpr {")
	p.indent++
	p.printf("op: %s\n", u.op.val)
	p.printf("expr: ")
	p.ignore = true
	u.expr.accept(p)
	p.indent--
	p.println("}")
}

func (p *printerVisitor) visitBoolExpr(b *boolExpr) {
	p.println("*boolExpr {")
	p.indent++
	p.printf("val: %t\n", b.val)
	p.indent--
	p.println("}")
}

func (p *printerVisitor) visitFloatExpr(f *floatExpr) {
	p.println("*floatExpr {")
	p.indent++
	p.printf("val: %f\n", f.val)
	p.indent--
	p.println("}")
}

func (p *printerVisitor) visitIntExpr(i *intExpr) {
	p.println("*intExpr {")
	p.indent++
	p.printf("val: %d\n", i.val)
	p.indent--
	p.println("}")
}

func (p *printerVisitor) visitParamExpr(e *paramExpr) {
	p.println("*identifier {")
	p.indent++
	p.printf("val: \"%s\"\n", e.identifier)
	p.indent--
	p.println("}")
}

func (p *printerVisitor) printIndent() {
	if p.ignore {
		p.ignore = false
		return
	}
	for i := 0; i < p.indent; i++ {
		fmt.Print(indentStr)
	}
}
