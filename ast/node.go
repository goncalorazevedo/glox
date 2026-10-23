package ast

import (
	"fmt"

	"github.com/goncalorazevedo/glox/lexer"
)

type Expr interface{ isExpr() }

func (*Primary) isExpr()    {}
func (*Grouping) isExpr()   {}
func (*Unary) isExpr()      {}
func (*Binary) isExpr()     {}
func (*Logical) isExpr()    {}
func (*Variable) isExpr()   {}
func (*Assignment) isExpr() {}
func (*Call) isExpr()       {}

type ExprVisitor[R any] interface {
	VisitPrimary(e *Primary) (R, error)
	VisitGrouping(e *Grouping) (R, error)
	VisitUnary(e *Unary) (R, error)
	VisitBinary(e *Binary) (R, error)
	VisitLogical(e *Logical) (R, error)
	VisitVariable(e *Variable) (R, error)
	VisitAssignment(e *Assignment) (R, error)
	VisitCall(c *Call) (R, error)
}

func AcceptExpr[R any](e Expr, v ExprVisitor[R]) (R, error) {
	switch expr := e.(type) {
	case *Primary:
		return v.VisitPrimary(expr)
	case *Grouping:
		return v.VisitGrouping(expr)
	case *Unary:
		return v.VisitUnary(expr)
	case *Binary:
		return v.VisitBinary(expr)
	case *Logical:
		return v.VisitLogical(expr)
	case *Variable:
		return v.VisitVariable(expr)
	case *Assignment:
		return v.VisitAssignment(expr)
	case *Call:
		return v.VisitCall(expr)
	default:
		panic(fmt.Sprintf("missing expression case for %T", e))
	}
}

type Primary struct {
	Value Value
}

type Grouping struct {
	Expr Expr
}

type Unary struct {
	Token lexer.Token
	Expr  Expr
}

type Binary struct {
	Token lexer.Token
	Left  Expr
	Right Expr
}

type Logical struct {
	Token lexer.Token
	Left  Expr
	Right Expr
}

type Variable struct {
	Identifier lexer.Token
}

type Assignment struct {
	Identifier lexer.Token
	Expr       Expr
}

type Call struct {
	Callee    Expr
	Arguments []Expr
}
