package logx

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
)

func logFileName(name string) string {
	return fmt.Sprintf("%s.log", name)
}

func logErrorFileName(name string) string {
	return fmt.Sprintf("%s.error.log", name)
}

func newFileCore(level zapcore.Level, dir, name string) zapcore.Core {
	// todo this field use option type
	hook := lumberjack.Logger{
		Filename:   path.Join(dir, logFileName(name)),
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     7,
		Compress:   true,
	}
	w := zapcore.AddSync(&hook)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		level,
	)
}

func newErrorFileCore(level zapcore.Level, dir, name string) zapcore.Core {
	// todo this field use option type
	hook := lumberjack.Logger{
		Filename:   path.Join(dir, logErrorFileName(name)),
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     7,
		Compress:   true,
	}
	w := zapcore.AddSync(&hook)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		zapcore.ErrorLevel,
	)
}

func InitZapFileLogger(level, dir, name string) (*zap.Logger, error) {
	var l zapcore.Level
	err := l.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}

	var allCore []zapcore.Core

	allCore = append(allCore,
		newFileCore(l, dir, name),
		newErrorFileCore(l, dir, name),
	)

	core := zapcore.NewTee(allCore...)
	logger := zap.New(core).WithOptions(zap.AddCaller())

	return logger, nil
}
