package traverse

import (
	"fmt"
	"os"
	"testing"
)

type TestStruct struct {
}

func (t *TestStruct) Test1(item *TraverseItem) {
	fmt.Println(item.Fs.Name())
}

func TestDirTraverse(t *testing.T) {
	p, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("path: %s\n", p)
	dir := NewDirTraverse(p)

	ts := TestStruct{}

	dir.Handle(ts.Test1)
}
