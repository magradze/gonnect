package logger

import (
	"fmt"
	"time"
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
	Gray   = "\033[90m"
)

// Logger defines the contract for logging within the framework.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// globalLogger holds the current logger instance.
var globalLogger Logger = &StandardLogger{}

// SetLogger replaces the global logger instance.
func SetLogger(l Logger) {
	globalLogger = l
}

// Global accessor functions
func Debug(msg string, args ...any) { globalLogger.Debug(msg, args...) }
func Info(msg string, args ...any)  { globalLogger.Info(msg, args...) }
func Warn(msg string, args ...any)  { globalLogger.Warn(msg, args...) }
func Error(msg string, args ...any) { globalLogger.Error(msg, args...) }

// Tag wraps a module name in brackets and colors it Purple.
// Usage: logger.Info("%s Started", logger.Tag("MyModule"))
func Tag(name string) string {
	return fmt.Sprintf("%s[%s]%s", Purple, name, Reset)
}

// StandardLogger implementation using fmt.Printf with colors.
type StandardLogger struct{}

// getTimestamp returns just the time part (HH:MM:SS.000).
// Since MCUs reset to 1970, the date is irrelevant/distracting.
func (l *StandardLogger) getTimestamp() string {
	return time.Now().Format("15:04:05.000")
}

func (l *StandardLogger) print(level, color, msg string, args ...any) {
	timestamp := l.getTimestamp()
	formattedMsg := fmt.Sprintf(msg, args...)
	
	// Format: TIME [LEVEL] MESSAGE
	// We use \r\n for better compatibility with serial monitors
	fmt.Printf("%s %s[%s]%s %s\r\n", 
		timestamp, 
		color, level, Reset, 
		formattedMsg,
	)
}

func (l *StandardLogger) Debug(msg string, args ...any) {
	l.print("DEBUG", Green, msg, args...)
}

func (l *StandardLogger) Info(msg string, args ...any) {
	l.print("INFO ", Blue, msg, args...)
}

func (l *StandardLogger) Warn(msg string, args ...any) {
	l.print("WARN ", Yellow, msg, args...)
}

func (l *StandardLogger) Error(msg string, args ...any) {
	l.print("ERROR", Red, msg, args...)
}