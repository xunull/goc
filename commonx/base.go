package commonx

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

func CheckErrOrFatal(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Println(debug.Stack())
		//log.Fatal().Err(err).Msgf("%v\n%s", err, string(debug.Stack()))

	}
}

func QuitWatch() {
	c := make(chan os.Signal)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-c
	os.Exit(-1)
}

func QuitListen() {
	c := make(chan os.Signal)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-c
}
