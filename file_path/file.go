package file_path

import (
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
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

func GetFileName(target string) string {
	fs := filepath.Base(target)
	ext := filepath.Ext(target)
	return strings.TrimSuffix(fs, ext)
}

func GetFullFileName(target string) string {
	return filepath.Base(target)
}
