// pkg/logger/interface.go
package logger

import (
	"fmt"
	"time"
)

// LogLevel controls the verbosity of the logger.
type LogLevel uint8

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelNone // Disables all logging
)

// ANSI Color Codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
)

// Logger defines the contract for logging within the framework.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	SetLevel(level LogLevel)
}

// globalLogger holds the current logger instance.
var globalLogger Logger = &StandardLogger{level: LevelInfo}

// SetLogger replaces the global logger instance.
func SetLogger(l Logger) {
	globalLogger = l
}

// SetLevel adjusts the verbosity of the global logger.
func SetLevel(level LogLevel) {
	globalLogger.SetLevel(level)
}

// Global accessor functions
func Debug(msg string, args ...any) { globalLogger.Debug(msg, args...) }
func Info(msg string, args ...any)  { globalLogger.Info(msg, args...) }
func Warn(msg string, args ...any)  { globalLogger.Warn(msg, args...) }
func Error(msg string, args ...any) { globalLogger.Error(msg, args...) }

// Tag wraps a module name in brackets and colors it Purple.
// Note: usage involves string allocation. Use sparingly in tight loops.
func Tag(name string) string {
	return fmt.Sprintf("%s[%s]%s", Cyan, name, Reset)
}

// StandardLogger implementation using standard fmt.Printf.
type StandardLogger struct {
	level LogLevel
}

func (l *StandardLogger) SetLevel(level LogLevel) {
	l.level = level
}

// print handles the formatting and output.
// We avoid time.Format to save binary size (heavy dependency).
func (l *StandardLogger) print(reqLevel LogLevel, color, label, msg string, args ...any) {
	if reqLevel < l.level {
		return
	}

	// Calculate timestamp manually to avoid 'time' package formatting overhead
	now := time.Now()
	// UnixNano is supported efficiently on most TinyGo targets
	nanos := now.UnixNano()
	
	sec := nanos / 1e9
	ms := (nanos % 1e9) / 1e6

	// Format: SSS.mmm [LEVEL] MESSAGE
	// We construct the prefix manually and let fmt handle the user args.
	prefix := fmt.Sprintf("%d.%03d %s[%s]%s ", sec, ms, color, label, Reset)
	
	// Print in one go. Using \r\n for serial terminal compatibility.
	fmt.Printf(prefix+msg+"\r\n", args...)
}

func (l *StandardLogger) Debug(msg string, args ...any) {
	l.print(LevelDebug, Green, "DEBUG", msg, args...)
}

func (l *StandardLogger) Info(msg string, args ...any) {
	l.print(LevelInfo, Blue, "INFO ", msg, args...)
}

func (l *StandardLogger) Warn(msg string, args ...any) {
	l.print(LevelWarn, Yellow, "WARN ", msg, args...)
}

func (l *StandardLogger) Error(msg string, args ...any) {
	l.print(LevelError, Red, "ERROR", msg, args...)
}