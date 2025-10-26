package log

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/dizzrt/ellie/log/zlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var global = &loggerAppliance{}

type loggerAppliance struct {
	LogWriter
	lock sync.RWMutex
}

func init() {
	writer, err := NewStdLoggerWriter("",
		zlog.OutputType(zlog.OutputType_Console),
		zlog.Level(zapcore.DebugLevel),
		zlog.ZapOpts(
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
			zap.AddCallerSkip(2),
		),
	)

	if err != nil {
		panic(err)
	}

	global.SetLogger(writer)
}

func (a *loggerAppliance) SetLogger(writer LogWriter) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.LogWriter = writer
}

func SetLogger(writer LogWriter) {
	global.SetLogger(writer)
}

func GetLogger() LogWriter {
	global.lock.RLock()
	defer global.lock.RUnlock()
	return global.LogWriter
}

func Sync() {
	writer := global.LogWriter
	asyncWriter, ok := writer.(LogAsyncWriter)
	if !ok {
		return
	}

	asyncWriter.Sync()
}

func Debug(a ...any) {
	global.Write(LevelDebug, DefaultMessageKey, fmt.Sprint(a...))
}

func Debugf(format string, a ...any) {
	global.Write(LevelDebug, DefaultMessageKey, fmt.Sprintf(format, a...))
}

func Debugw(keyvals ...any) {
	global.Write(LevelDebug, keyvals...)
}

func Info(a ...any) {
	global.Write(LevelInfo, DefaultMessageKey, fmt.Sprint(a...))
}

func Infof(format string, a ...any) {
	global.Write(LevelInfo, DefaultMessageKey, fmt.Sprintf(format, a...))
}

func Infow(keyvals ...any) {
	global.Write(LevelInfo, keyvals...)
}

func Warn(a ...any) {
	global.Write(LevelWarn, DefaultMessageKey, fmt.Sprint(a...))
}

func Warnf(format string, a ...any) {
	global.Write(LevelWarn, DefaultMessageKey, fmt.Sprintf(format, a...))
}

func Warnw(keyvals ...any) {
	global.Write(LevelWarn, keyvals...)
}

func Error(a ...any) {
	global.Write(LevelError, DefaultMessageKey, fmt.Sprint(a...))
}

func Errorf(format string, a ...any) {
	global.Write(LevelError, DefaultMessageKey, fmt.Sprintf(format, a...))
}

func Errorw(keyvals ...any) {
	global.Write(LevelError, keyvals...)
}

func Fatal(a ...any) {
	global.Write(LevelFatal, DefaultMessageKey, fmt.Sprint(a...))
	os.Exit(1)
}

func Fatalf(format string, a ...any) {
	global.Write(LevelFatal, DefaultMessageKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

func Fatalw(keyvals ...any) {
	global.Write(LevelFatal, keyvals...)
	os.Exit(1)
}

// log with context

func fromCtx(ctx context.Context, kvs ...any) []any {
	if ctx == nil {
		return kvs
	}

	for _, k := range fromCtxKeys {
		if v := ctx.Value(k); v != nil {
			kvs = append(kvs, k, v)
		}
	}

	return kvs
}

func CtxDebug(ctx context.Context, a ...any) {
	kvs := fromCtx(ctx, DefaultMessageKey, fmt.Sprint(a...))
	global.Write(LevelDebug, kvs...)
}

func CtxDebugf(ctx context.Context, format string, a ...any) {
	kvs := fromCtx(ctx, DefaultMessageKey, fmt.Sprintf(format, a...))
	global.Write(LevelDebug, kvs...)
}

func CtxDebugw(ctx context.Context, keyvals ...any) {
	kvs := fromCtx(ctx, keyvals...)
	global.Write(LevelDebug, kvs...)
}

func CtxInfo(ctx context.Context, a ...any) {
	kvs := fromCtx(ctx, DefaultMessageKey, fmt.Sprint(a...))
	global.Write(LevelInfo, kvs...)
}

func CtxInfof(ctx context.Context, format string, a ...any) {
	kvs := fromCtx(ctx, DefaultMessageKey, fmt.Sprintf(format, a...))
	global.Write(LevelInfo, kvs...)
}

func CtxInfow(ctx context.Context, keyvals ...any) {
	kvs := fromCtx(ctx, keyvals...)
	global.Write(LevelInfo, kvs...)
}

func CtxWarn(ctx context.Context, a ...any) {
	kvs := fromCtx(ctx, DefaultMessageKey, fmt.Sprint(a...))
	global.Write(LevelWarn, kvs...)
}

func CtxWarnf(ctx context.Context, format string, a ...any) {
	kvs := fromCtx(ctx, DefaultMessageKey, fmt.Sprintf(format, a...))
	global.Write(LevelWarn, kvs...)
}

func CtxWarnw(ctx context.Context, keyvals ...any) {
	kvs := fromCtx(ctx, keyvals...)
	global.Write(LevelWarn, kvs...)
}

func CtxError(ctx context.Context, a ...any) {
	kvs := fromCtx(ctx, DefaultMessageKey, fmt.Sprint(a...))
	global.Write(LevelError, kvs...)
}

func CtxErrorf(ctx context.Context, format string, a ...any) {
	kvs := fromCtx(ctx, DefaultMessageKey, fmt.Sprintf(format, a...))
	global.Write(LevelError, kvs...)
}

func CtxErrorw(ctx context.Context, keyvals ...any) {
	kvs := fromCtx(ctx, keyvals...)
	global.Write(LevelError, kvs...)
}

func CtxFatal(ctx context.Context, a ...any) {
	kvs := fromCtx(ctx, DefaultMessageKey, fmt.Sprint(a...))
	global.Write(LevelFatal, kvs...)
	os.Exit(1)
}

func CtxFatalf(ctx context.Context, format string, a ...any) {
	kvs := fromCtx(ctx, DefaultMessageKey, fmt.Sprintf(format, a...))
	global.Write(LevelFatal, kvs...)
	os.Exit(1)
}

func CtxFatalw(ctx context.Context, keyvals ...any) {
	kvs := fromCtx(ctx, keyvals...)
	global.Write(LevelFatal, kvs...)
	os.Exit(1)
}
