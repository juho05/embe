package main

import (
	"fmt"
	"os"

	"github.com/Bananenpro/embe/parser"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stdin, "USAGE: %s <file>\n", os.Args[0])
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	check(err)
	defer file.Close()

	tokens, _, err := parser.Scan(file)
	check(err)

	fmt.Println(tokens)

	err = parser.Parse(tokens)
	check(err)
}

func check(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stdin, "ERROR: %s\n", err.Error())
	os.Exit(1)
}
