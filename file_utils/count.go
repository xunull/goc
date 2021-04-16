package file_utils

import (
	"io/ioutil"
	"path"
)

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
