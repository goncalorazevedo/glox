package lox

import (
	"bytes"
	"fmt"

	"github.com/goncalorazevedo/glox/ast"
)

type astPrinter struct{}

var _ ast.ExprVisitor[string, Value] = (*astPrinter)(nil)

func PrintExpr(expr ast.Expr) (string, error) {
	p := &astPrinter{}
	return ast.AcceptExpr(expr, p)
}

func (p *astPrinter) VisitPrimary(b *ast.Primary[Value]) (string, error) {
	return fmt.Sprintf("%s", b.Value.String()), nil
}

func (p *astPrinter) VisitGrouping(g *ast.Grouping) (string, error) {
	e, err := ast.AcceptExpr(g.Expr, p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("(%s)", e), nil
}

func (p *astPrinter) VisitUnary(u *ast.Unary) (string, error) {
	e, err := ast.AcceptExpr(u.Expr, p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", u.Token.Lexeme, e), nil
}

func (p *astPrinter) VisitBinary(b *ast.Binary) (string, error) {
	left, err := ast.AcceptExpr(b.Left, p)
	if err != nil {
		return "", err
	}
	right, err := ast.AcceptExpr(b.Right, p)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(%s %s %s)", b.Token.Lexeme, left, right), nil
}

func (p *astPrinter) VisitLogical(l *ast.Logical) (string, error) {
	left, err := ast.AcceptExpr(l.Left, p)
	if err != nil {
		return "", err
	}
	right, err := ast.AcceptExpr(l.Right, p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("(%s %s %s)", l.Token.Lexeme, left, right), nil
}

func (p *astPrinter) VisitVariable(v *ast.Variable) (string, error) {
	return v.Identifier.Lexeme, nil
}

func (p *astPrinter) VisitAssignment(a *ast.Assignment) (string, error) {
	expr, err := ast.AcceptExpr(a.Expr, p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("(= %s %s)", a.Identifier.Lexeme, expr), nil
}

func (p *astPrinter) VisitCall(c *ast.Call) (string, error) {
	var buffer bytes.Buffer
	for _, a := range c.Arguments {
		s, err := ast.AcceptExpr(a, p)
		if err != nil {
			return "", err
		}
		buffer.WriteString(s)
	}
	callee, err := ast.AcceptExpr(c.Callee, p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("(call %s (%s))", callee, buffer.String()), nil
}
