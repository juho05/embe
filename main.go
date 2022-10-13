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

var (
	stderr         = colorable.NewColorableStderr()
	version string = "dev"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(stderr, "Compile embe source code to .mblock files.\n\n")
		fmt.Fprintf(stderr, "USAGE:\n  %s <files...>\n\n", os.Args[0])
		fmt.Fprintln(stderr, "COMMANDS:")
		fmt.Fprintln(stderr, "  version    print the embe version number")
		fmt.Fprintln(stderr, "  update     update embe to the latest release version")
		fmt.Fprintln(stderr, "  uninstall  uninstall embe")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "version":
		printVersion()
	case "update":
		update()
	case "uninstall":
		uninstall()
	default:
		run()
	}
}

func run() {
	versionCheck(true, false)

	var inFileNameBase string

	allBlocks := make([]map[string]*blocks.Block, 0, len(os.Args)-1)
	allDefinitions := make([]analyzer.Definitions, 0, len(os.Args)-1)

	var error bool
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("Compiling %s...\n", os.Args[i])
		file, err := os.Open(os.Args[i])
		if err != nil {
			printError(err, nil)
			error = true
			continue
		}
		if i == 1 {
			inFileNameBase = filepath.Base(file.Name())
		}

		tokens, lines, err := parser.Scan(file)
		file.Close()
		if err != nil {
			printError(err, lines)
			error = true
			continue
		}

		statements, errs := parser.Parse(tokens, lines)
		if len(errs) > 0 {
			for _, err := range errs {
				printError(err, lines)
			}
			error = true
			continue
		}

		statements, analyzerResult := analyzer.Analyze(statements, lines)
		for _, w := range analyzerResult.Warnings {
			printError(w, lines)
		}
		if len(analyzerResult.Errors) > 0 {
			for _, err := range analyzerResult.Errors {
				printError(err, lines)
			}
			error = true
			continue
		}

		blocks, errs := generator.GenerateBlocks(statements, analyzerResult.Definitions, lines)
		if len(errs) > 0 {
			for _, err := range errs {
				printError(err, lines)
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
		printError(err, nil)
		os.Exit(1)
	}
	defer outFile.Close()
	err = generator.Package(outFile, allBlocks, allDefinitions)
	if err != nil {
		printError(err, nil)
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Println("embe", version)
}

func uninstall() {
	panic("not implemented")
}
