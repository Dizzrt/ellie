package zlog

import (
	"time"

	"go.uber.org/zap/zapcore"
)

func defaultEncoder() zapcore.Encoder {
	config := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		FunctionKey:   zapcore.OmitKey,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime:    timeEncoder,
		EncodeLevel:   levelEncoder,
		EncodeCaller:  callerEncoder,
	}

	return zapcore.NewConsoleEncoder(config)
}

func levelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format(time.DateTime) + "]")
}

func callerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.FullPath() + "]")
}
