package parser

import (
	"github.com/goncalorazevedo/glox/ast"
	"github.com/goncalorazevedo/glox/lexer"
)

type parser struct {
	tokens []lexer.Token
	index  int
}

func Parse(tokens []lexer.Token) ([]ast.Stmnt, error) {
	p := parser{
		tokens: tokens,
		index:  0,
	}
	return p.parse()
}

func ParseExpr(tokens []lexer.Token) (ast.Expr, error) {
	p := parser{
		tokens: tokens,
		index:  0,
	}
	return p.expression()
}
