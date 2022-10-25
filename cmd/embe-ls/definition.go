package main

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/Bananenpro/embe/parser"
)

func textDocumentDefinition(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	document, ok := getDocument(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}

	var token parser.Token
	for _, t := range document.tokens {
		if t.Pos.Line == int(params.Position.Line) && int(params.Position.Character) >= t.Pos.Column-1 && int(params.Position.Character) <= t.Pos.Column+len(t.Lexeme)-1 {
			token = t
			break
		}
	}
	if token.Type != parser.TkIdentifier {
		Warn("Goto definition on non identifier at %v.", params.Position)
		return nil, nil
	}
	Trace("Goto definition at position: %v; token: %v", params.Position, token)
	identifierName := token.Lexeme

	var start parser.Position
	var end parser.Position

functions:
	for _, f := range document.functions {
		if int(params.Position.Line) >= f.StartLine && int(params.Position.Line) <= f.EndLine {
			for _, p := range f.Params {
				if p.Name.Lexeme == identifierName {
					start = p.Name.Pos
					end = p.Name.EndPos
					break functions
				}
			}
		}
	}

	if start == (parser.Position{}) || end == (parser.Position{}) {
		if e, ok := document.events[identifierName]; ok {
			start = e.Name.Pos
			end = e.Name.EndPos
		} else if f, ok := document.functions[identifierName]; ok {
			start = f.Name.Pos
			end = f.Name.EndPos
		} else if v, ok := document.variables[identifierName]; ok {
			start = v.Name.Pos
			end = v.Name.EndPos
		} else if l, ok := document.lists[identifierName]; ok {
			start = l.Name.Pos
			end = l.Name.EndPos
		} else if c, ok := document.constants[identifierName]; ok {
			start = c.Name.Pos
			end = c.Name.EndPos
		} else if d, ok := document.defines.GetDefine(identifierName, token.Pos); ok {
			start = d.Name.Pos
			end = d.Name.EndPos
		} else {
			return nil, nil
		}
	}

	return &protocol.Location{
		URI: document.uri,
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
