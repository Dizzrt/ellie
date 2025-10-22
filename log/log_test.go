package log

import (
	"testing"

	"github.com/dizzrt/ellie/log/zlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestStdLoggerWriter(t *testing.T) {
	writer, err := NewStdLoggerWriter("logs/test.log",
		zlog.ZapOpts(
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		),
	)

	if err != nil {
		t.Fatal(err)
	}

	writer.Write(LevelDebug, "msg", "debug message")
	writer.Write(LevelInfo, "msg", "info message")
	writer.Write(LevelWarn, "msg", "warn message")
	writer.Write(LevelError, "msg", "error message")

	writer.(*stdLoggerWriter).Sync()
}

func TestLogger(t *testing.T) {
	writer, err := NewStdLoggerWriter("logs/log",
		zlog.OutputType(zlog.OutputType_Both),
		zlog.Level(zapcore.DebugLevel),
		zlog.ZapOpts(
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
			zap.AddCallerSkip(2),
		),
	)

	if err != nil {
		t.Fatal(err)
	}

	logger := NewLogger(writer)

	logger.Debug("debug message")
	logger.Debugf("debugf message: %d", 123)
	logger.Debugw("msgx", "xxx", "key1", "value1", "key2", 2)
	logger.Infow("msg", "infow message", "key1", "value1", "key2", 2, "key1", "value11")
	logger.Info("info message")
	logger.Warn("warn message")
	// logger.Error("error message")
	// logger.DPanic("dpanic message")
	// logger.Panic("panic message")
	// logger.Fatal("fatal message")

	logger.writer.(*stdLoggerWriter).Sync()
}
