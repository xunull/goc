package commonx

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

func MD5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

func StrMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func MD5Salt(str, salt string) string {
	return MD5([]byte(str + salt))
}

func MD5File(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	md5 := fmt.Sprintf("%x", md5.Sum(body))
	runtime.GC()
	return md5, nil
}
