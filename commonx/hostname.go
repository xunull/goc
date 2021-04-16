package commonx

import (
	"github.com/rs/zerolog/log"
	"os"
)

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
