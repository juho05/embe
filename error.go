package main

import (
	"fmt"
	"strings"

	"github.com/Bananenpro/embe/analyzer"
	"github.com/Bananenpro/embe/generator"
	"github.com/Bananenpro/embe/parser"
)

func generateErrorText(message string, lines [][]rune, start, end parser.Position, warning bool) string {
	errColor := "\x1b[4m\x1b[31m" // red underlined
	errLabel := "\x1b[31mERROR"
	if warning {
		errColor = "\x1b[4m\x1b[33m" // yellow underlined
		errLabel = "\x1b[33mWARNING"
	}

	errorLines := make([]string, 0, end.Line-start.Line+1)
	for l := start.Line; l <= end.Line; l++ {
		line := []rune(strings.TrimPrefix(strings.TrimPrefix(string(lines[l]), " "), "\t"))
		startCol := 0
		if start.Line == l {
			startCol = start.Column
		}
		endCol := len(line) - 1
		if end.Line == l {
			endCol = end.Column
		}
		startCol = startCol - (len(lines[l]) - len(line))
		endCol = endCol - (len(lines[l]) - len(line)) + 1

		if endCol > len(line) {
			endCol = len(line)
		}
		if startCol < 0 || startCol >= endCol {
			startCol = 0
		}

		errorLine := string(line[:startCol])
		errorLine = errorLine + errColor
		errorLine = errorLine + string(line[startCol:endCol])
		errorLine = errorLine + "\x1b[0m"
		errorLine = errorLine + string(line[endCol:])
		errorLine = fmt.Sprintf("\x1b[2m[%d]  \x1b[0m%s", l+1, errorLine)
		errorLines = append(errorLines, errorLine)
	}

	return fmt.Sprintf("%s\x1b[2m%s\x1b[0m\n%s\n\x1b[2m%s\x1b[0m", fmt.Sprintf("%s\x1b[0m [%d:%d]: %s\n", errLabel, start.Line+1, start.Column+1, message), strings.Repeat("-", 30), strings.Join(errorLines, "\n"), strings.Repeat("-", 30))
}

func printError(err error, lines [][]rune) {
	switch e := err.(type) {
	case parser.ScanError:
		fmt.Fprintln(stderr, generateErrorText(e.Message, lines, e.Pos, parser.Position{
			Line:   e.Pos.Line,
			Column: e.Pos.Column,
		}, false))
	case parser.ParseError:
		fmt.Fprintln(stderr, generateErrorText(e.Message, lines, e.Token.Pos, parser.Position{
			Line:   e.Token.Pos.Line,
			Column: e.Token.Pos.Column + len(e.Token.Lexeme) - 1,
		}, false))
	case analyzer.AnalyzerError:
		fmt.Fprintln(stderr, generateErrorText(e.Message, lines, e.Start, e.End, e.Warning))
	case generator.GenerateError:
		fmt.Fprintln(stderr, generateErrorText(e.Message, lines, e.Start, e.End, e.Warning))
	default:
		fmt.Fprintf(stderr, "\x1b[31mERROR\x1b[0m: %s\n", err.Error())
	}
}
