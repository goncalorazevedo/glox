package lox

import (
	"fmt"
)

type environment struct {
	values    map[string]Value
	enclosing *environment
}

func newEnvironment(enclosing *environment) *environment {
	return &environment{
		values:    make(map[string]Value),
		enclosing: enclosing,
	}
}

func (e *environment) define(s string, v Value) error {
	if _, exists := e.values[s]; exists {
		return fmt.Errorf("variable %s already defined in this scope", s)
	}
	e.values[s] = v
	return nil
}

func (e *environment) assign(s string, v Value) error {
	if _, exists := e.values[s]; !exists {
		return e.enclosing.assign(s, v)
	}
	e.values[s] = v
	return nil
}

func (e *environment) lookup(s string) Value {
	if v, exists := e.values[s]; exists {
		return v
	}
	if e.enclosing != nil {
		return e.enclosing.lookup(s)
	}
	return nil
}
