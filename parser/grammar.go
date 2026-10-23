// recursive descent parser for the Lox grammar https://craftinginterpreters.com/appendix-i.html
package parser

import (
	"fmt"

	"github.com/goncalorazevedo/glox/ast"
	"github.com/goncalorazevedo/glox/lexer"
)

func (p *parser) parse() ([]ast.Stmnt, error) {
	stmnts := []ast.Stmnt{}
	for !p.atEnd() {
		stmnt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmnts = append(stmnts, stmnt)
	}
	return stmnts, nil
}

func (p *parser) declaration() (ast.Stmnt, error) {
	switch t := p.currToken(); t.Type {
	case lexer.Var:
		p.index++
		return p.varDeclaration()
	case lexer.Fun:
		p.index++
		return p.funDeclaration()
	case lexer.Class:
	}
	return p.statement()
}

func (p *parser) varDeclaration() (ast.Stmnt, error) {
	if p.atEnd() || !p.match(lexer.Identifier) {
		return nil, fmt.Errorf(`line %d: missing identifier on variable declaration`, p.prevToken().Line)
	}
	identifier := p.currToken()
	p.index++

	var init ast.Expr
	if !p.atEnd() && p.match(lexer.Equal) {
		p.index++
		var err error
		init, err = p.expression()
		if err != nil {
			return nil, fmt.Errorf(`line %d: failed parsing initializer on "%s" declaration: %w`, identifier.Line, identifier.Lexeme, err)
		}
	}
	if p.atEnd() || !p.match(lexer.Semicolon) {
		return nil, fmt.Errorf(`line %d: missing ";" on variable declaration`, p.prevToken().Line)
	}
	p.index++

	return &ast.VarDecl{Identifier: identifier, Init: init}, nil
}

func (p *parser) funDeclaration() (ast.Stmnt, error) {
	if p.atEnd() || !p.match(lexer.Identifier) {
		return nil, fmt.Errorf(`line %d: missing identifier on function declaration`, p.prevToken().Line)
	}
	identifier := p.currToken()

	p.index++
	if p.atEnd() || !p.match(lexer.LeftParen) {
		return nil, fmt.Errorf(`line %d: missing "(" on function declaration`, identifier.Line)
	}
	p.index++

	if p.atEnd() {
		return nil, fmt.Errorf(`line %d: unexpected EOF parsing function "%s" parameters`, identifier.Line, identifier.Lexeme)
	}
	params := []lexer.Token{}
	if p.match(lexer.RightParen) {
		p.index++ // no parameters
	} else {
		for {
			if !p.match(lexer.Identifier) {
				return nil, fmt.Errorf(`line %d: unexpected token while parsing "%s" parameters`, identifier.Line, identifier.Lexeme)
			}
			params = append(params, p.currToken())
			p.index++
			if p.atEnd() {
				return nil, fmt.Errorf(`line %d: unexpected EOF parsing function "%s" parameters`, identifier.Line, identifier.Lexeme)
			}
			if p.match(lexer.Comma) && p.isPeekable() {
				p.index++
				continue
			} else if p.match(lexer.RightParen) {
				p.index++
				break
			} else {
				return nil, fmt.Errorf(`line %d: unexpected EOF parsing function "%s" parameters`, identifier.Line, identifier.Lexeme)
			}
		}
	}
	body, err := p.block()
	if err != nil {
		return nil, fmt.Errorf(`line %d: unable to parse function "%s" body: %w`, identifier.Line, identifier.Lexeme, err)
	}
	return &ast.FunDecl{Identifier: identifier, Params: params, Body: body}, nil
}

func (p *parser) statement() (ast.Stmnt, error) {
	switch p.currToken().Type {
	case lexer.Print:
		p.index++
		return p.printSmnt()
	case lexer.LeftBrace:
		return p.block()
	case lexer.If:
		p.index++
		return p.ifStmnt()
	case lexer.While:
		p.index++
		return p.whileStmnt()
	case lexer.For:
		p.index++
		return p.forStmnt()
	case lexer.Return:
		p.index++
		return p.returnStmnt()
	default:
		return p.expressionStmnt()
	}
}

func (p *parser) printSmnt() (ast.Stmnt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, fmt.Errorf(`line %d: missing expression on print statement`, p.prevToken().Line)
	}
	if p.atEnd() || !p.match(lexer.Semicolon) {
		return nil, fmt.Errorf(`line %d: missing ";" on print statement`, p.prevToken().Line)
	}
	p.index++
	return &ast.PrintSmnt{Expr: expr}, nil
}

func (p *parser) block() (ast.Stmnt, error) {
	if p.atEnd() || !p.match(lexer.LeftBrace) {
		return nil, fmt.Errorf("unexpected token while parsing block statement")
	}
	p.index++
	stmnts := []ast.Stmnt{}
	for {
		if p.atEnd() {
			return nil, fmt.Errorf("unexpected EOF while parsing block statement")
		}
		if p.match(lexer.RightBrace) {
			p.index++
			break
		}
		stmnt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmnts = append(stmnts, stmnt)
	}
	return &ast.Block{Stmnts: stmnts}, nil
}

func (p *parser) ifStmnt() (ast.Stmnt, error) {
	ifToken := p.prevToken()
	if p.atEnd() || !p.match(lexer.LeftParen) {
		return nil, fmt.Errorf(`line %d: missing '(' in if statement`, ifToken.Line)
	}
	p.index++
	cond, err := p.expression()
	if err != nil {
		return nil, fmt.Errorf(`line %d: invalid condition in if statement: %w`, ifToken.Line, err)
	}
	if p.atEnd() || !p.match(lexer.RightParen) {
		return nil, fmt.Errorf(`line %d: missing ')' in if statement`, ifToken.Line)
	}
	p.index++
	then, err := p.statement()
	if err != nil {
		return nil, fmt.Errorf(`line %d: invalid statement in if statement: %w`, ifToken.Line, err)
	}
	var elseThen ast.Stmnt
	if !p.atEnd() && p.match(lexer.Else) {
		p.index++
		var err error
		elseThen, err = p.statement()
		if err != nil {
			return nil, fmt.Errorf(`line %d: invalid else branch in if statement: %w`, ifToken.Line, err)
		}
	}
	return &ast.IfStmnt{Condition: cond, Then: then, ElseThen: elseThen}, nil
}

func (p *parser) whileStmnt() (ast.Stmnt, error) {
	whileToken := p.prevToken()
	if p.atEnd() || !p.match(lexer.LeftParen) {
		return nil, fmt.Errorf(`line %d: missing '(' in while statement`, whileToken.Line)
	}
	p.index++
	condition, err := p.expression()
	if err != nil {
		return nil, fmt.Errorf(`line %d: invalid condition in while statement: %w`, whileToken.Line, err)
	}
	if p.atEnd() || !p.match(lexer.RightParen) {
		return nil, fmt.Errorf(`line %d: missing ')' in while statement`, whileToken.Line)
	}
	p.index++
	stmnt, err := p.statement()
	if err != nil {
		return nil, fmt.Errorf(`line %d: invalid statement in while statement: %w`, whileToken.Line, err)
	}
	return &ast.WhileStmnt{Condition: condition, Stmnt: stmnt}, nil
}

func (p *parser) forStmnt() (ast.Stmnt, error) {
	forToken := p.prevToken()
	if p.atEnd() || !p.match(lexer.LeftParen) {
		return nil, fmt.Errorf(`line %d: missing '(' in for statement`, forToken.Line)
	}

	p.index++
	var initializer ast.Stmnt
	if p.atEnd() {
		return nil, fmt.Errorf(`line %d: unexpected EOF while parsing for tatement`, forToken.Line)
	}
	if p.match(lexer.Semicolon) {
		p.index++
	} else {
		var err error
		initializer, err = p.declaration()
		if err != nil {
			return nil, fmt.Errorf(`line %d: unable to parse initializer in for statement: %w`, forToken.Line, err)
		}
		_, isVarDecl := initializer.(*ast.VarDecl)
		_, isExprStmnt := initializer.(*ast.ExprStmnt)
		if !isVarDecl && !isExprStmnt {
			return nil, fmt.Errorf(`line %d: invalid initializer in for statement initializer`, forToken.Line)
		}
	}

	var condition ast.Expr
	if p.atEnd() {
		return nil, fmt.Errorf(`line %d: unexpected EOF while parsing for statement condition`, forToken.Line)
	}
	if !p.match(lexer.Semicolon) {
		var err error
		condition, err = p.expression()
		if err != nil {
			return nil, fmt.Errorf(`line %d: unable to parse condition in for statement: %w`, forToken.Line, err)
		}
		if p.atEnd() || !p.match(lexer.Semicolon) {
			return nil, fmt.Errorf(`line %d: missing ';' in for statement condition`, forToken.Line)
		}
	}
	p.index++

	var increment ast.Expr
	if p.atEnd() {
		return nil, fmt.Errorf(`line %d: unexpected EOF while parsing for statement increment`, forToken.Line)
	}
	if !p.match(lexer.RightParen) {
		var err error
		increment, err = p.expression()
		if err != nil {
			return nil, fmt.Errorf(`line %d: unable to parse increment in for statement: %w`, forToken.Line, err)
		}
		if p.atEnd() || !p.match(lexer.RightParen) {
			return nil, fmt.Errorf(`line %d: missing ')' in for statement`, forToken.Line)
		}
	}
	p.index++

	if p.atEnd() {
		return nil, fmt.Errorf(`line %d: unexpected EOF while parsing for statement increment`, forToken.Line)
	}

	stmnt, err := p.statement()
	if err != nil {
		return nil, fmt.Errorf(`line %d: unable to parse body in for statement: %w`, forToken.Line, err)
	}

	var body ast.Stmnt
	if increment == nil {
		body = stmnt
	} else {
		body = &ast.Block{Stmnts: []ast.Stmnt{stmnt, &ast.ExprStmnt{Expr: increment}}}
	}
	if condition == nil {
		condition = &ast.Primary{Value: ast.Bool(true)}
	}
	loop := &ast.WhileStmnt{Condition: condition, Stmnt: body}
	if initializer == nil {
		return loop, nil
	} else {
		return &ast.Block{Stmnts: []ast.Stmnt{initializer, loop}}, nil
	}
}

func (p *parser) returnStmnt() (ast.Stmnt, error) {
	var expr ast.Expr
	if !p.atEnd() && !p.match(lexer.Semicolon) {
		var err error
		expr, err = p.expression()
		if err != nil {
			return nil, fmt.Errorf("line %d: missing expression on return", p.prevToken().Line)
		}
	}
	if p.atEnd() || !p.match(lexer.Semicolon) {
		return nil, fmt.Errorf(`line %d: missing ";" on return statement`, p.prevToken().Line)
	}
	p.index++
	return &ast.Return{Expr: expr}, nil
}

func (p *parser) expressionStmnt() (ast.Stmnt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if p.atEnd() || !p.match(lexer.Semicolon) {
		return nil, fmt.Errorf(`line %d: missing ";" on expression statement`, p.prevToken().Line)
	}
	p.index++
	return &ast.ExprStmnt{Expr: expr}, nil
}

func (p *parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *parser) assignment() (ast.Expr, error) {
	or, err := p.or()
	if err != nil {
		return nil, err
	}
	if !p.atEnd() && p.match(lexer.Equal) {
		lval, ok := or.(*ast.Variable)
		if !ok {
			return nil, fmt.Errorf("line %d: invalid lval for assginment", p.currToken().Line)
		}
		p.index++
		expr, err := p.assignment()
		if err != nil {
			return nil, err
		}
		return &ast.Assignment{Identifier: lval.Identifier, Expr: expr}, nil
	}
	return or, nil
}

func (p *parser) or() (ast.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for !p.atEnd() && p.match(lexer.Or) {
		token := p.currToken()
		p.index++
		if p.atEnd() {
			return nil, fmt.Errorf(`unexpected EOF while parsing "%s" expression`, token.Lexeme)
		}
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &ast.Logical{Left: expr, Token: token, Right: right}
	}
	return expr, nil
}

func (p *parser) and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	for !p.atEnd() && p.match(lexer.And) {
		token := p.currToken()
		p.index++
		if p.atEnd() {
			return nil, fmt.Errorf(`unexpected EOF while parsing "%s" expression`, token.Lexeme)
		}
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &ast.Logical{Left: expr, Token: token, Right: right}
	}
	return expr, nil
}

func (p *parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for !p.atEnd() && p.matchAny(lexer.EqualEqual, lexer.BangEqual) {
		token := p.currToken()
		p.index++
		if p.atEnd() {
			return nil, fmt.Errorf(`unexpected EOF while parsing "%s" expression`, token.Lexeme)
		}
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Token: token, Right: right}
	}
	return expr, nil
}

func (p *parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for !p.atEnd() && p.matchAny(lexer.GreaterEqual, lexer.Greater, lexer.LessEqual, lexer.Less) {
		token := p.currToken()
		p.index++
		if p.atEnd() {
			return nil, fmt.Errorf(`unexpected EOF while parsing "%s" expression`, token.Lexeme)
		}
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Token: token, Right: right}
	}
	return expr, nil
}

func (p *parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for !p.atEnd() && p.matchAny(lexer.Plus, lexer.Minus) {
		token := p.tokens[p.index]
		p.index++
		if p.atEnd() {
			return nil, fmt.Errorf(`unexpected EOF while parsing "%s" expression`, token.Lexeme)
		}
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Token: token, Right: right}
	}
	return expr, nil
}

func (p *parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for !p.atEnd() && p.matchAny(lexer.Star, lexer.Slash) {
		token := p.tokens[p.index]
		p.index++
		if p.atEnd() {
			return nil, fmt.Errorf(`unexpected EOF while parsing "%s" expression`, token.Lexeme)
		}
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Token: token, Right: right}
	}
	return expr, nil
}

func (p *parser) unary() (ast.Expr, error) {
	if !p.atEnd() && p.matchAny(lexer.Bang, lexer.Minus) {
		token := p.tokens[p.index]
		p.index++
		if p.atEnd() {
			return nil, fmt.Errorf(`unexpected EOF while parsing "%s" expression`, token.Lexeme)
		}
		expr, err := p.unary()
		if err != nil {
			return nil, fmt.Errorf(`failed parsing "%s" expression: %w`, token.Lexeme, err)
		}
		return &ast.Unary{Token: token, Expr: expr}, nil
	}
	return p.call()
}

func (p *parser) call() (ast.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	if !p.atEnd() && p.match(lexer.LeftParen) {
		p.index++
		if p.atEnd() {
			return nil, fmt.Errorf("unexpected EOF parsing call arguments")
		}
		args := []ast.Expr{}
		if p.match(lexer.RightParen) {
			p.index++
		} else {
			for {
				arg, err := p.expression()
				if err != nil {
					return nil, fmt.Errorf("unable to parse function call arguments")
				}
				args = append(args, arg)
				if p.atEnd() {
					return nil, fmt.Errorf(`unexpected EOF parsing function call arguments`)
				}
				if p.match(lexer.Comma) && p.isPeekable() {
					p.index++
					continue
				} else if p.match(lexer.RightParen) {
					p.index++
					break
				} else {
					return nil, fmt.Errorf("unexpected token while pargin call arguments")
				}
			}
		}
		return &ast.Call{Callee: expr, Arguments: args}, nil
	}
	return expr, nil
}

func (p *parser) primary() (ast.Expr, error) {
	if p.atEnd() {
		return nil, fmt.Errorf("unexpected EOF parsing expression")
	}

	switch t := p.currToken(); t.Type {
	case lexer.True, lexer.False:
		p.index++
		return &ast.Primary{Value: ast.Bool(t.Value.(bool))}, nil
	case lexer.String:
		p.index++
		return &ast.Primary{Value: ast.String(t.Value.(string))}, nil
	case lexer.Number:
		p.index++
		return &ast.Primary{Value: ast.Number(t.Value.(float64))}, nil
	case lexer.Nil:
		p.index++
		return &ast.Primary{Value: ast.Nil(t.Value.(struct{}))}, nil
	case lexer.Identifier:
		p.index++
		return &ast.Variable{Identifier: p.prevToken()}, nil
	case lexer.LeftParen:
		p.index++
		if p.atEnd() {
			return nil, fmt.Errorf("unexpected EOF while parsing primary expression")
		}
		expr, err := p.expression()
		if err != nil {
			return nil, fmt.Errorf("failed parsing grouping expression: %w", err)
		}
		if p.atEnd() || !p.match(lexer.RightParen) {
			return nil, fmt.Errorf(`missing closing ")" on grouping expression`)
		}
		p.index++
		return &ast.Grouping{Expr: expr}, nil
	default:
		return nil, fmt.Errorf("unexpected token while parsing expression")
	}
}

func (p *parser) currToken() lexer.Token {
	return p.tokens[p.index]
}

func (p *parser) prevToken() lexer.Token {
	return p.tokens[p.index-1]
}

func (p *parser) match(t lexer.TokenType) bool {
	return p.tokens[p.index].Type == t
}

func (p *parser) matchAny(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.tokens[p.index].Type == t {
			return true
		}
	}
	return false
}

// returns false if next token is EOF
func (p *parser) isPeekable() bool {
	if p.index+1 >= len(p.tokens) {
		return false
	}
	return p.tokens[p.index+1].Type != lexer.EOF
}

func (p *parser) atEnd() bool {
	if p.index >= len(p.tokens) {
		return true
	}
	return p.tokens[p.index].Type == lexer.EOF
}
