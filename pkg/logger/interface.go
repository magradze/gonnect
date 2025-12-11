package logger

import (
	"log"
	"os"
)

// Logger defines the contract for logging within the framework.
// Implementations can be swapped to support UART, USB Serial, or File logging.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// globalLogger holds the current logger instance.
// Defaults to a standard library logger until replaced.
var globalLogger Logger = &StandardLogger{
	std: log.New(os.Stdout, "", log.LstdFlags),
}

// SetLogger replaces the global logger instance.
// This should be called early in the main application setup.
func SetLogger(l Logger) {
	globalLogger = l
}

// Global accessor functions for convenience.

func Debug(msg string, args ...any) { globalLogger.Debug(msg, args...) }
func Info(msg string, args ...any)  { globalLogger.Info(msg, args...) }
func Warn(msg string, args ...any)  { globalLogger.Warn(msg, args...) }
func Error(msg string, args ...any) { globalLogger.Error(msg, args...) }

// StandardLogger is a basic implementation using Go's standard log package.
// It serves as the default logger for the framework.
type StandardLogger struct {
	std *log.Logger
}

func (l *StandardLogger) Debug(msg string, args ...any) {
	l.std.Printf("[DEBUG] "+msg, args...)
}

func (l *StandardLogger) Info(msg string, args ...any) {
	l.std.Printf("[INFO]  "+msg, args...)
}

func (l *StandardLogger) Warn(msg string, args ...any) {
	l.std.Printf("[WARN]  "+msg, args...)
}

func (l *StandardLogger) Error(msg string, args ...any) {
	l.std.Printf("[ERROR] "+msg, args...)
}