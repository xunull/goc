package file_utils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CopyFile(dst, src string) error {
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	_, err = io.Copy(d, s)
	return err
}

func CopyFileIntoDirOrFatal(src string, dir string) error {
	name := filepath.Base(src)
	dstName := filepath.Join(dir, name)
	return CopyFile(dstName, src)
}

func CopyFileIntoDir(src string, dir string) error {
	name := filepath.Base(src)
	dstName := filepath.Join(dir, name)
	return CopyFile(dstName, src)
}

func ListFileNames(dirname string) ([]string, error) {

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	} else {
		res := make([]string, 0, len(files))
		for _, file := range files {
			res = append(res, file.Name())
		}
		return res, nil
	}
}
