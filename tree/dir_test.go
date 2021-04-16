package tree

import (
	"fmt"
	"os"
	"testing"
)

func TestDirTree(t *testing.T) {
	p, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("path: %s", p)
	ppt := DirTree(p)
	res := TreeIt([]TreeAble{ppt})

	for _, item := range res {
		fmt.Println(item)
	}

}
