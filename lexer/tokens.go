package lexer

import "fmt"

type TokenType int

type Token struct {
	Type   TokenType
	Line   int
	Lexeme string
	Value  any
}

const (
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Star
	Slash
	Semicolon

	Bang
	BangEqual
	Equal
	EqualEqual
	Less
	LessEqual
	Greater
	GreaterEqual

	True
	False
	Identifier
	String
	Number
	Nil
	And
	Or
	If
	Else
	For
	While
	Var
	Fun
	Class
	This
	Super
	Return
	Print

	EOF
)

func (t TokenType) String() string {
	if int(t) >= 0 && int(t) < len(tokenNames) {
		return tokenNames[t]
	}
	return fmt.Sprintf("TokenType(%d)", int(t))
}

var tokenNames = [...]string{
	LeftParen: "LeftParen", RightParen: "RightParen", LeftBrace: "LeftBrace", RightBrace: "RightBrace",
	Comma: "Comma", Dot: "Dot", Minus: "Minus", Plus: "Plus", Star: "Star", Slash: "Slash", Semicolon: "Semicolon",

	Bang: "Bang", BangEqual: "BangEqual",
	Equal: "Equal", EqualEqual: "EqualEqual",
	Less: "Less", LessEqual: "LessEqual",
	Greater: "Greater", GreaterEqual: "GreaterEqual",

	True: "True", False: "False", Identifier: "Identifier", String: "String", Number: "Number", Nil: "Nil",
	And: "And", Or: "Or", If: "If", Else: "Else", For: "For", While: "While", Var: "Var", Fun: "Fun",
	Class: "Class", This: "This", Super: "Super", Return: "Return", Print: "Print",

	EOF: "EOF",
}
