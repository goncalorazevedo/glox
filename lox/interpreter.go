package lox

import (
	"fmt"

	"github.com/goncalorazevedo/glox/ast"
	"github.com/goncalorazevedo/glox/lexer"
)

type interpreter struct {
	env *environment
}

var _ ast.ExprVisitor[Value, Value] = (*interpreter)(nil)
var _ ast.StmntVisitor = (*interpreter)(nil)

func Evaluate(stmnts []ast.Stmnt) error {
	interp := interpreter{
		env: newEnvironment(nil),
	}
	if err := interp.env.define("clock", CallableVal{Callable: getClockNative()}); err != nil {
		panic(err)
	}
	if err := interp.env.define("sleep", CallableVal{Callable: getSleepNative()}); err != nil {
		panic(err)
	}
	for _, stmnt := range stmnts {
		if err := interp.evaluateStmnt(stmnt); err != nil {
			panic(err)
		}
	}
	return nil
}

func (i *interpreter) evaluateStmnt(s ast.Stmnt) error {
	return ast.AcceptStmnt(s, i)
}

func (i *interpreter) evaluateExpr(expr ast.Expr) (Value, error) {
	return ast.AcceptExpr(expr, i)
}

func (i *interpreter) VisitPrintStmnt(p *ast.PrintSmnt) error {
	value, err := i.evaluateExpr(p.Expr)
	if err == nil {
		fmt.Println(value)
	}
	return err
}

func (i *interpreter) VisitExprStmnt(e *ast.ExprStmnt) error {
	_, err := i.evaluateExpr(e.Expr)
	return err
}

func (i *interpreter) VisitVarDecl(v *ast.VarDecl) error {
	var value Value = Nil{}
	if v.Init != nil {
		var err error
		value, err = i.evaluateExpr(v.Init)
		if err != nil {
			return err
		}
	}
	if err := i.env.define(v.Identifier.Lexeme, value); err != nil {
		return err
	}
	return nil
}

func (i *interpreter) VisitFunDecl(f *ast.FunDecl) error {
	params := []string{}
	for _, p := range f.Params {
		params = append(params, p.Lexeme)
	}
	fun := &Function{
		Name:    f.Identifier.Lexeme,
		Params:  params,
		Stmnt:   f.Body,
		Closure: i.env,
	}
	if err := i.env.define(f.Identifier.Lexeme, CallableVal{Callable: fun}); err != nil {
		return err
	}
	return nil
}

func (i *interpreter) VisitBlock(b *ast.Block) error {
	i.env = newEnvironment(i.env)
	defer func() { i.env = i.env.enclosing }()
	for _, stmnt := range b.Stmnts {
		if err := i.evaluateStmnt(stmnt); err != nil {
			return err
		}
	}
	return nil
}

func (i *interpreter) VisitIfStmnt(s *ast.IfStmnt) error {
	value, err := i.evaluateExpr(s.Condition)
	if err != nil {
		return err
	}
	if value.IsTruthy() {
		return i.evaluateStmnt(s.Then)
	} else if s.ElseThen != nil {
		return i.evaluateStmnt(s.ElseThen)
	} else {
		return nil
	}
}

func (i *interpreter) VisitWhileStmnt(w *ast.WhileStmnt) error {
	for {
		value, err := i.evaluateExpr(w.Condition)
		if err != nil {
			return err
		}
		if !value.IsTruthy() {
			return nil
		}
		if err := i.evaluateStmnt(w.Stmnt); err != nil {
			return err
		}
	}
}

func (i *interpreter) VisitPrimary(p *ast.Primary[Value]) (Value, error) {
	return p.Value, nil
}

func (i *interpreter) VisitGrouping(g *ast.Grouping) (Value, error) {
	return i.evaluateExpr(g.Expr)
}

func (i *interpreter) VisitUnary(u *ast.Unary) (Value, error) {
	value, err := i.evaluateExpr(u.Expr)
	if err != nil {
		return nil, err
	}
	switch t := u.Token; t.Type {
	case lexer.Bang:
		switch v := value.(type) {
		case Bool:
			return Bool(!v), nil
		default:
			return nil, fmt.Errorf("line %d: operand of '!' must be a boolean ", u.Token.Line)
		}
	case lexer.Minus:
		switch v := value.(type) {
		case Number:
			return Number(-v), nil
		default:
			return nil, fmt.Errorf("line %d: operand of '-' must be a number", u.Token.Line)
		}
	default:
		return nil, fmt.Errorf(`line %d: invalid operator "%s" for unary expression`, t.Line, t.Lexeme)
	}
}

func (i *interpreter) VisitBinary(b *ast.Binary) (Value, error) {
	left, err := i.evaluateExpr(b.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluateExpr(b.Right)
	if err != nil {
		return nil, err
	}

	switch t := b.Token; t.Type {
	case lexer.EqualEqual:
		return Bool(left == right), nil // Go interface values make this work for lox
	case lexer.BangEqual:
		return Bool(left != right), nil // Go interface values make this work for lox
	case lexer.Less:
		less, err := isLess(left, right)
		if err != nil {
			return nil, err
		}
		return Bool(less), nil
	case lexer.LessEqual:
		less, err := isLess(left, right)
		if err != nil {
			return nil, err
		}
		return Bool(less || left == right), nil
	case lexer.Greater:
		less, err := isLess(left, right)
		if err != nil {
			return nil, err
		}
		return Bool(!less && left != right), nil
	case lexer.GreaterEqual:
		less, err := isLess(left, right)
		if err != nil {
			return nil, err
		}
		return Bool(!less), nil
	case lexer.Plus, lexer.Minus, lexer.Star, lexer.Slash:
		result, err := lookupBinaryOperation(t.Type)(left, right)
		if err != nil {
			return nil, fmt.Errorf(`line %d: invalid operation for "%s": %w`, t.Line, left, err)
		}
		return result, nil
	default:
		panic(fmt.Sprintf(`line %d: invalid operator "%s" for binary expression`, t.Line, t.Lexeme))
	}
}

func (i *interpreter) VisitLogical(l *ast.Logical) (Value, error) {
	left, err := i.evaluateExpr(l.Left)
	if err != nil {
		return nil, err
	}
	// short circuit
	switch t := l.Token; t.Type {
	case lexer.And:
		if !left.IsTruthy() {
			return left, nil
		}
	case lexer.Or:
		if left.IsTruthy() {
			return left, nil
		}
	default:
		panic(fmt.Sprintf(`line %d: invalid operator "%s" for logical expression`, t.Line, t.Lexeme))
	}
	return i.evaluateExpr(l.Right)
}

func (i *interpreter) VisitVariable(v *ast.Variable) (Value, error) {
	value := i.env.lookup(v.Identifier.Lexeme)
	if value == nil {
		return nil, fmt.Errorf(`line %d: variable "%s" not defined in scope`, v.Identifier.Line, v.Identifier.Lexeme)
	}
	return value, nil
}

func (i *interpreter) VisitAssignment(a *ast.Assignment) (Value, error) {
	value, err := i.evaluateExpr(a.Expr)
	if err != nil {
		return nil, err
	}
	if err := i.env.assign(a.Identifier.Lexeme, value); err != nil {
		return nil, fmt.Errorf(`line %d: variable "%s" not defined in scope`, a.Identifier.Line, a.Identifier.Lexeme)
	}
	return value, nil
}

// following the flow of the book, which uses exceptions to surface the return value,
// even though this is much more rudimentary it works similarly as we manually check
// for errors everywhere
type returnVal struct {
	value Value
}

func (r *returnVal) Error() string {
	return fmt.Sprintf("%s", r.value)
}

func (i *interpreter) VisitCall(c *ast.Call) (Value, error) {
	callee, err := i.evaluateExpr(c.Callee)
	if err != nil {
		return nil, err
	}
	callable, ok := callee.(CallableVal)
	if !ok {
		return nil, fmt.Errorf("invalid call expression, not a callable value")
	}
	if callable.Arity() != len(c.Arguments) {
		return nil, fmt.Errorf("invalid number of arguments to %s, got %d", callable, len(c.Arguments))
	}
	var args []Value
	for _, a := range c.Arguments {
		v, err := i.evaluateExpr(a)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}
	ret, err := callable.Call(i, args...)
	if err != nil {
		if ret, ok := err.(*returnVal); ok {
			return ret.value, nil
		}
		return nil, err
	}
	return ret, nil
}

func (i *interpreter) VisitReturn(r *ast.Return) error {
	var value Value = Nil{}
	if r.Expr != nil {
		var err error
		value, err = i.evaluateExpr(r.Expr)
		if err != nil {
			return err
		}
	}
	return &returnVal{value: value}
}
