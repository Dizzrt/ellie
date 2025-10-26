package log

import (
	"context"
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

func TestGlobalCtxLog(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, LogID, "123456789abc")
	ctx = context.WithValue(ctx, TraceID, "trace123456789")
	ctx = context.WithValue(ctx, SpanID, "span123456789")

	CtxDebug(ctx, "debug message")
	CtxDebugf(ctx, "debugf message: %d", 123)
	CtxDebugw(ctx, "msgx", "xxx", "key1", "value1", "key2", 2)
	CtxInfow(ctx, "msg", "infow message", "key1", "value1", "key2", 2, "key1", "value11")
	CtxInfo(ctx, "info message")
	CtxWarn(ctx, "warn message")
	CtxError(ctx, "error message")
	CtxErrorf(ctx, "errorf message: %d", 123)
	CtxErrorw(ctx, "msgx", "xxx", "key1", "value1", "key2", 2)
	// CtxFatal(ctx, "fatal message")
	// CtxFatalf(ctx, "fatalf message: %d", 123)

	Sync()
}
