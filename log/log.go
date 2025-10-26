package log

import (
	"context"
	"fmt"
)

const DefaultMessageKey = "msg"

type LogWriter interface {
	Write(level Level, keyvals ...any) error
}

type LogAsyncWriter interface {
	Sync() error
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

func (logger *Logger) isValidLevel(level Level) bool {
	// TODO
	return true
}

func (logger *Logger) Writer() LogWriter {
	return logger.writer
}

func (logger *Logger) Write(level Level, keyvals ...any) {
	_ = logger.writer.Write(level, keyvals...)
}

func (logger *Logger) Debug(a ...any) {
	if !logger.isValidLevel(LevelDebug) {
		return
	}

	_ = logger.writer.Write(LevelDebug, logger.msgKey, logger.sprint(a...))
}

func (logger *Logger) Debugf(format string, a ...any) {
	if !logger.isValidLevel(LevelDebug) {
		return
	}

	_ = logger.writer.Write(LevelDebug, logger.msgKey, logger.sprintf(format, a...))
}

func (logger *Logger) Debugw(keyvals ...any) {
	if !logger.isValidLevel(LevelDebug) {
		return
	}

	_ = logger.writer.Write(LevelDebug, keyvals...)
}

func (logger *Logger) Info(a ...any) {
	if !logger.isValidLevel(LevelInfo) {
		return
	}

	_ = logger.writer.Write(LevelInfo, logger.msgKey, logger.sprint(a...))
}

func (logger *Logger) Infof(format string, a ...any) {
	if !logger.isValidLevel(LevelInfo) {
		return
	}

	_ = logger.writer.Write(LevelInfo, logger.msgKey, logger.sprintf(format, a...))
}

func (logger *Logger) Infow(keyvals ...any) {
	if !logger.isValidLevel(LevelInfo) {
		return
	}

	_ = logger.writer.Write(LevelInfo, keyvals...)
}

func (logger *Logger) Warn(a ...any) {
	if !logger.isValidLevel(LevelWarn) {
		return
	}

	_ = logger.writer.Write(LevelWarn, logger.msgKey, logger.sprint(a...))
}

func (logger *Logger) Warnf(format string, a ...any) {
	if !logger.isValidLevel(LevelWarn) {
		return
	}

	_ = logger.writer.Write(LevelWarn, logger.msgKey, logger.sprintf(format, a...))
}

func (logger *Logger) Warnw(keyvals ...any) {
	if !logger.isValidLevel(LevelWarn) {
		return
	}

	_ = logger.writer.Write(LevelWarn, keyvals...)
}

func (logger *Logger) Error(a ...any) {
	if !logger.isValidLevel(LevelError) {
		return
	}

	_ = logger.writer.Write(LevelError, logger.msgKey, logger.sprint(a...))
}

func (logger *Logger) Errorf(format string, a ...any) {
	if !logger.isValidLevel(LevelError) {
		return
	}

	_ = logger.writer.Write(LevelError, logger.msgKey, logger.sprintf(format, a...))
}

func (logger *Logger) Errorw(keyvals ...any) {
	if !logger.isValidLevel(LevelError) {
		return
	}

	_ = logger.writer.Write(LevelError, keyvals...)
}

func (logger *Logger) DPanic(a ...any) {
	if !logger.isValidLevel(LevelDPanic) {
		return
	}

	_ = logger.writer.Write(LevelDPanic, logger.msgKey, logger.sprint(a...))
}

func (logger *Logger) DPanicf(format string, a ...any) {
	if !logger.isValidLevel(LevelDPanic) {
		return
	}

	_ = logger.writer.Write(LevelDPanic, logger.msgKey, logger.sprintf(format, a...))
}

func (logger *Logger) DPanicw(keyvals ...any) {
	if !logger.isValidLevel(LevelDPanic) {
		return
	}

	_ = logger.writer.Write(LevelDPanic, keyvals...)
}

func (logger *Logger) Panic(a ...any) {
	if !logger.isValidLevel(LevelPanic) {
		return
	}

	_ = logger.writer.Write(LevelPanic, logger.msgKey, logger.sprint(a...))
}

func (logger *Logger) Panicf(format string, a ...any) {
	if !logger.isValidLevel(LevelPanic) {
		return
	}

	_ = logger.writer.Write(LevelPanic, logger.msgKey, logger.sprintf(format, a...))
}

func (logger *Logger) Panicw(keyvals ...any) {
	if !logger.isValidLevel(LevelPanic) {
		return
	}

	_ = logger.writer.Write(LevelPanic, keyvals...)
}

func (logger *Logger) Fatal(a ...any) {
	if !logger.isValidLevel(LevelFatal) {
		return
	}

	_ = logger.writer.Write(LevelFatal, logger.msgKey, logger.sprint(a...))
}

func (logger *Logger) Fatalf(format string, a ...any) {
	if !logger.isValidLevel(LevelFatal) {
		return
	}

	_ = logger.writer.Write(LevelFatal, logger.msgKey, logger.sprintf(format, a...))
}

func (logger *Logger) Fatalw(keyvals ...any) {
	if !logger.isValidLevel(LevelFatal) {
		return
	}

	_ = logger.writer.Write(LevelFatal, keyvals...)
}
