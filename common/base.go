package common

import (
	"bytes"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func CheckErrOrFatal(err error) {
	if err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}
}

func FormatJson(b []byte) (string, error) {
	var str bytes.Buffer
	err := json.Indent(&str, b, "", "    ")
	return str.String(), err
}

func FormatJsonStr(src string) (string, error) {
	var str bytes.Buffer
	err := json.Indent(&str, []byte(src), "", "    ")
	return str.String(), err
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

func HostName() (string, error) {
	return os.Hostname()
}

func HostnameOrFatal() string {
	if name, err := os.Hostname(); err == nil {
		return name
	} else {
		log.Fatal().Err(err)
	}

	return ""
}
