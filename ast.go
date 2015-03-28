package gocalc

import "fmt"

type expr interface {
	print(indent int)
}

type literal struct {
	val string
}

type callExpr struct {
	function string
	args     []expr
}

type identifier struct {
	val string
}

type unaryExpr struct {
	expr expr
	op   *token
}

type binaryExpr struct {
	left  expr
	right expr
	//opPos int
	op *token
}

func print(expr expr) {
	if expr == nil {
		fmt.Println("nil")
	} else {
		expr.print(0)
	}
}

var indentStr = ".   "

func printIndent(indent int) {
	for indent > 0 {
		fmt.Print(indentStr)
		indent--
	}
}

func (lit *literal) print(indent int) {
	fmt.Print("*literal {\n")
	printIndent(indent + 1)
	fmt.Printf("val: \"%s\"\n", lit.val)
	printIndent(indent)
	fmt.Print("}\n")
}

func (iden *identifier) print(indent int) {
	fmt.Print("*identifier {\n")
	printIndent(indent + 1)
	fmt.Printf("val: \"%s\"\n", iden.val)
	printIndent(indent)
	fmt.Print("}\n")
}

func (expr *callExpr) print(indent int) {
	fmt.Print("*callExpr {\n")
	printIndent(indent + 1)
	fmt.Printf("func: %s\n", expr.function)
	printIndent(indent + 1)
	fmt.Printf("args (len: %d):\n", len(expr.args))
	for _, arg := range expr.args {
		printIndent(indent + 1)
		arg.print(indent + 1)
	}
	printIndent(indent)
	fmt.Print("}\n")
}

func (expr *unaryExpr) print(indent int) {
	fmt.Print("*unaryExpr {\n")
	printIndent(indent + 1)
	fmt.Printf("op: %s\n", expr.op.val)
	printIndent(indent + 1)
	fmt.Print("expr: ")
	expr.expr.print(indent + 1)
	printIndent(indent)
	fmt.Print("}\n")
}

func (expr *binaryExpr) print(indent int) {
	fmt.Print("*binaryExpr {\n")
	printIndent(indent + 1)
	fmt.Print("lhs: ")
	expr.left.print(indent + 1)
	printIndent(indent + 1)
	fmt.Printf("op: %s\n", expr.op.val)
	printIndent(indent + 1)
	fmt.Print("rhs: ")
	expr.right.print(indent + 1)
	printIndent(indent)
	fmt.Print("}\n")
}
