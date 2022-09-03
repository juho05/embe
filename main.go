package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Bananenpro/embe/generator"
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

	blocks, variables, _, errs := generator.GenerateBlocks(statements, lines)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}

	inFileNameBase := filepath.Base(file.Name())
	outName := strings.TrimSuffix(inFileNameBase, filepath.Ext(inFileNameBase)) + ".mblock"
	outFile, err := os.Create(outName)
	check(err)
	defer outFile.Close()
	err = generator.Package(outFile, blocks, variables)
	check(err)
}

func check(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stdin, "ERROR: %s\n", err.Error())
	os.Exit(1)
}
