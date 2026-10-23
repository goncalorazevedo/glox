package ast

import (
	"bytes"
	"fmt"
)

type astPrinter struct{}

var _ ExprVisitor[string] = (*astPrinter)(nil)

func PrintExpr(expr Expr) (string, error) {
	p := &astPrinter{}
	return AcceptExpr(expr, p)
}

func (p *astPrinter) VisitPrimary(b *Primary) (string, error) {
	return fmt.Sprintf("%s", b.Value.String()), nil
}

func (p *astPrinter) VisitGrouping(g *Grouping) (string, error) {
	e, err := AcceptExpr(g.Expr, p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("(%s)", e), nil
}

func (p *astPrinter) VisitUnary(u *Unary) (string, error) {
	e, err := AcceptExpr(u.Expr, p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", u.Token.Lexeme, e), nil
}

func (p *astPrinter) VisitBinary(b *Binary) (string, error) {
	left, err := AcceptExpr(b.Left, p)
	if err != nil {
		return "", err
	}
	right, err := AcceptExpr(b.Right, p)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(%s %s %s)", b.Token.Lexeme, left, right), nil
}

func (p *astPrinter) VisitLogical(l *Logical) (string, error) {
	left, err := AcceptExpr(l.Left, p)
	if err != nil {
		return "", err
	}
	right, err := AcceptExpr(l.Right, p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("(%s %s %s)", l.Token.Lexeme, left, right), nil
}

func (p *astPrinter) VisitVariable(v *Variable) (string, error) {
	return v.Identifier.Lexeme, nil
}

func (p *astPrinter) VisitAssignment(a *Assignment) (string, error) {
	expr, err := AcceptExpr(a.Expr, p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("(= %s %s)", a.Identifier.Lexeme, expr), nil
}

func (p *astPrinter) VisitCall(c *Call) (string, error) {
	var buffer bytes.Buffer
	for _, a := range c.Arguments {
		s, err := AcceptExpr(a, p)
		if err != nil {
			return "", err
		}
		buffer.WriteString(s)
	}
	callee, err := AcceptExpr(c.Callee, p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("(call %s (%s))", callee, buffer.String()), nil
}
