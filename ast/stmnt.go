package ast

import (
	"fmt"

	"github.com/goncalorazevedo/glox/lexer"
)

type Stmnt interface{ isStmnt() }

func (*ExprStmnt) isStmnt()  {}
func (*PrintSmnt) isStmnt()  {}
func (*VarDecl) isStmnt()    {}
func (*Block) isStmnt()      {}
func (*IfStmnt) isStmnt()    {}
func (*WhileStmnt) isStmnt() {}
func (*FunDecl) isStmnt()    {}
func (*Return) isStmnt()     {}

type StmntVisitor interface {
	VisitExprStmnt(*ExprStmnt) error
	VisitPrintStmnt(*PrintSmnt) error
	VisitVarDecl(*VarDecl) error
	VisitBlock(*Block) error
	VisitIfStmnt(*IfStmnt) error
	VisitWhileStmnt(*WhileStmnt) error
	VisitFunDecl(*FunDecl) error
	VisitReturn(*Return) error
}

func AcceptStmnt(s Stmnt, v StmntVisitor) error {
	switch stmnt := s.(type) {
	case *ExprStmnt:
		return v.VisitExprStmnt(stmnt)
	case *PrintSmnt:
		return v.VisitPrintStmnt(stmnt)
	case *VarDecl:
		return v.VisitVarDecl(stmnt)
	case *Block:
		return v.VisitBlock(stmnt)
	case *IfStmnt:
		return v.VisitIfStmnt(stmnt)
	case *WhileStmnt:
		return v.VisitWhileStmnt(stmnt)
	case *FunDecl:
		return v.VisitFunDecl(stmnt)
	case *Return:
		return v.VisitReturn(stmnt)
	default:
		panic(fmt.Sprintf("missing statement case for %T", s))
	}
}

type ExprStmnt struct {
	Expr Expr
}

type PrintSmnt struct {
	Expr Expr
}

type VarDecl struct {
	Identifier lexer.Token
	Init       Expr
}

type Block struct {
	Stmnts []Stmnt
}

type IfStmnt struct {
	Condition Expr
	Then      Stmnt
	ElseThen  Stmnt
}

type WhileStmnt struct {
	Condition Expr
	Stmnt     Stmnt
}

type FunDecl struct {
	Identifier lexer.Token
	Params     []lexer.Token
	Body       Stmnt
}

type Return struct {
	Expr Expr
}
