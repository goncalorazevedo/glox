package lexer

import (
	"errors"
	"fmt"
	"strconv"
)

func Parse(src string) ([]Token, error) {
	lexer := newLexer(src)
	return lexer.parse()
}

type lexer struct {
	src   string
	index int
	line  int
}

func newLexer(src string) lexer {
	return lexer{
		src:   src,
		index: 0,
		line:  1,
	}
}

var lexEOF = errors.New("lexer reached end of file")

func (l *lexer) parse() ([]Token, error) {
	tokens := []Token{}
	for {
		token, err := l.parseToken()
		if err != nil {
			if errors.Is(err, lexEOF) {
				break
			}
			return nil, fmt.Errorf("failed parsing source: %w", err)
		}
		tokens = append(tokens, token)
	}
	tokens = append(tokens, l.emitToken(EOF, ""))
	return tokens, nil
}

func (l *lexer) parseToken() (Token, error) {
	for {
		if l.atEnd() {
			return Token{}, lexEOF
		}
		switch c := l.src[l.index]; c {
		case '\r', '\t', ' ':
			l.index++
		case '\n':
			l.line++
			l.index++
		case '}', '{', '(', ')', ',', '.', ';', '+', '-', '*':
			l.index++
			return l.emitToken(l.lookupSymbolType(c), string(c)), nil
		case '/':
			l.index++
			// checking for comment
			if l.atEnd() || l.src[l.index] != '/' {
				return l.emitToken(Slash, "/"), nil
			} else {
				// if comment discard until LF
				for {
					l.index++
					if l.atEnd() {
						break
					} else if l.src[l.index] == '\n' {
						l.line++
						l.index++ // discard the checked LF
						break
					}
				}
			}
		case '!', '=', '<', '>':
			l.index++
			if l.atEnd() || l.src[l.index] != '=' {
				return l.emitToken(l.lookupSymbolType(c), string(c)), nil
			}
			symbol := string(l.src[l.index-1 : l.index+1])
			l.index++
			return l.emitToken(l.lookupTwoCharSymbolType(symbol), symbol), nil
		case '"':
			start := l.index + 1
			for {
				l.index++
				if l.atEnd() {
					return Token{}, fmt.Errorf(`line %d: invalid string literal, missing closing '"'`, l.line)
				}
				if l.src[l.index] == '"' {
					l.index++
					break
				}
			}
			literal := l.src[start : l.index-1]
			return l.emitTokenWithValue(String, string(literal), literal), nil
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			start := l.index
			var isFloat bool = false
			for {
				l.index++
				if l.atEnd() {
					break
				}
				if l.src[l.index] == '.' {
					if isFloat {
						return Token{}, fmt.Errorf(`line %d: invalid number syntax "%s"`, l.line, l.src[start:l.index+1])
					}
					isFloat = true
					continue
				}
				if l.src[l.index] < '0' || l.src[l.index] > '9' {
					break
				}
			}
			literal := l.src[start:l.index]
			num, err := strconv.ParseFloat(literal, 64)
			if err != nil {
				return Token{}, fmt.Errorf(`line %d: unsupported float syntax "%s": %w`, l.line, literal, err)
			}
			return l.emitTokenWithValue(Number, literal, num), nil
		default:
			if !isAlpha(c) {
				return Token{}, fmt.Errorf("line %d: unexpected character %c", l.line, c)
			}
			start := l.index
			for {
				l.index++
				if l.atEnd() || !isAlpha(l.src[l.index]) {
					break
				}
			}
			word := l.src[start:l.index]
			tokenType := l.lookupKeywordSymbolType(word)
			var value any
			switch tokenType {
			case True:
				value = true
			case False:
				value = false
			case Nil:
				value = struct{}{}
			}
			return l.emitTokenWithValue(tokenType, word, value), nil
		}
	}
}

func isAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func (l *lexer) emitToken(typ TokenType, lexem string) Token {
	return Token{
		Type:   typ,
		Line:   l.line,
		Lexeme: lexem,
		Value:  nil,
	}
}

func (l *lexer) emitTokenWithValue(typ TokenType, lexem string, value any) Token {
	return Token{
		Type:   typ,
		Line:   l.line,
		Lexeme: lexem,
		Value:  value,
	}
}

func (l *lexer) lookupSymbolType(c byte) TokenType {
	t, valid := symbolsMap[c]
	if !valid {
		panic(fmt.Sprintf("parsed an invalid symbol '%s' in line %d", string(c), l.line))
	}
	return t
}

func (l *lexer) lookupTwoCharSymbolType(symbol string) TokenType {
	if len(symbol) != 2 {
		panic(fmt.Sprintf("parsed and invalid two char symbol '%s' in line %d", symbol, l.line))
	}
	t, valid := twoCharSymbolsMap[symbol]
	if !valid {
		panic(fmt.Sprintf("parsed an invalid two char symbol symbol '%s' in line %d", symbol, l.line))
	}
	return t
}

func (l *lexer) lookupKeywordSymbolType(symbol string) TokenType {
	t, valid := keywordsMap[symbol]
	if !valid {
		return Identifier
	}
	return t
}

func (l *lexer) atEnd() bool {
	return l.index >= len(l.src)
}
