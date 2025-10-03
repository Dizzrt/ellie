package zlog

import (
	"time"

	"github.com/Dizzrt/filerotator"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(*config)

func Level(level zapcore.Level) Option {
	return func(conf *config) {
		conf.Level = level
	}
}

func Suffix(suffix string) Option {
	return func(conf *config) {
		conf.Suffix = suffix
	}
}

func LinkFile(linkFile string) Option {
	return func(conf *config) {
		conf.LinkFile = linkFile
	}
}

func Clock(clock filerotator.Clock) Option {
	return func(conf *config) {
		conf.Clock = clock
	}
}

func RotateType(rotateType filerotator.RotateType) Option {
	return func(conf *config) {
		conf.RotateType = rotateType
	}
}

func MaxAge(maxAge time.Duration) Option {
	return func(conf *config) {
		conf.MaxAge = maxAge
	}
}

func MaxBackups(maxBackups uint) Option {
	return func(conf *config) {
		conf.MaxBackups = maxBackups
	}
}

func RotationTime(rotationTime time.Duration) Option {
	return func(conf *config) {
		conf.RotationTime = rotationTime
	}
}

func RotationSize(rotationSize int64) Option {
	return func(conf *config) {
		conf.RotationSize = rotationSize
	}
}

func OutputType(outputType outputType) Option {
	return func(conf *config) {
		conf.OutputType = outputType
	}
}

func AsyncWrite(asyncWrite bool) Option {
	return func(conf *config) {
		conf.AsyncWrite = asyncWrite
	}
}

func BufferSize(bufferSize int) Option {
	return func(conf *config) {
		conf.BufferSize = bufferSize
	}
}

func FlushInterval(flushInterval time.Duration) Option {
	return func(conf *config) {
		conf.FlushInterval = flushInterval
	}
}

func ZapOpts(zapOpts ...zap.Option) Option {
	return func(conf *config) {
		conf.ZapOpts = append(conf.ZapOpts, zapOpts...)
	}
}
