package gocalc

import "fmt"

type printer struct {
	indent int
	ignore bool
}

var indentStr = ".   "

func newPrinter() *printer {
	return &printer{}
}

func (p *printer) printf(format string, args ...interface{}) {
	p.printIndent()
	fmt.Printf(format, args...)
}

func (p *printer) println(s string) {
	p.printIndent()
	fmt.Println(s)
}

func (p *printer) visitBinaryExpr(b *binaryExpr) {
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

func (p *printer) visitFuncExpr(f *funcExpr) {
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

func (p *printer) visitUnaryExpr(u *unaryExpr) {
	p.println("*unaryExpr {")
	p.indent++
	p.printf("op: %s\n", u.op.val)
	p.printf("expr: ")
	p.ignore = true
	u.expr.accept(p)
	p.indent--
	p.println("}")
}

func (p *printer) visitValueExpr(v *valueExpr) {
	p.println("*valueExpr {")
	p.indent++
	p.printf("val: \"%s\"\n", v.val)
	p.indent--
	p.println("}")
}

func (p *printer) visitParamExpr(e *paramExpr) {
	p.println("*identifier {")
	p.indent++
	p.printf("val: \"%s\"\n", e.identifier)
	p.indent--
	p.println("}")
}

func (p *printer) printIndent() {
	if p.ignore {
		p.ignore = false
		return
	}
	for i := 0; i < p.indent; i++ {
		fmt.Print(indentStr)
	}
}
