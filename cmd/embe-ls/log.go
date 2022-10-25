package main

import (
	"fmt"
	"os"
)

type LogLevel int

const (
	LogLevelTrace LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelFatal
	LogLevelNone
)

var (
	logFile  *os.File
	logLevel LogLevel
)

func initLog() {
	switch ConfLogLevel {
	case "":
		logLevel = LogLevelWarning
	case "trace":
		logLevel = LogLevelTrace
	case "info":
		logLevel = LogLevelInfo
	case "warning":
		logLevel = LogLevelWarning
	case "error":
		logLevel = LogLevelError
	case "fatal":
		logLevel = LogLevelFatal
	case "none":
		logLevel = LogLevelNone
	default:
		fmt.Fprintln(os.Stderr, "Invalid log level:", ConfLogLevel)
		logLevel = LogLevelWarning
	}

	if ConfLogFile == "" {
		logFile = os.Stderr
		return
	}

	var err error
	logFile, err = os.Create(ConfLogFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open log file:", err)
		logFile = os.Stderr
	}
}

func Trace(format string, a ...any) {
	if logLevel > LogLevelTrace {
		return
	}
	fmt.Fprintf(logFile, fmt.Sprintf("[TRACE] %s\n", format), a...)
}

func Info(format string, a ...any) {
	if logLevel > LogLevelInfo {
		return
	}
	fmt.Fprintf(logFile, fmt.Sprintf("[INFO]  %s\n", format), a...)
}

func Warn(format string, a ...any) {
	if logLevel > LogLevelWarning {
		return
	}
	fmt.Fprintf(logFile, fmt.Sprintf("[WARN]  %s\n", format), a...)
}

func Error(format string, a ...any) {
	if logLevel > LogLevelError {
		return
	}
	fmt.Fprintf(logFile, fmt.Sprintf("[ERROR] %s\n", format), a...)
}

func Fatal(format string, a ...any) {
	if logLevel > LogLevelFatal {
		return
	}
	fmt.Fprintf(logFile, fmt.Sprintf("[FATAL] %s\n", format), a...)
	logFile.Close()
	os.Exit(1)
}
