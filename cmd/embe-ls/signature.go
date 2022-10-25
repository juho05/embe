package main

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/Bananenpro/embe/analyzer"
	"github.com/Bananenpro/embe/parser"
)

func textDocumentSignatureHelp(context *glsp.Context, params *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	document, ok := getDocument(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}

	Trace("Signature help at %v.", params.Position)

	pos := params.Position
	startPos := pos
	startPos.Character = 0
	endPos := startPos.EndOfLineIn(document.content)

	line := document.content[startPos.IndexIn(document.content):endPos.IndexIn(document.content)]

	parenIndices := make([]int, 0, 5)
	for i, c := range line {
		if i >= int(pos.Character) {
			break
		}
		if c == '(' {
			parenIndices = append(parenIndices, i)
		} else if c == ')' && len(parenIndices) > 0 {
			parenIndices = parenIndices[:len(parenIndices)-1]
		}
	}
	if len(parenIndices) == 0 {
		Warn("Requesting signature help but not in a function at %v. Line: %s", pos, line)
		return nil, nil
	}

	var identifier parser.Token
	var identifierIndex int
	for i, t := range document.tokens {
		if i == 0 {
			continue
		}
		if t.Pos.Line == int(params.Position.Line) && t.Pos.Column == parenIndices[len(parenIndices)-1] {
			identifier = document.tokens[i-1]
			identifierIndex = i
			break
		}
	}

	if (identifier == parser.Token{}) {
		Error("Couldn't find token under cursor.")
		return nil, nil
	}

	var signatures []analyzer.Signature
	if f, ok := analyzer.ExprFuncCalls[identifier.Lexeme]; ok {
		signatures = f.Signatures
	} else if f, ok := analyzer.FuncCalls[identifier.Lexeme]; ok {
		signatures = f.Signatures
	} else if f, ok := document.functions[identifier.Lexeme]; ok {
		params := make([]analyzer.Param, 0)
		for _, p := range f.Params {
			params = append(params, analyzer.Param{
				Name: p.Name.Lexeme,
				Type: p.Type.DataType,
			})
		}
		signatures = []analyzer.Signature{
			{
				FuncName: f.Name.Lexeme,
				Params:   params,
			},
		}
	} else {
		Error("Couldn't find signature for token: %s", identifier)
		return nil, nil
	}

	parens := 1
	paramCount := 1
	var paramIndex uint32
	for i := identifierIndex + 2; i < len(document.tokens) && parens > 0; i++ {
		token := document.tokens[i]
		switch token.Type {
		case parser.TkOpenParen:
			parens++
		case parser.TkCloseParen:
			parens--
		case parser.TkComma:
			paramCount++
			if token.Pos.Line == int(pos.Line) && token.Pos.Column <= int(pos.Character) {
				paramIndex++
			}
		}
	}
	if identifierIndex+2 < len(document.tokens) && document.tokens[identifierIndex+2].Type == parser.TkCloseParen {
		paramCount = 0
	}

	var activeSignature uint32
	for i, s := range signatures {
		if len(s.Params) == paramCount {
			activeSignature = uint32(i)
			break
		}
	}

	signatureInformation := make([]protocol.SignatureInformation, len(signatures))
	for i, s := range signatures {
		parameters := make([]protocol.ParameterInformation, len(s.Params))
		for j, p := range s.Params {
			parameters[j] = protocol.ParameterInformation{
				Label: p.Name + ": " + string(p.Type),
			}
		}
		signatureInformation[i] = protocol.SignatureInformation{
			Label:      s.String(),
			Parameters: parameters,
		}
	}

	Trace("Sending signature help: activeSignature: %d; activeParameter: %d, signatures: %v", activeSignature, paramIndex, signatureInformation)

	return &protocol.SignatureHelp{
		Signatures:      signatureInformation,
		ActiveSignature: &activeSignature,
		ActiveParameter: &paramIndex,
	}, nil
}
