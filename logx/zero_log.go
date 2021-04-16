package logx

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"strings"
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

func Focus(format string, v ...interface{}) {

	f := "\n" + strings.Repeat("+", 25) + "\n" +
		fmt.Sprintf(format, v...) +
		"\n" + strings.Repeat("-", 25) + "\n"
	// todo \n
	//log.Debug().Msg(f)
	//log.Print(f)
	fmt.Println(f)

}

func FocusErr(err error, format string, v ...interface{}) {
	f := "\n" + strings.Repeat("+", 25) + "\n" +
		format +
		"\n" + strings.Repeat("-", 25) + "\n"
	log.Debug().Err(err).Msgf(f, v...)
}
