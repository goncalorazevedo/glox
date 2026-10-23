package lox

import (
	"fmt"

	"github.com/goncalorazevedo/glox/ast"
)

type environment struct {
	values    map[string]ast.Value
	enclosing *environment
}

func newEnvironment(enclosing *environment) *environment {
	return &environment{
		values:    make(map[string]ast.Value),
		enclosing: enclosing,
	}
}

func (e *environment) define(s string, v ast.Value) error {
	if _, exists := e.values[s]; exists {
		return fmt.Errorf("variable %s already defined in this scope", s)
	}
	e.values[s] = v
	return nil
}

func (e *environment) assign(s string, v ast.Value) error {
	if _, exists := e.values[s]; !exists {
		return e.enclosing.assign(s, v)
	}
	e.values[s] = v
	return nil
}

func (e *environment) lookup(s string) ast.Value {
	if v, exists := e.values[s]; exists {
		return v
	}
	if e.enclosing != nil {
		return e.enclosing.lookup(s)
	}
	return nil
}
