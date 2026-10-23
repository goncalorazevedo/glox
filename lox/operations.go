package lox

import (
	"errors"
	"fmt"

	"github.com/goncalorazevedo/glox/ast"
	"github.com/goncalorazevedo/glox/lexer"
)

func isLess(a, b ast.Value) (bool, error) {
	switch v := a.(type) {
	case ast.Number:
		v2, ok := b.(ast.Number)
		if !ok {
			return false, fmt.Errorf("")
		}
		return v < v2, nil
	case ast.String:
		v2, ok := b.(ast.String)
		if !ok {
			return false, fmt.Errorf("")
		}
		return len(v) < len(v2), nil
	default:
		return false, errors.New("operands must be two numbers or two strings")
	}
}

func plus(a, b ast.Value) (ast.Value, error) {
	switch v1 := a.(type) {
	case ast.Number:
		v2, ok := b.(ast.Number)
		if !ok {
			return nil, fmt.Errorf("incompatible types %T and %T", a, b)
		}
		return ast.Number(v1 + v2), nil
	case ast.String:
		v2, ok := b.(ast.String)
		if !ok {
			return nil, fmt.Errorf("incompatible types %T and %T", a, b)
		}
		return ast.String(v1 + v2), nil
	default:
		return nil, fmt.Errorf("incompatible value %s of type %T", a.String(), a)
	}
}

func minus(a, b ast.Value) (ast.Value, error) {
	switch v1 := a.(type) {
	case ast.Number:
		v2, ok := b.(ast.Number)
		if !ok {
			return nil, fmt.Errorf("incompatible types %T and %T", a, b)
		}
		return ast.Number(v1 - v2), nil
	default:
		return nil, fmt.Errorf("incompatible value %s of type %T", a.String(), a)
	}
}

func star(a, b ast.Value) (ast.Value, error) {
	switch v1 := a.(type) {
	case ast.Number:
		v2, ok := b.(ast.Number)
		if !ok {
			return nil, fmt.Errorf("incompatible types %T and %T", a, b)
		}
		return ast.Number(v1 * v2), nil
	default:
		return nil, fmt.Errorf("incompatible value %s of type %T", a.String(), a)
	}
}

func slash(a, b ast.Value) (ast.Value, error) {
	switch v1 := a.(type) {
	case ast.Number:
		v2, ok := b.(ast.Number)
		if !ok {
			return nil, fmt.Errorf("incompatible types %T and %T", a, b)
		}
		return ast.Number(v1 / v2), nil
	default:
		return nil, fmt.Errorf("incompatible value %s of type %T", a.String(), a)
	}
}

var binaryFunctionsMap = map[lexer.TokenType]func(a, b ast.Value) (ast.Value, error){
	lexer.Plus:  plus,
	lexer.Minus: minus,
	lexer.Star:  star,
	lexer.Slash: slash,
}

func lookupBinaryOperation(t lexer.TokenType) func(a, b ast.Value) (ast.Value, error) {
	operation, exists := binaryFunctionsMap[t]
	if !exists {
		panic("")
	}
	return operation
}
