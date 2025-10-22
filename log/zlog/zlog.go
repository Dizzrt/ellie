package zlog

import (
	"os"

	"github.com/dizzrt/filerotator"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(file string, opts ...Option) (*zap.Logger, error) {
	conf := defaultConfig()
	for _, opt := range opts {
		opt(conf)
	}

	core, err := buildCore(file, conf)
	if err != nil {
		return nil, err
	}

	zl := zap.New(core).WithOptions(conf.ZapOpts...)
	return zl, nil
}

func buildCore(file string, conf *config) (zapcore.Core, error) {
	cores := make([]zapcore.Core, 0, 1)

	if conf.OutputType == OutputType_File || conf.OutputType == OutputType_Both {
		rotatorOpts := conf.toFileRotatorOptions()
		rotator, err := filerotator.New(file, rotatorOpts...)
		if err != nil {
			return nil, err
		}

		if conf.AsyncWrite {
			writeSyncer := &zapcore.BufferedWriteSyncer{
				WS:            zapcore.AddSync(rotator),
				Size:          conf.BufferSize,
				FlushInterval: conf.FlushInterval,
			}

			cores = append(cores, zapcore.NewCore(defaultEncoder(), writeSyncer, conf.Level))
		} else {
			cores = append(cores, zapcore.NewCore(defaultEncoder(), zapcore.AddSync(rotator), conf.Level))
		}
	}

	if conf.OutputType == OutputType_Console || conf.OutputType == OutputType_Both {
		cores = append(cores, zapcore.NewCore(defaultEncoder(), zapcore.Lock(os.Stdout), conf.Level))
	}

	tee := zapcore.NewTee(cores...)
	return tee, nil
}
