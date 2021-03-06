package file_path

import (
	"github.com/rs/zerolog/log"
	"os"
)

func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// not found
			return false
		}
		log.Info().Msgf("failed to check file %s: %s\n", path, err)
		return false
	}

	return fi.Mode().IsRegular()
}
