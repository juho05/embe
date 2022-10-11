package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-colorable"

	"github.com/Bananenpro/embe/analyzer"
	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/generator"
	"github.com/Bananenpro/embe/parser"
)

var stderr = colorable.NewColorableStderr()

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "USAGE: %s <files...>\n", os.Args[0])
		os.Exit(1)
	}

	var inFileNameBase string

	allBlocks := make([]map[string]*blocks.Block, 0, len(os.Args)-1)
	allDefinitions := make([]analyzer.Definitions, 0, len(os.Args)-1)

	var error bool
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("Compiling %s...\n", os.Args[i])
		file, err := os.Open(os.Args[i])
		if err != nil {
			printError(err)
			error = true
			continue
		}
		if i == 1 {
			inFileNameBase = filepath.Base(file.Name())
		}

		tokens, lines, err := parser.Scan(file)
		file.Close()
		if err != nil {
			fmt.Fprintln(stderr, err)
			error = true
			continue
		}

		statements, errs := parser.Parse(tokens, lines)
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Fprintln(stderr, err)
			}
			error = true
			continue
		}

		statements, analyzerResult := analyzer.Analyze(statements, lines)
		for _, w := range analyzerResult.Warnings {
			fmt.Fprintln(stderr, w)
		}
		if len(analyzerResult.Errors) > 0 {
			for _, err := range analyzerResult.Errors {
				fmt.Fprintln(stderr, err)
			}
			error = true
			continue
		}

		blocks, errs := generator.GenerateBlocks(statements, analyzerResult.Definitions, lines)
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Fprintln(stderr, err)
			}
			error = true
			continue
		}
		allBlocks = append(allBlocks, blocks)
		allDefinitions = append(allDefinitions, analyzerResult.Definitions)
	}

	if error {
		os.Exit(1)
	}

	outName := strings.TrimSuffix(inFileNameBase, filepath.Ext(inFileNameBase)) + ".mblock"

	fmt.Printf("Writing output to %s...\n", outName)

	outFile, err := os.Create(outName)
	if err != nil {
		printError(err)
		os.Exit(1)
	}
	defer outFile.Close()
	err = generator.Package(outFile, allBlocks, allDefinitions)
	if err != nil {
		printError(err)
		os.Exit(1)
	}
}

func printError(err error) {
	fmt.Fprintf(stderr, "\x1b[31mERROR\x1b[0m: %s\n", err.Error())
}
