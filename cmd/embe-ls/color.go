package main

import (
	"fmt"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/juho05/embe/parser"
)

var colors = map[string]protocol.Color{
	"gray": {
		Red:   0.61,
		Green: 0.61,
		Blue:  0.61,
	},
	"red": {
		Red:   0.81,
		Green: 0.01,
		Blue:  0.1,
	},
	"orange": {
		Red:   0.96,
		Green: 0.65,
		Blue:  0.14,
	},
	"yellow": {
		Red:   0.97,
		Green: 0.91,
		Blue:  0.11,
	},
	"green": {
		Red:   0.49,
		Green: 0.83,
		Blue:  0.13,
	},
	"cyan": {
		Red:   0.31,
		Green: 0.83,
		Blue:  0.76,
	},
	"blue": {
		Red:   0.29,
		Green: 0.56,
		Blue:  0.89,
	},
	"magenta": {
		Red:   0.74,
		Green: 0.06,
		Blue:  0.88,
	},
	"white": {
		Red:   1,
		Green: 1,
		Blue:  1,
	},
}

func textDocumentColor(context *glsp.Context, params *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	document, ok := getDocument(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}

	Trace("Collecting color information...")

	colorInformation := make([]protocol.ColorInformation, 0)
	for _, t := range document.tokens {
		if t.Type != parser.TkLiteral || t.DataType != parser.DTString {
			continue
		}
		color := t.Literal.(string)
		if c, ok := colors[color]; ok {
			colorInformation = append(colorInformation, protocol.ColorInformation{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(t.Pos.Line),
						Character: uint32(t.Pos.Column),
					},
					End: protocol.Position{
						Line:      uint32(t.Pos.Line),
						Character: uint32(t.Pos.Column + len(t.Lexeme)),
					},
				},
				Color: protocol.Color{
					Red:   c.Red,
					Green: c.Green,
					Blue:  c.Blue,
					Alpha: 1,
				},
			})
			continue
		}
		var r int
		var g int
		var b int
		_, err := fmt.Sscanf(color, "#%02x%02x%02x", &r, &g, &b)
		if err != nil {
			continue
		}
		colorInformation = append(colorInformation, protocol.ColorInformation{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(t.Pos.Line),
					Character: uint32(t.Pos.Column),
				},
				End: protocol.Position{
					Line:      uint32(t.Pos.Line),
					Character: uint32(t.Pos.Column + len(t.Lexeme)),
				},
			},
			Color: protocol.Color{
				Red:   float32(r) / 255,
				Green: float32(g) / 255,
				Blue:  float32(b) / 255,
				Alpha: 1,
			},
		})
	}
	Trace("Collected color information.")

	return colorInformation, nil
}

func textDocumentColorPresentation(context *glsp.Context, params *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	r := int(params.Color.Red * 255)
	g := int(params.Color.Green * 255)
	b := int(params.Color.Blue * 255)
	label := fmt.Sprintf("\"#%02x%02x%02x\"", r, g, b)
	Trace("Color presentation. Input: rgb(%f, %f, %f), Output: %s", params.Color.Red, params.Color.Green, params.Color.Blue, label)
	return []protocol.ColorPresentation{{
		Label: label,
	}}, nil
}
