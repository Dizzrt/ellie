package log

import (
	"context"
	"fmt"
)

const DefaultMessageKey = "msg"

type LogWriter interface {
	Write(level Level, keyvals ...any) error
}

type Logger struct {
	writer  LogWriter
	msgKey  string
	sprint  func(...any) string
	sprintf func(format string, a ...any) string
}

func NewLogger(writer LogWriter, opts ...Option) *Logger {
	logger := &Logger{
		writer:  writer,
		msgKey:  DefaultMessageKey,
		sprint:  fmt.Sprint,
		sprintf: fmt.Sprintf,
	}

	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

func (logger *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		// logger:
		// TODO
	}
}
