package config

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"
)

const logsFlag = log.Ldate | log.Ltime

const (
	prefixInfo    = "\033[32m[INFO]: "
	prefixDebug   = "\033[34m[DEBUG]: "
	prefixError   = "\033[31m[ERROR]: "
	prefixWarning = "\033[33m[WARNING]: "
)

// Logger log wrapper with 4 log levels: INFO, DEBUG, ERROR, WARNING
// at production mode write logs to files
type Logger struct {
	ErrOut   *log.Logger
	InfoOut  *log.Logger
	WarnOut  *log.Logger
	DebugOut *log.Logger

	outputs []*os.File
}

// Close all output files
func (l *Logger) Close() {
	for _, out := range l.outputs {
		out.Close()
	}
}

// NamedFmt named strings formatter
func (l *Logger) NamedFmt(msg string, args any) string {
	t, err := template.New("").Parse(msg)

	if err != nil {
		return msg
	}

	b := bytes.NewBufferString(msg)

	if err = t.Execute(b, args); err == nil {
		return b.String()
	}

	return msg
}

// Info log any messages
func (l *Logger) Printf(format string, args ...interface{}) {
	l.ErrOut.Printf(format, args...)
}

// Info log INFO messages
func (l *Logger) Info(val ...interface{}) {
	l.InfoOut.Println(val...)
}

// Infof log formated INFO messages
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warn log WARN messages
func (l *Logger) Warn(val ...interface{}) {
	l.WarnOut.Println(val...)
}

// Warnf log formated WARN messages
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// Debug log DEBUG messages
func (l *Logger) Debug(val ...interface{}) {
	l.DebugOut.Println(val...)
}

// Debugf log formated DEBUG messages
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

// Error log ERROR messages
func (l *Logger) Error(val ...interface{}) {
	l.ErrOut.Println(val...)
}

// Errorf log formated ERROR messages
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// Fatalln log ERROR messages and call os.Exit(1)
func (l *Logger) Fatalln(args ...interface{}) {
	l.Error(args...)
	os.Exit(1)
}

// Fatalf log formated ERROR messages and call os.Exit(1)
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Errorf(format, args...)
	os.Exit(1)
}

// NewLogger create new Logger instance
// In debug mode write DEBUG logs into stdout, otherwise into dev/null
// In prod mode write INFO, WARNNING, ERROR into output files
func NewLogger(debug bool, appMode string, logsPath string) *Logger {
	l := new(Logger)

	l.WarnOut = log.New(os.Stdout, prefixWarning, logsFlag)
	l.InfoOut = log.New(os.Stdout, prefixInfo, logsFlag)
	l.ErrOut = log.New(os.Stderr, prefixError, logsFlag)

	if debug {
		l.DebugOut = log.New(os.Stdout, prefixDebug, logsFlag)
	} else {
		l.DebugOut = log.New(io.Discard, prefixDebug, logsFlag)
	}

	return l
}
