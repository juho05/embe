package main

import (
	"fmt"
	"os"

	"github.com/Bananenpro/embe/debug"
	"github.com/Bananenpro/embe/parser"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stdin, "USAGE: %s <file>\n", os.Args[0])
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	check(err)

	tokens, lines, err := parser.Scan(file)
	file.Close()
	check(err)

	statements, errs := parser.Parse(tokens, lines)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}

	debug.PrintAST(statements)
}

func check(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stdin, "ERROR: %s\n", err.Error())
	os.Exit(1)
}
