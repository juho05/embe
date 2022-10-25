package main

import (
	"bytes"
	"strings"
	"sync"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/Bananenpro/embe/analyzer"
	"github.com/Bananenpro/embe/generator"
	"github.com/Bananenpro/embe/parser"
)

type Document struct {
	uri         protocol.DocumentUri
	content     string
	tokens      []parser.Token
	changed     bool
	diagnostics []protocol.Diagnostic
	defines     *parser.Defines
	variables   map[string]*analyzer.Variable
	lists       map[string]*analyzer.List
	constants   map[string]*analyzer.Constant
	functions   map[string]*analyzer.Function
	events      map[string]*analyzer.CustomEvent
}

var documents sync.Map

func (d *Document) validate(notify glsp.NotifyFunc) {
	if !d.changed {
		return
	}
	d.changed = false

	if !strings.HasSuffix(d.content, "\n") {
		d.content += "\n"
	}

	Trace("Validating document...")

	defer d.sendDiagnostics(notify)

	severityWarning := protocol.DiagnosticSeverityWarning
	severityError := protocol.DiagnosticSeverityError

	d.diagnostics = d.diagnostics[:0]

	tokens, _, errs := parser.Scan(bytes.NewBufferString(d.content))
	if len(errs) > 0 {
		for _, err := range errs {
			if e, ok := err.(parser.ScanError); ok {
				d.diagnostics = append(d.diagnostics, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{
							Line:      uint32(e.Pos.Line),
							Character: uint32(e.Pos.Column),
						},
						End: protocol.Position{
							Line:      uint32(e.Pos.Line),
							Character: uint32(e.Pos.Column + 1),
						},
					},
					Severity: &severityError,
					Message:  e.Message,
				})
			} else {
				Error("Failed to scan '%s': %s", d.uri, err)
			}
		}
		return
	}
	d.tokens = make([]parser.Token, len(tokens))
	copy(d.tokens, tokens)

	tokens, defines, errs := parser.Preprocess(tokens)
	if len(errs) > 0 {
		for _, err := range errs {
			if e, ok := err.(parser.ParseError); ok {
				d.diagnostics = append(d.diagnostics, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{
							Line:      uint32(e.Token.Pos.Line),
							Character: uint32(e.Token.Pos.Column),
						},
						End: protocol.Position{
							Line:      uint32(e.Token.Pos.Line),
							Character: uint32(e.Token.Pos.Column + len(e.Token.Lexeme)),
						},
					},
					Severity: &severityError,
					Message:  e.Message,
				})
			} else {
				Error("Failed to preprocess '%s': %s", d.uri, err)
			}
		}
		return
	}
	d.defines = defines

	statements, errs := parser.Parse(tokens)
	if len(errs) > 0 {
		for _, err := range errs {
			if e, ok := err.(parser.ParseError); ok {
				d.diagnostics = append(d.diagnostics, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{
							Line:      uint32(e.Token.Pos.Line),
							Character: uint32(e.Token.Pos.Column),
						},
						End: protocol.Position{
							Line:      uint32(e.Token.Pos.Line),
							Character: uint32(e.Token.Pos.Column + len(e.Token.Lexeme)),
						},
					},
					Severity: &severityError,
					Message:  e.Message,
				})
			} else {
				Error("Failed to parse '%s': %s", d.uri, err)
			}
		}
		return
	}

	statements, analyzerResult := analyzer.Analyze(statements)
	for _, warning := range analyzerResult.Warnings {
		if w, ok := warning.(analyzer.AnalyzerError); ok {
			d.diagnostics = append(d.diagnostics, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(w.Start.Line),
						Character: uint32(w.Start.Column),
					},
					End: protocol.Position{
						Line:      uint32(w.End.Line),
						Character: uint32(w.End.Column + 1),
					},
				},
				Severity: &severityWarning,
				Message:  w.Message,
			})
		}
	}
	if len(analyzerResult.Errors) > 0 {
		for _, err := range analyzerResult.Errors {
			if e, ok := err.(analyzer.AnalyzerError); ok {
				d.diagnostics = append(d.diagnostics, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{
							Line:      uint32(e.Start.Line),
							Character: uint32(e.Start.Column),
						},
						End: protocol.Position{
							Line:      uint32(e.End.Line),
							Character: uint32(e.End.Column + 1),
						},
					},
					Severity: &severityError,
					Message:  e.Message,
				})
			} else {
				Error("Failed to parse '%s': %s", d.uri, err)
			}
		}
		return
	}
	d.variables = analyzerResult.Definitions.Variables
	d.lists = analyzerResult.Definitions.Lists
	d.constants = analyzerResult.Definitions.Constants
	d.functions = analyzerResult.Definitions.Functions
	d.events = analyzerResult.Definitions.Events

	_, errs = generator.GenerateBlocks(statements, analyzerResult.Definitions)
	if len(errs) > 0 {
		for _, err := range errs {
			if e, ok := err.(generator.GenerateError); ok {
				d.diagnostics = append(d.diagnostics, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{
							Line:      uint32(e.Start.Line),
							Character: uint32(e.Start.Column),
						},
						End: protocol.Position{
							Line:      uint32(e.End.Line),
							Character: uint32(e.End.Column + 1),
						},
					},
					Severity: &severityError,
					Message:  e.Message,
				})
			} else {
				Error("Failed to generate blocks for '%s': %s", d.uri, err)
			}
		}
		return
	}
}

func (d *Document) sendDiagnostics(notify glsp.NotifyFunc) {
	Trace("Sending diagnostics...")
	notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
		URI:         d.uri,
		Diagnostics: d.diagnostics,
	})
}

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	Trace("Document did open: %s", params.TextDocument.URI)
	document := &Document{
		uri:         params.TextDocument.URI,
		content:     params.TextDocument.Text,
		tokens:      make([]parser.Token, 0),
		changed:     true,
		diagnostics: make([]protocol.Diagnostic, 0),
		defines:     parser.NewDefines(),
		variables:   make(map[string]*analyzer.Variable),
		lists:       make(map[string]*analyzer.List),
		constants:   make(map[string]*analyzer.Constant),
		functions:   make(map[string]*analyzer.Function),
		events:      make(map[string]*analyzer.CustomEvent),
	}
	documents.Store(params.TextDocument.URI, document)
	go document.validate(context.Notify)
	return nil
}

func textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	if document, ok := getDocument(params.TextDocument.URI); ok {
		Trace("Document did change: %s", document.uri)
		content := document.content
		for _, change := range params.ContentChanges {
			if c, ok := change.(protocol.TextDocumentContentChangeEvent); ok {
				start, end := c.Range.IndexesIn(content)
				content = content[:start] + c.Text + content[end:]
				Trace("Applied change type 'partial'.")
			} else if c, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
				content = c.Text
				Trace("Applied change type 'whole'.")
			}
		}
		document.content = content
		document.changed = len(params.ContentChanges) > 0
		go document.validate(context.Notify)
	}
	return nil
}

func textDocumentDidClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	Trace("Document did close: %s", params.TextDocument.URI)
	_, ok := documents.LoadAndDelete(params.TextDocument.URI)
	if ok {
		go context.Notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
			URI:         params.TextDocument.URI,
			Diagnostics: make([]protocol.Diagnostic, 0),
		})
	}
	return nil
}

func getDocument(uri protocol.DocumentUri) (*Document, bool) {
	doc, ok := documents.Load(uri)
	if !ok {
		return nil, false
	}
	return doc.(*Document), true
}
