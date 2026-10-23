package lexer

var symbolsMap = map[byte]TokenType{
	'{': LeftBrace,
	'}': RightBrace,
	'(': LeftParen,
	')': RightParen,
	',': Comma,
	'.': Dot,
	';': Semicolon,
	'+': Plus,
	'-': Minus,
	'*': Star,
	'/': Slash,
	'!': Bang,
	'=': Equal,
	'<': Less,
	'>': Greater,
}

var twoCharSymbolsMap = map[string]TokenType{
	"==": EqualEqual,
	"!=": BangEqual,
	"<=": LessEqual,
	">=": GreaterEqual,
}

var keywordsMap = map[string]TokenType{
	"true":   True,
	"false":  False,
	"nil":    Nil,
	"and":    And,
	"or":     Or,
	"if":     If,
	"else":   Else,
	"for":    For,
	"while":  While,
	"var":    Var,
	"fun":    Fun,
	"class":  Class,
	"this":   This,
	"super":  Super,
	"return": Return,
	"print":  Print,
}
