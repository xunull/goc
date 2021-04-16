package file_utils

import (
	"testing"
)

func TestMakeTempFiles(t *testing.T) {
	err := MakeTempFiles("N:\\Temp\\418", 5, 3)
	if err != nil {
		t.Logf("open url failed, error is :%v", err)
	}
}
