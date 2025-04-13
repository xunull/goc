package file_utils

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// FileCount 这个会递归的调用
func FileCount(dirname string) (int, error) {
	count := 0
	sub := make([]string, 0)
	if files, err := ioutil.ReadDir(dirname); err == nil {
		for _, file := range files {
			if file.IsDir() {
				sub = append(sub, file.Name())
			} else {
				count += 1
			}
		}

		for _, file := range sub {
			if c, err := FileCount(path.Join(dirname, file)); err == nil {
				count += c
			} else {
				return 0, nil
			}
		}
		return count, nil
	} else {
		return count, err
	}
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func CountDirFiles(dirName string) (int, error) {
	if !IsDir(dirName) {
		return 0, nil
	}
	var count int
	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		count++
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}
