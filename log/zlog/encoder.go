package zlog

import (
	"go.uber.org/zap/zapcore"
)

func defaultEncoder() zapcore.Encoder {
	config := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		FunctionKey:    zapcore.OmitKey,
	}

	return zapcore.NewJSONEncoder(config)
}
