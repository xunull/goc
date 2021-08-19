package commonx

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

func CheckErrOrFatal(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("\n" + string(debug.Stack()))
		os.Exit(1)
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
}
