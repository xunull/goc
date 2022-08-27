package logx

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"time"
)

func InitZeroFileLog(level zerolog.Level, name, p string) *zerolog.Logger {
	hook := lumberjack.Logger{
		Filename:   path.Join(p, logFileName(name)),
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     3,
		Compress:   true,
	}
	logger := log.Output(zerolog.ConsoleWriter{Out: &hook, TimeFormat: time.RFC3339})
	logger.Level(level)
	return &logger
}

func InitZeroConsoleLog(level zerolog.Level) *zerolog.Logger {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	logger.Level(level)
	return &logger
}
