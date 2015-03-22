package gocalc

import (
	"fmt"
)

type node interface {
	print(indent int)
}

type expr interface {
	node
}

type literal struct {
	val string
}

type binaryExpr struct {
	left  expr
	right expr
	//opPos int
	op token
}

func print(node node) {
	if node == nil {
		fmt.Println("nil")
	} else {
		node.print(0)
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
