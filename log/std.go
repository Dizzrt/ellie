package log

import (
	"fmt"

	"github.com/Dizzrt/ellie/log/zlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ LogWriter = (*stdLoggerWriter)(nil)

type stdLoggerWriter struct {
	zapLogger *zap.Logger
	sugared   *zap.SugaredLogger
}

func NewStdLoggerWriterWithCustomZap(zapLogger *zap.Logger) (LogWriter, error) {
	return &stdLoggerWriter{
		zapLogger: zapLogger,
		sugared:   zapLogger.Sugar(),
	}, nil
}

func NewStdLoggerWriter(file string, opts ...zlog.Option) (LogWriter, error) {
	zapLogger, err := zlog.New(file, opts...)
	if err != nil {
		return nil, err
	}

	return &stdLoggerWriter{
		zapLogger: zapLogger,
		sugared:   zapLogger.Sugar(),
	}, nil
}

func (logger *stdLoggerWriter) Write(level Level, keyvals ...any) error {
	zlevel := zapcore.Level(level)

	keyLen := len(keyvals)
	if keyLen == 0 || keyLen&1 == 1 {
		logger.zapLogger.Warn(fmt.Sprint("keyvals must appear in pairs: ", keyvals))
		return nil
	}

	hasMsg := false
	for i := 0; i < keyLen; i += 2 {
		if keyvals[i] == "msg" {
			hasMsg = true
			break
		}
	}

	if hasMsg {
		msg := ""
		data := make([]zap.Field, 0, (keyLen>>1)+1)
		for i := 0; i < keyLen; i += 2 {
			if keyvals[i].(string) == "msg" {
				msg, _ = keyvals[i+1].(string)
				continue
			}

			data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
		}

		switch zlevel {
		case zapcore.DebugLevel:
			logger.zapLogger.Debug(msg, data...)
		case zapcore.InfoLevel:
			logger.zapLogger.Info(msg, data...)
		case zapcore.WarnLevel:
			logger.zapLogger.Warn(msg, data...)
		case zapcore.ErrorLevel:
			logger.zapLogger.Error(msg, data...)
		case zapcore.DPanicLevel:
			logger.zapLogger.DPanic(msg, data...)
		case zapcore.PanicLevel:
			logger.zapLogger.Panic(msg, data...)
		case zapcore.FatalLevel:
			logger.zapLogger.Fatal(msg, data...)
		}
	} else {
		s := ""
		for i := 0; i < keyLen; i += 2 {
			v := keyvals[i+1]
			if _, ok := v.(string); ok {
				v = fmt.Sprintf("\"%s\"", v)
			}

			s += fmt.Sprintf("\"%s\": %v", keyvals[i], v)
			if i != keyLen-2 {
				s += ", "
			}
		}
		s = "{" + s + "}"

		switch zlevel {
		case zapcore.DebugLevel:
			logger.sugared.Debug(s)
		case zapcore.InfoLevel:
			logger.sugared.Info(s)
		case zapcore.WarnLevel:
			logger.sugared.Warn(s)
		case zapcore.ErrorLevel:
			logger.sugared.Error(s)
		case zapcore.DPanicLevel:
			logger.sugared.DPanic(s)
		case zapcore.PanicLevel:
			logger.sugared.Panic(s)
		case zapcore.FatalLevel:
			logger.sugared.Fatal(s)
		}
	}

	return nil
}

func (logger *stdLoggerWriter) Sync() error {
	return logger.zapLogger.Sync()
}

func (logger *stdLoggerWriter) Close() error {
	return logger.Sync()
}
