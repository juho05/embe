package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

var (
	ConfLogFile     string
	ConfLogLevel    string
	ConfGLSPLogFile *string
)

type Config struct {
	LogFile     string  `json:"log_file"`
	LogLevel    string  `json:"log_level"`
	GLSPLogFile *string `json:"lsp_log_file"`
}

func loadConfig() {
	path := filepath.Join(xdg.ConfigHome, "embe-ls", "config.json")
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode config file: %s", err)
		return
	}

	ConfLogFile = config.LogFile
	ConfLogLevel = config.LogLevel
	ConfGLSPLogFile = config.GLSPLogFile
	if ConfGLSPLogFile != nil && *ConfGLSPLogFile == "" {
		ConfGLSPLogFile = nil
	}

	if ConfGLSPLogFile != nil {
		os.Remove(*ConfGLSPLogFile)
	}
}
