package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

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
		repl()
	} else {
		src, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}

		tokens, err := lexer.Parse(string(src))
		if err != nil {
			fmt.Println(err)
			return
		}

		stmnts, err := parser.Parse(tokens)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := lox.Evaluate(stmnts); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func repl() {
	trimLinePrefixFunc := func(s string) string {
		linePrefix := regexp.MustCompile(`^line \d+:\s*`)
		return linePrefix.ReplaceAllString(s, "")
	}

	reader := bufio.NewReader(os.Stdin)
	// FIXME add support for a repl session to preserve environmnent between prompts
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		tokens, err := lexer.Parse(text)
		if err != nil {
			fmt.Println(trimLinePrefixFunc(err.Error()))
			continue
		}
		stmnts, err := parser.Parse(tokens)
		if err != nil {
			fmt.Println(trimLinePrefixFunc(err.Error()))
			continue
		}
		if err := lox.Evaluate(stmnts); err != nil {
			fmt.Println(trimLinePrefixFunc(err.Error()))
			continue
		}
	}
}
