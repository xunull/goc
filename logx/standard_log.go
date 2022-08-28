package logx

import (
	"fmt"
	"github.com/xunull/goc/commonx"
	"log"
	"os"
	"path"
	"time"
)

func InitStandardFileLogger(name, p string) *log.Logger {
	file, err := os.OpenFile(path.Join(p, logFileName(name)), os.O_RDWR|os.O_CREATE, 0666)
	fmt.Println(file)
	time.Sleep(1 * time.Second)
	commonx.CheckErrOrFatal(err)
	logger := log.Default()
	logger.SetOutput(file)
	return logger
}
