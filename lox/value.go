package lox

import (
	"fmt"
)

type Value interface {
	isValue()
	String() string
	IsTruthy() bool
}

type String string
type Number float64
type Nil struct{}
type Bool bool
type CallableVal struct {
	Callable
}

func (s String) isValue()      {}
func (n Number) isValue()      {}
func (b Bool) isValue()        {}
func (n Nil) isValue()         {}
func (n CallableVal) isValue() {}

func (s String) String() string { return string(s) }

// FIXME changing the format specifier %g might recurse forever, create more robust impl
func (n Number) String() string { return fmt.Sprintf("%g", n) }

func (b Bool) String() string {
	if bool(b) == true {
		return "true"
	} else {
		return "false"
	}
}

func (n Nil) String() string {
	return "nil"
}

func (s String) IsTruthy() bool    { return true }
func (n Number) IsTruthy() bool    { return n != 0 }
func (b Bool) IsTruthy() bool      { return bool(b) }
func (Nil) IsTruthy() bool         { return false }
func (CallableVal) IsTruthy() bool { return true }
