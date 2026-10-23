package main

import (
	"os"

	"github.com/goncalorazevedo/glox/lexer"
	"github.com/goncalorazevedo/glox/lox"
	"github.com/goncalorazevedo/glox/parser"
)

func main() {
	var file string
	if len(os.Args) >= 2 {
		file = os.Args[1]
	}

	if file == "" {
		lox.Repl()
	} else {
		src, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}

		tokens, err := lexer.Parse(string(src))
		if err != nil {
			panic(err)
		}

		stmnts, err := parser.Parse(tokens)
		if err != nil {
			panic(err)
		}

		if err := lox.Evaluate(stmnts); err != nil {
			panic(err)
		}
	}
}
