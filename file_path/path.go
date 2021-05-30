package file_path

import (
	"github.com/mitchellh/go-homedir"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"
)

func IsAbs(p string) bool {
	return filepath.IsAbs(p)
}

func GetAbsPath(p string) (string, error) {
	t, err := homedir.Expand(p)
	if err != nil {
		return p, err
	}
	if filepath.IsAbs(t) {
		return t, nil
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return p, err
		} else {
			return filepath.Join(wd, t), nil
		}
	}
}

func NameNoExt(file fs.FileInfo) string {
	name := file.Name()
	ext := path.Ext(name)
	return strings.TrimSuffix(name, ext)
}

func PathIsDir(path string) (bool, error) {
	if info, err := os.Stat(path); err == nil {
		return info.IsDir(), err
	} else {
		return false, err
	}
}

func PathIsDirOrFatal(path string) bool {
	if info, err := os.Stat(path); err == nil {
		return info.IsDir()
	} else {
		log.Fatal(err)
	}
	return false
}

func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func PathExistOrCreate(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return false, err
			} else {
				return true, nil
			}
		} else {
			return false, err
		}
	}
	return false, nil
}

func PathExistsOrFatal(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			log.Fatal(err)
		}
	}
	return true
}

func GetBinDir() (string, error) {
	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		debug.PrintStack()
		log.Printf("can not get bin dir %s", err)
		return "", err
	} else {
		return dir, nil
	}
}

func GetBinDirOrFatal() string {
	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		log.Fatal(err)
	} else {
		return dir
	}
	return ""
}

func TargetContains(p, sub string) bool {
	pb := filepath.Base(p)
	return strings.Contains(pb, sub)
}

func RemovePrefixN(target string, n int) string {
	return strings.Join(strings.Split(target, "/")[n:], "/")
}

func RemoveSuffixN(target string, n int) string {
	list := strings.Split(target, "/")
	length := len(list)
	return strings.Join(list[:length-n], "/")
}
