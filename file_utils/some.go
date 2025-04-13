package file_utils

import (
	"fmt"
	"os"
)

func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

func MkDirs(dirs ...string) error {
	if len(dirs) == 0 {
		return nil
	}
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create %s, %v", dir, err)
		}
	}
	return nil
}
