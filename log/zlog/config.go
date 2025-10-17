package zlog

import (
	"time"

	"github.com/Dizzrt/filerotator"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type outputType uint8

const (
	OutputType_None outputType = iota
	OutputType_Both
	OutputType_File
	OutputType_Console
)

type config struct {
	Level   zapcore.Level
	Symlink string

	Clock        filerotator.Clock
	RotateType   filerotator.RotateType
	MaxAge       time.Duration
	MaxBackups   uint
	RotationTime time.Duration
	RotationSize int64

	OutputType    outputType
	AsyncWrite    bool
	BufferSize    int
	FlushInterval time.Duration

	ZapOpts []zap.Option
}

func defaultConfig() *config {
	return &config{
		Level:   zapcore.InfoLevel,
		Symlink: "",

		Clock:        filerotator.Local,          // default to local time
		RotateType:   filerotator.RotateTypeTime, // default to roatate by time
		MaxAge:       0,                          // default to unlimited
		MaxBackups:   0,                          // default to unlimited
		RotationTime: 1 * time.Hour,              // default to rotate every hour
		RotationSize: 10 * 1024 * 1024,           // default to rotate every 10MB

		OutputType: OutputType_File,
		AsyncWrite: true,
		BufferSize: 8 * 1024, // default to 8KB
	}
}

func (conf *config) toFileRotatorOptions() []filerotator.Option {
	opts := []filerotator.Option{
		filerotator.WithClock(conf.Clock),
		filerotator.WithRotateType(conf.RotateType),
		filerotator.WithRotationTime(conf.RotationTime),
		filerotator.WithRotationSize(conf.RotationSize),
		filerotator.WithMaxAge(conf.MaxAge),
		filerotator.WithMaxBackup(conf.MaxBackups),
		filerotator.WithSymlink(conf.Symlink),
	}

	return opts
}

func ParseOutputType(outputType string) outputType {
	switch outputType {
	case "none":
		return OutputType_None
	case "both":
		return OutputType_Both
	case "file":
		return OutputType_File
	case "console":
		return OutputType_Console
	default:
		return OutputType_File
	}
}

func ParseRotationType(rotateType string) filerotator.RotateType {
	return filerotator.ParseRotationType(rotateType)
}
