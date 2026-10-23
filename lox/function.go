package lox

import (
	"fmt"

	"github.com/goncalorazevedo/glox/ast"
)

type Function struct {
	Name   string
	Params []string
	Stmnt  ast.Stmnt
}

var _ ast.Callable = (*Function)(nil)

func (f *Function) String() string {
	return fmt.Sprintf("<fun %s(%d)>", f.Name, f.Arity())
}

func (f *Function) Arity() int {
	return len(f.Params)
}

func (f *Function) Call(ctx any, args ...ast.Value) (ast.Value, error) {
	interp, ok := ctx.(*interpreter)
	if !ok {
		panic(fmt.Sprintf("missing Lox interpreter context in call to %s", f.String()))
	}
	if len(f.Params) != len(args) {
		panic(fmt.Sprintf("invalid number of arguments provided to %s", f.String()))
	}

	interp.env = newEnvironment(interp.env)
	for i, p := range f.Params {
		interp.env.define(p, args[i])
	}
	err := interp.evaluateStmnt(f.Stmnt)
	interp.env = interp.env.enclosing
	return &ast.Nil{}, err
}
