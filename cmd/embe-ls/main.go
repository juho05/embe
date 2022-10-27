package main

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"github.com/tliron/kutil/logging"
	_ "github.com/tliron/kutil/logging/simple"
)

var (
	name           = "embe-ls"
	version string = "dev"
)

var handler protocol.Handler

func main() {
	loadConfig()
	initLog()
	Info("Starting %s %s...", name, version)
	glspLogLevel := 0
	if ConfGLSPLogFile != nil {
		glspLogLevel = 2
	}
	logging.Configure(glspLogLevel, ConfGLSPLogFile)

	handler = protocol.Handler{
		Initialize:                    initialize,
		Initialized:                   initialized,
		Shutdown:                      shutdown,
		SetTrace:                      setTrace,
		TextDocumentDidOpen:           textDocumentDidOpen,
		TextDocumentDidChange:         textDocumentDidChange,
		TextDocumentDidClose:          textDocumentDidClose,
		TextDocumentCompletion:        textDocumentCompletion,
		TextDocumentSignatureHelp:     textDocumentSignatureHelp,
		TextDocumentHover:             textDocumentHover,
		TextDocumentColor:             textDocumentColor,
		TextDocumentColorPresentation: textDocumentColorPresentation,
		TextDocumentDefinition:        textDocumentDefinition,
	}

	var protocol string
	pflag.StringVarP(&protocol, "protocol", "p", "stdio", "The protocol to use. ('stdio', 'tcp', 'websocket', 'node-ipc')")
	var address string
	pflag.StringVarP(&address, "address", "a", ":4389", "The address to use for a TCP or WebSocket protocol.")
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.Parse()

	server := server.NewServer(&handler, name, ConfGLSPLogFile != nil)

	var err error
	switch protocol {
	case "stdio":
		Info("Protocol: STDIO")
		err = server.RunStdio()
	case "tcp":
		Info("Protocol: TCP")
		err = server.RunTCP(address)
	case "websocket":
		Info("Protocol: WebSocket")
		err = server.RunWebSocket(address)
	case "node-ipc":
		Info("Protocol: Node IPC")
		err = server.RunNodeJs()
	default:
		err = fmt.Errorf("Unsupported protocol: %s", protocol)
	}
	if err != nil {
		Fatal(err.Error())
	}
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	Trace("Initializing capabilities...")
	capabilities := handler.CreateServerCapabilities()
	capabilities.TextDocumentSync = protocol.TextDocumentSyncKindIncremental
	capabilities.CompletionProvider = &protocol.CompletionOptions{
		TriggerCharacters: []string{"@", "#", "."},
	}
	capabilities.SignatureHelpProvider = &protocol.SignatureHelpOptions{
		TriggerCharacters: []string{"(", ","},
	}
	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    name,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	Trace("Initialized.")
	return nil
}

func shutdown(context *glsp.Context) error {
	Info("Shutdown.")
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
