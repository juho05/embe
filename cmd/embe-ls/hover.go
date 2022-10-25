package main

import (
	"fmt"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/Bananenpro/embe/analyzer"
	"github.com/Bananenpro/embe/parser"
)

func textDocumentHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	document, ok := getDocument(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}

	var token parser.Token
	var tokenIndex int
	for i, t := range document.tokens {
		if t.Pos.Line == int(params.Position.Line) && int(params.Position.Character) >= t.Pos.Column-1 && int(params.Position.Character) <= t.Pos.Column+len(t.Lexeme)-1 {
			token = t
			tokenIndex = i
			break
		}
	}
	if token.Type != parser.TkIdentifier {
		Warn("Hover on non identifier at %v.", params.Position)
		return nil, nil
	}
	Trace("Document hover at position: %v; token: %v", params.Position, token)

	identifierName := token.Lexeme

	var signature string

functions:
	for _, f := range document.functions {
		if int(params.Position.Line) >= f.StartLine && int(params.Position.Line) <= f.EndLine {
			for _, p := range f.Params {
				if p.Name.Lexeme == identifierName {
					signature = fmt.Sprintf("var %s: %s", p.Name.Lexeme, p.Type.DataType)
					break functions
				}
			}
		}
	}

	if signature == "" {
		if e, ok := analyzer.Events[token.Lexeme]; ok && tokenIndex > 0 && document.tokens[tokenIndex-1].Type == parser.TkAt {
			signature = e.String()
			identifierName = "@" + identifierName
		} else if f, ok := analyzer.FuncCalls[token.Lexeme]; ok {
			paramCount := getParamCount(document.tokens, tokenIndex+2)
			for _, s := range f.Signatures {
				if len(s.Params) == paramCount {
					signature = "func " + s.String()
					break
				}
			}
		} else if ef, ok := analyzer.ExprFuncCalls[token.Lexeme]; ok {
			paramCount := getParamCount(document.tokens, tokenIndex+2)
			for _, s := range ef.Signatures {
				if len(s.Params) == paramCount {
					signature = "func " + s.String()
					break
				}
			}
		} else if v, ok := analyzer.Variables[token.Lexeme]; ok {
			signature = v.String()
		} else if cv, ok := document.variables[token.Lexeme]; ok {
			signature = fmt.Sprintf("var %s: %s", cv.Name.Lexeme, cv.DataType)
		} else if l, ok := document.lists[token.Lexeme]; ok {
			signature = fmt.Sprintf("var %s: %s", l.Name.Lexeme, l.DataType)
		} else if c, ok := document.constants[token.Lexeme]; ok {
			signature = fmt.Sprintf("const %s: %s = %s", c.Name.Lexeme, c.Type, toString(c.Value))
		} else if cf, ok := document.functions[token.Lexeme]; ok {
			signature = "func " + cf.Name.Lexeme + "("
			for i, p := range cf.Params {
				if i > 0 {
					signature += ", "
				}
				signature += p.Name.Lexeme + ": " + string(p.Type.DataType)
			}
			signature += ")"
		} else if d, ok := document.defines.GetDefine(token.Lexeme, token.Pos); ok {
			signature = d.String()
		} else if ce, ok := document.events[token.Lexeme]; ok {
			signature = fmt.Sprintf("event %s", ce.Name.Lexeme)
		}
	}

	if signature == "" {
		Error("No hover signature found.")
		return nil, nil
	}
	Trace("Found signature: %s", signature)

	value := fmt.Sprintf("```embe\n%s\n```", signature)

	if docs, ok := documentation[identifierName]; ok {
		value += "\n---\n" + docs
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: value,
		},
	}, nil
}

func getParamCount(tokens []parser.Token, start int) int {
	parens := 1
	paramCount := 1
	for i := start; i < len(tokens) && parens > 0 && tokens[i].Type != parser.TkNewLine; i++ {
		switch tokens[i].Type {
		case parser.TkOpenParen:
			parens++
		case parser.TkCloseParen:
			parens--
		case parser.TkComma:
			if parens == 1 {
				paramCount++
			}
		}
	}
	if parens != 0 || (start < len(tokens) && tokens[start].Type == parser.TkCloseParen) {
		return 0
	}
	return paramCount
}

func toString(value any) string {
	if v, ok := value.(string); ok {
		return fmt.Sprintf("\"%v\"", v)
	}
	return fmt.Sprintf("%v", value)
}
