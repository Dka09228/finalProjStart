package jsonlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo  Level = iota // Has the value of 0.
	LevelError              // Has the value of 1.
	LevelFatal              // Has the value of 2.
	LevelOff                // Has the value of 3.
)

// String returns a human-friendly string for the severity level.
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func NewLogger(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Logger) PrintPanic(err interface{}, stackTrace []byte) {
	properties := map[string]string{
		"stackTrace": string(stackTrace),
	}
	l.print(LevelFatal, fmt.Sprintf("PANIC: %v", err), properties)
	os.Exit(1)
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {

	if level < l.minLevel {
		return 0, nil
	}

	logEntry := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	logLine, err := json.Marshal(logEntry)
	if err != nil {
		logLine = []byte(LevelError.String() + ": unable to marshal log message:" + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(logLine, '\n'))
}

func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
