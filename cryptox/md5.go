package cryptox

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
)

func MD5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

func Md5Salt(str, salt string) string {
	return MD5([]byte(str + salt))
}

func Md5FilePath(fp string) (string, error) {
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return "", err
	}
	h := md5.New()
	return hex.EncodeToString(h.Sum(b)), nil
}

func Md5File(file *os.File) (string, error) {
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	h := md5.New()
	return hex.EncodeToString(h.Sum(b)), nil
}
