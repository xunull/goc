package color_out

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func Fatal(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	color.Red("[Fatal]:%s", s)
	os.Exit(1)
}

func Error(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	color.Red("[Error]:%s", s)
}

func Info(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	color.Green("[Info]:%s", s)
}

func Warn(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	color.White("[Warn]:%s", s)
}

func Debug(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	color.Yellow("[Debug]:%s", s)
}
