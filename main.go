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
		fmt.Fprintf(os.Stderr, "USAGE: %s <files...>\n", os.Args[0])
		os.Exit(1)
	}

	var inFileNameBase string
	results := make([]generator.GeneratorResult, len(os.Args)-1)

	var error bool
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("Compiling %s...\n", os.Args[i])
		file, err := os.Open(os.Args[i])
		check(err)
		if i == 1 {
			inFileNameBase = filepath.Base(file.Name())
		}

		tokens, lines, err := parser.Scan(file)
		file.Close()
		check(err)

		statements, errs := parser.Parse(tokens, lines)
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Fprintln(os.Stderr, err)
			}
			error = true
		}

		result := generator.GenerateBlocks(statements, lines)
		for _, w := range result.Warnings {
			fmt.Fprintln(os.Stderr, w)
		}
		if len(result.Errors) > 0 {
			for _, err := range result.Errors {
				fmt.Fprintln(os.Stderr, err)
			}
			error = true
		}
		results[i-1] = result
	}

	if error {
		os.Exit(1)
	}

	outName := strings.TrimSuffix(inFileNameBase, filepath.Ext(inFileNameBase)) + ".mblock"

	fmt.Printf("Writing output to %s...\n", outName)

	outFile, err := os.Create(outName)
	check(err)
	defer outFile.Close()
	err = generator.Package(outFile, results)
	check(err)
}

func check(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
	os.Exit(1)
}
