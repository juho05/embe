package main

import (
	"bytes"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/juho05/embe/analyzer"
	"github.com/juho05/embe/generator"
	"github.com/juho05/embe/parser"
)

type Document struct {
	uri        protocol.DocumentUri
	path       string
	content    string
	tokens     []parser.Token
	validating bool
	defines    *parser.Defines
	variables  map[string]*analyzer.Variable
	lists      map[string]*analyzer.List
	constants  map[string]*analyzer.Constant
	functions  map[string]*analyzer.Function
	events     map[string]*analyzer.CustomEvent
}

var documents sync.Map

var (
	// document path -> outer document path
	innerDocuments     = make(map[string]string, 0)
	innerDocumentsLock sync.RWMutex
)

type reader struct {
	content io.Reader
}

func (r *reader) Read(p []byte) (n int, err error) {
	return r.content.Read(p)
}

func (r *reader) Close() error {
	return nil
}

func (d *Document) validate(notify glsp.NotifyFunc) {
	if d.validating {
		return
	}
	d.validating = true
	defer func() { d.validating = false }()

	if !strings.HasSuffix(d.content, "\n") {
		d.content += "\n"
	}

	Trace("Validating document %s...", d.uri)

	severityWarning := protocol.DiagnosticSeverityWarning
	severityError := protocol.DiagnosticSeverityError

	innerDocumentsLock.RLock()
	diagnostics := make(map[string][]protocol.Diagnostic, 1+len(innerDocuments))
	diagnostics[d.path] = make([]protocol.Diagnostic, 0, 5)
	innerDocumentsLock.RUnlock()

	var errs []error
	var defines *parser.Defines
	var files map[string][][]rune
	var statements []parser.Stmt
	var analyzerResult analyzer.AnalyzerResult

	tokens, _, errs := parser.Scan(bytes.NewBufferString(d.content), d.path)
	if len(errs) > 0 {
		for _, err := range errs {
			if e, ok := err.(parser.ScanError); ok {
				if _, ok := diagnostics[e.Pos.Path]; !ok {
					diagnostics[e.Pos.Path] = make([]protocol.Diagnostic, 0, 5)
				}
				diagnostics[e.Pos.Path] = append(diagnostics[e.Pos.Path], protocol.Diagnostic{
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
		goto diagnostics
	}
	d.tokens = make([]parser.Token, len(tokens))
	copy(d.tokens, tokens)

	tokens, files, defines, _, errs = parser.Preprocess(tokens, d.path, func(name string) (io.ReadCloser, error) {
		if doc, ok := getDocument(pathToURI(name)); ok {
			return &reader{
				content: bytes.NewReader([]byte(doc.content)),
			}, nil
		}
		return os.Open(name)
	}, nil, nil)
	for f := range files {
		if _, ok := diagnostics[f]; !ok {
			diagnostics[f] = make([]protocol.Diagnostic, 0, 5)
		}
	}
	if len(errs) > 0 {
		for _, err := range errs {
			if e, ok := err.(parser.ParseError); ok {
				diagnostics[e.Token.Pos.Path] = append(diagnostics[e.Token.Pos.Path], protocol.Diagnostic{
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
		goto diagnostics
	}
	d.defines = defines

	statements, errs = parser.Parse(tokens)
	if len(errs) > 0 {
		for _, err := range errs {
			if e, ok := err.(parser.ParseError); ok {
				diagnostics[e.Token.Pos.Path] = append(diagnostics[e.Token.Pos.Path], protocol.Diagnostic{
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
		goto diagnostics
	}

	statements, analyzerResult = analyzer.Analyze(statements)
	for _, warning := range analyzerResult.Warnings {
		if w, ok := warning.(analyzer.AnalyzerError); ok {
			diagnostics[w.Start.Path] = append(diagnostics[w.Start.Path], protocol.Diagnostic{
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
				diagnostics[e.Start.Path] = append(diagnostics[e.Start.Path], protocol.Diagnostic{
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
		goto diagnostics
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
				diagnostics[e.Start.Path] = append(diagnostics[e.Start.Path], protocol.Diagnostic{
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
		goto diagnostics
	}

diagnostics:
	innerDocumentsLock.Lock()
	if _, ok := innerDocuments[d.path]; ok {
		innerDocumentsLock.Unlock()
		return
	}
	for f, ds := range diagnostics {
		sendDiagnostics(notify, pathToURI(f), ds)
		if runtime.GOOS == "windows" {
			f = strings.ToLower(f)
		}
		if f != d.path {
			innerDocuments[f] = d.path
		}
	}
	for inner, outer := range innerDocuments {
		if outer == d.path {
			if _, ok := diagnostics[inner]; !ok {
				delete(innerDocuments, inner)
				if d, ok := getDocument(inner); ok {
					go d.validate(notify)
				}
			}
		}
	}
	innerDocumentsLock.Unlock()
}

func pathToURI(path string) string {
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			Error(err.Error())
			return path
		}
	}
	path = url.PathEscape(filepath.ToSlash(path))
	if runtime.GOOS == "windows" {
		path = strings.ToLower(path)
	}
	if !strings.HasPrefix(path, "/") {
		return "file:///" + path
	}
	return "file://" + path
}

func sendDiagnostics(notify glsp.NotifyFunc, uri string, diagnostics []protocol.Diagnostic) {
	Trace("Sending diagnostics for %s: %v", uri, diagnostics)
	notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
}

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	Trace("Document did open: %s", params.TextDocument.URI)
	path := strings.TrimPrefix(params.TextDocument.URI, "file://")
	p, err := url.PathUnescape(path)
	if err == nil {
		path = p
	}
	if runtime.GOOS == "windows" {
		path = strings.ReplaceAll(path, "/", "\\")
		path = strings.TrimPrefix(path, "\\")
		path = strings.ToLower(path)
		params.TextDocument.URI = strings.ToLower(params.TextDocument.URI)
	}
	document := &Document{
		uri:        params.TextDocument.URI,
		path:       path,
		content:    params.TextDocument.Text,
		tokens:     make([]parser.Token, 0),
		validating: false,
		defines:    parser.NewDefines(),
		variables:  make(map[string]*analyzer.Variable),
		lists:      make(map[string]*analyzer.List),
		constants:  make(map[string]*analyzer.Constant),
		functions:  make(map[string]*analyzer.Function),
		events:     make(map[string]*analyzer.CustomEvent),
	}
	documents.Store(params.TextDocument.URI, document)

	innerDocumentsLock.RLock()
	if outer, ok := innerDocuments[document.path]; ok {
		if d, ok := getDocument(pathToURI(outer)); ok {
			go d.validate(context.Notify)
			innerDocumentsLock.RUnlock()
			return nil
		}
	}
	innerDocumentsLock.RUnlock()
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

		innerDocumentsLock.RLock()
		if outer, ok := innerDocuments[document.path]; ok {
			if d, ok := getDocument(pathToURI(outer)); ok {
				go d.validate(context.Notify)
			} else {
				Warn("Cannot validate parent document %s: not loaded", outer)
			}
		}
		innerDocumentsLock.RUnlock()
		go document.validate(context.Notify)
	}
	return nil
}

func textDocumentDidClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	Trace("Document did close: %s", params.TextDocument.URI)
	if runtime.GOOS == "windows" {
		params.TextDocument.URI = strings.ToLower(params.TextDocument.URI)
	}
	d, ok := documents.LoadAndDelete(params.TextDocument.URI)
	if ok {
		innerDocumentsLock.Lock()
		delete(innerDocuments, d.(*Document).path)
		for inner, outer := range innerDocuments {
			if outer == d.(*Document).path {
				delete(innerDocuments, inner)
			}
		}
		innerDocumentsLock.Unlock()
		go context.Notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
			URI:         params.TextDocument.URI,
			Diagnostics: make([]protocol.Diagnostic, 0),
		})
	}
	return nil
}

func getDocument(uri protocol.DocumentUri) (*Document, bool) {
	if runtime.GOOS == "windows" {
		uri = strings.ToLower(uri)
	}
	doc, ok := documents.Load(uri)
	if !ok {
		return nil, false
	}
	return doc.(*Document), true
}
