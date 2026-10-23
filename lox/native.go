package lox

import (
	"fmt"
	"time"

	"github.com/goncalorazevedo/glox/ast"
)

type Native struct {
	Name  string
	arity int
	Fn    func(...ast.Value) (ast.Value, error)
}

var _ ast.Callable = (*Native)(nil)

func (n *Native) Arity() int {
	return n.arity
}

func (n *Native) String() string {
	return fmt.Sprintf("<native  %s(%d)>", n.Name, n.Arity())
}

func (n *Native) Call(ctx any, args ...ast.Value) (ast.Value, error) {
	if n.Arity() != len(args) {
		panic(fmt.Sprintf("invalid number of arguments provided to %s", n.String()))
	}
	return n.Fn(args...)
}

func getClockNative() *Native {
	return &Native{
		Name:  "clock",
		arity: 0,
		Fn: func(_ ...ast.Value) (ast.Value, error) {
			return ast.Number(float64(time.Now().UnixNano()) / 1e6), nil
		},
	}
}

func getSleepNative() *Native {
	return &Native{
		Name:  "sleep",
		arity: 1,
		Fn: func(args ...ast.Value) (ast.Value, error) {
			duration, ok := args[0].(ast.Number)
			if !ok {
				return nil, fmt.Errorf("invalid argument to sleep, expected Number")
			}
			time.Sleep(time.Duration(duration) * time.Millisecond)
			return ast.Nil{}, nil
		},
	}
}
