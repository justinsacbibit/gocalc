package gocalc

import "testing"

type mockExprVisitor struct {
	visited []expr
}

func (m *mockExprVisitor) add(e expr) {
	m.visited = append(m.visited, e)
}

func (m *mockExprVisitor) visitBinaryExpr(b *binaryExpr) { m.add(b) }
func (m *mockExprVisitor) visitFuncExpr(f *funcExpr)     { m.add(f) }
func (m *mockExprVisitor) visitParamExpr(p *paramExpr)   { m.add(p) }
func (m *mockExprVisitor) visitUnaryExpr(u *unaryExpr)   { m.add(u) }
func (m *mockExprVisitor) visitBoolExpr(b *boolExpr)     { m.add(b) }
func (m *mockExprVisitor) visitFloatExpr(f *floatExpr)   { m.add(f) }
func (m *mockExprVisitor) visitIntExpr(i *intExpr)       { m.add(i) }

func newMockExprVisitor() *mockExprVisitor {
	return &mockExprVisitor{}
}

var acceptExprs = []expr{
	&binaryExpr{},
	&funcExpr{},
	&paramExpr{},
	&unaryExpr{},
	&boolExpr{},
	&floatExpr{},
	&intExpr{},
}

func TestAccepts(t *testing.T) {
	for _, accepter := range acceptExprs {
		m := newMockExprVisitor()
		accepter.accept(m)

		if len(m.visited) > 1 {
			t.Errorf("More than one expr visited: %#v", m.visited)
		} else if len(m.visited) < 1 {
			t.Error("No expr visited")
		} else if visited := m.visited[0]; visited != accepter {
			t.Errorf("Wrong expr was visited: %T", visited)
		}
	}
}
