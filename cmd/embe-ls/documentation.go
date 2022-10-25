package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

//go:embed documentation.md
var documentationFile []byte

var documentation = make(map[string]string)

func getDocs(identifier string) any {
	if docs, ok := documentation[identifier]; ok {
		return protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: docs,
		}
	}
	Warn("Docs not available for: %s", identifier)
	return nil
}

func init() {
	scanner := bufio.NewScanner(bytes.NewBuffer(documentationFile))

	var name string
	var sbuilder strings.Builder
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(text, "//") {
			continue
		} else if name == "" {
			name = text
		} else if text == "---" {
			documentation[name] = strings.TrimSpace(sbuilder.String())
			name = ""
			sbuilder.Reset()
		} else {
			sbuilder.WriteRune('\n')
			sbuilder.WriteString(text)
		}
	}
}
