package lox

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/goncalorazevedo/glox/lexer"
	"github.com/goncalorazevedo/glox/parser"
)

func Repl() {
	reader := bufio.NewReader(os.Stdin)
	interp := interpreter{
		env: newEnvironment(nil),
	}
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		tokens, err := lexer.Parse(text)
		if err != nil {
			fmt.Println(trimLinePrefix(err.Error()))
			continue
		}
		stmnts, err := parser.Parse(tokens)
		if err != nil {
			fmt.Println(trimLinePrefix(err.Error()))
			continue
		}
		for _, stmnt := range stmnts {
			if err := interp.evaluateStmnt(stmnt); err != nil {
				fmt.Println(trimLinePrefix(err.Error()))
				continue
			}
		}
	}
}

func trimLinePrefix(s string) string {
	linePrefix := regexp.MustCompile(`^line \d+:\s*`)
	return linePrefix.ReplaceAllString(s, "")
}
