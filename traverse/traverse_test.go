package traverse

import (
	"fmt"
	"os"
	"testing"
)

type TestStruct struct {
}

func (t *TestStruct) Test1(item *TraverseItem) {
	fmt.Println(item.FileInfo.Name())
}

func TestDirTraverse(t *testing.T) {
	p, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("path: %s\n", p)
	ts := TestStruct{}
	dir := NewDirTraverse(p, ts.Test1)
	defer dir.Close() // 添加资源清理

	err = dir.Handle()
	if err != nil {
		t.Logf("Handle completed with errors: %v", err)
	}

	// 检查是否有错误
	if dir.HasErrors() {
		t.Logf("Collected errors: %v", dir.Errors())
	}
}
