package commonx

import "github.com/rs/zerolog/log"

func Must(err error) {
	if err != nil {
		log.Error().Err(err).Msg("")
		panic(err)
	}
}
