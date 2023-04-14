package main

import (
	"os"
	"path/filepath"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/juho05/embe/parser"
)

func textDocumentDefinition(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
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

	if token.Type == parser.TkLiteral && token.DataType == parser.DTString && tokenIndex > 0 && document.tokens[tokenIndex-1].Type == parser.TkPreprocessor && document.tokens[tokenIndex-1].Lexeme == "#include" {
		target := filepath.Join(filepath.Dir(document.path), token.Literal.(string))
		if filepath.Ext(target) != ".mb" {
			target += ".mb"
		}
		if _, ok := getDocument(pathToURI(target)); !ok {
			if _, err := os.Stat(target); err != nil {
				return nil, nil
			}
		}
		return &protocol.Location{
			URI: pathToURI(target),
		}, nil
	}

	if token.Type != parser.TkIdentifier {
		Warn("Goto definition on non identifier at %v.", params.Position)
		return nil, nil
	}
	Trace("Goto definition at position: %v; token: %v", params.Position, token)
	identifierName := token.Lexeme

	var start parser.Position
	var end parser.Position
	var path string

functions:
	for _, f := range document.functions {
		if int(params.Position.Line) >= f.StartLine && int(params.Position.Line) <= f.EndLine {
			for _, p := range f.Params {
				if p.Name.Lexeme == identifierName {
					start = p.Name.Pos
					end = p.Name.EndPos
					path = p.Name.Pos.Path
					break functions
				}
			}
		}
	}

	if start == (parser.Position{}) || end == (parser.Position{}) {
		if e, ok := document.events[identifierName]; ok {
			start = e.Name.Pos
			end = e.Name.EndPos
			path = e.Name.Pos.Path
		} else if f, ok := document.functions[identifierName]; ok {
			start = f.Name.Pos
			end = f.Name.EndPos
			path = f.Name.Pos.Path
		} else if v, ok := document.variables[identifierName]; ok {
			start = v.Name.Pos
			end = v.Name.EndPos
			path = v.Name.Pos.Path
		} else if l, ok := document.lists[identifierName]; ok {
			start = l.Name.Pos
			end = l.Name.EndPos
			path = l.Name.Pos.Path
		} else if c, ok := document.constants[identifierName]; ok {
			start = c.Name.Pos
			end = c.Name.EndPos
			path = c.Name.Pos.Path
		} else if d, ok := document.defines.GetDefine(identifierName, token.Pos); ok {
			start = d.Name.Pos
			end = d.Name.EndPos
			path = d.Name.Pos.Path
		} else {
			return nil, nil
		}
	}

	return &protocol.Location{
		URI: pathToURI(path),
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      uint32(start.Line),
				Character: uint32(start.Column),
			},
			End: protocol.Position{
				Line:      uint32(end.Line),
				Character: uint32(end.Column),
			},
		},
	}, nil
}
