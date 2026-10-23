package lexer

import (
	"fmt"
	"os"
	"testing"
)

type validTestCase struct {
	name     string
	expected []Token
}

// we ignore token values here because it might change in the future
var validTestCases = []validTestCase{
	{name: "tc1.lox", expected: []Token{
		{Type: Number, Line: 1, Lexeme: "1", Value: nil},
		{Type: And, Line: 1, Lexeme: "and", Value: nil},
		{Type: Number, Line: 1, Lexeme: "2", Value: nil},
		{Type: Or, Line: 1, Lexeme: "or", Value: nil},
		{Type: True, Line: 1, Lexeme: "true", Value: nil},
		{Type: And, Line: 1, Lexeme: "and", Value: nil},
		{Type: String, Line: 1, Lexeme: "false", Value: nil},
		{Type: Semicolon, Line: 1, Lexeme: ";", Value: nil},
		{Type: EOF, Line: 1, Lexeme: "", Value: nil},
	}},
	{name: "tc2.lox", expected: []Token{
		{Type: Fun, Line: 1, Lexeme: "fun", Value: nil},
		{Type: Identifier, Line: 1, Lexeme: "printer", Value: nil},
		{Type: LeftParen, Line: 1, Lexeme: "(", Value: nil},
		{Type: Identifier, Line: 1, Lexeme: "a", Value: nil},
		{Type: RightParen, Line: 1, Lexeme: ")", Value: nil},
		{Type: LeftBrace, Line: 1, Lexeme: "{", Value: nil},
		{Type: Print, Line: 2, Lexeme: "print", Value: nil},
		{Type: Identifier, Line: 2, Lexeme: "a", Value: nil},
		{Type: Semicolon, Line: 2, Lexeme: ";", Value: nil},
		{Type: RightBrace, Line: 3, Lexeme: "}", Value: nil},
		{Type: Identifier, Line: 5, Lexeme: "printer", Value: nil},
		{Type: LeftParen, Line: 5, Lexeme: "(", Value: nil},
		{Type: String, Line: 5, Lexeme: "hello world", Value: nil},
		{Type: RightParen, Line: 5, Lexeme: ")", Value: nil},
		{Type: Semicolon, Line: 5, Lexeme: ";", Value: nil},
		{Type: EOF, Line: 5, Lexeme: "", Value: nil},
	}},
	{name: "tc3.lox", expected: []Token{
		{Type: Var, Line: 1, Lexeme: "var", Value: nil},
		{Type: Identifier, Line: 1, Lexeme: "a", Value: nil},
		{Type: Equal, Line: 1, Lexeme: "=", Value: nil},
		{Type: Number, Line: 1, Lexeme: "2.0", Value: nil},
		{Type: Semicolon, Line: 1, Lexeme: ";", Value: nil},
		{Type: If, Line: 2, Lexeme: "if", Value: nil},
		{Type: LeftParen, Line: 2, Lexeme: "(", Value: nil},
		{Type: Identifier, Line: 2, Lexeme: "a", Value: nil},
		{Type: EqualEqual, Line: 2, Lexeme: "==", Value: nil},
		{Type: Number, Line: 2, Lexeme: "2", Value: nil},
		{Type: RightParen, Line: 2, Lexeme: ")", Value: nil},
		{Type: LeftBrace, Line: 2, Lexeme: "{", Value: nil},
		{Type: Print, Line: 3, Lexeme: "print", Value: nil},
		{Type: Identifier, Line: 3, Lexeme: "a", Value: nil},
		{Type: Semicolon, Line: 3, Lexeme: ";", Value: nil},
		{Type: RightBrace, Line: 4, Lexeme: "}", Value: nil},
		{Type: Else, Line: 4, Lexeme: "else", Value: nil},
		{Type: LeftBrace, Line: 4, Lexeme: "{", Value: nil},
		{Type: RightBrace, Line: 5, Lexeme: "}", Value: nil},
		{Type: EOF, Line: 5, Lexeme: "", Value: nil},
	}},
	{name: "tc4.lox", expected: []Token{
		{Type: EOF, Line: 1, Lexeme: "", Value: nil},
	}},
	{name: "tc5.lox", expected: []Token{
		{Type: Nil, Line: 1, Lexeme: "nil", Value: nil},
		{Type: LessEqual, Line: 1, Lexeme: "<=", Value: nil},
		{Type: False, Line: 1, Lexeme: "false", Value: nil},
		{Type: Semicolon, Line: 1, Lexeme: ";", Value: nil},
		{Type: Bang, Line: 2, Lexeme: "!", Value: nil},
		{Type: True, Line: 2, Lexeme: "true", Value: nil},
		{Type: BangEqual, Line: 2, Lexeme: "!=", Value: nil},
		{Type: Bang, Line: 2, Lexeme: "!", Value: nil},
		{Type: True, Line: 2, Lexeme: "true", Value: nil},
		{Type: Semicolon, Line: 2, Lexeme: ";", Value: nil},
		{Type: EOF, Line: 2, Lexeme: "", Value: nil},
	}},
	{name: "tc6.lox", expected: []Token{
		{Type: For, Line: 1, Lexeme: "for", Value: nil},
		{Type: LeftParen, Line: 1, Lexeme: "(", Value: nil},
		{Type: Var, Line: 1, Lexeme: "var", Value: nil},
		{Type: Identifier, Line: 1, Lexeme: "a", Value: nil},
		{Type: Equal, Line: 1, Lexeme: "=", Value: nil},
		{Type: Number, Line: 1, Lexeme: "0", Value: nil},
		{Type: Semicolon, Line: 1, Lexeme: ";", Value: nil},
		{Type: Identifier, Line: 1, Lexeme: "i", Value: nil},
		{Type: Less, Line: 1, Lexeme: "<", Value: nil},
		{Type: Number, Line: 1, Lexeme: "10", Value: nil},
		{Type: Semicolon, Line: 1, Lexeme: ";", Value: nil},
		{Type: Identifier, Line: 1, Lexeme: "i", Value: nil},
		{Type: Equal, Line: 1, Lexeme: "=", Value: nil},
		{Type: Identifier, Line: 1, Lexeme: "i", Value: nil},
		{Type: Plus, Line: 1, Lexeme: "+", Value: nil},
		{Type: Number, Line: 1, Lexeme: "1", Value: nil},
		{Type: RightParen, Line: 1, Lexeme: ")", Value: nil},
		{Type: LeftBrace, Line: 1, Lexeme: "{", Value: nil},
		{Type: Print, Line: 2, Lexeme: "print", Value: nil},
		{Type: Identifier, Line: 2, Lexeme: "i", Value: nil},
		{Type: Semicolon, Line: 2, Lexeme: ";", Value: nil},
		{Type: RightBrace, Line: 3, Lexeme: "}", Value: nil},
		{Type: While, Line: 5, Lexeme: "while", Value: nil},
		{Type: LeftParen, Line: 5, Lexeme: "(", Value: nil},
		{Type: True, Line: 5, Lexeme: "true", Value: nil},
		{Type: RightParen, Line: 5, Lexeme: ")", Value: nil},
		{Type: LeftBrace, Line: 5, Lexeme: "{", Value: nil},
		{Type: Print, Line: 6, Lexeme: "print", Value: nil},
		{Type: String, Line: 6, Lexeme: "hi", Value: nil},
		{Type: Semicolon, Line: 6, Lexeme: ";", Value: nil},
		{Type: RightBrace, Line: 7, Lexeme: "}", Value: nil},
		{Type: EOF, Line: 7, Lexeme: "", Value: nil},
	}},
	{name: "tc7.lox", expected: []Token{
		{Type: Number, Line: 1, Lexeme: "1", Value: nil},
		{Type: GreaterEqual, Line: 1, Lexeme: ">=", Value: nil},
		{Type: Number, Line: 1, Lexeme: "2", Value: nil},
		{Type: Semicolon, Line: 1, Lexeme: ";", Value: nil},
		{Type: EOF, Line: 3, Lexeme: "", Value: nil},
	}},
}

func TestLexerValidTCs(t *testing.T) {
	for _, tc := range validTestCases {
		src, err := os.ReadFile(fmt.Sprintf("testdata/valid/%s", tc.name))
		if err != nil {
			t.Errorf("failed reading test data at %s: %s", tc.name, err)
			t.Fail()
		}
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := Parse(string(src))
			if err != nil {
				t.Errorf("failed parsing test case: %s", err)
				t.Fail()
			}
			if len(tokens) != len(tc.expected) {
				t.Errorf("failed parsing test case: %s", err)
				t.Fail()
			}

			for i, token := range tokens {
				if token.Type != tc.expected[i].Type {
					t.Errorf("invalid token type at index %d, expected %s got %s",
						i, token.Type.String(), tc.expected[i].Type.String())
					t.Fail()
				}
				if token.Line != tc.expected[i].Line {
					t.Errorf("invalid token line at index %d, expected %d got %d",
						i, token.Line, tc.expected[i].Line)
					t.Fail()
				}
			}
		})
	}
}

type invalidTestCase struct {
	name     string
	expected error
}

var invalidTestCases = []invalidTestCase{
	// TODO
}

func TestLexerInvalidTCs(t *testing.T) {
	for _, tc := range invalidTestCases {
		_, err := os.ReadFile(fmt.Sprintf("testdata/invalid/%s", tc.name))
		if err != nil {
			t.Errorf("failed reading test data at %s: %s", tc.name, err)
			t.Fail()
		}
		t.Run(tc.name, func(t *testing.T) {
			// TODO
		})
	}
}
