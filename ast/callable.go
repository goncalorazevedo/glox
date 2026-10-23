package ast

type Callable interface {
	String() string
	Arity() int
	Call(any, ...Value) (Value, error)
}
