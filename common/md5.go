package common

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

func MD5Salt(str, salt string) string {
	return MD5([]byte(str + salt))
}
