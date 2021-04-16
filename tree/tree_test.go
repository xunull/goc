package tree

import (
	"fmt"
	"testing"
)

type TestS struct {
	Path string
}

func TestPlantTree(t *testing.T) {
	t1 := TestS{
		Path: "aaa/bbb/ccc/ddd",
	}
	t2 := TestS{
		Path: "111/222/333/444",
	}
	t3 := TestS{
		Path: "aaa/bbb/eee",
	}
	t4 := TestS{
		Path: "aaa/bbb/eee/ccc",
	}

	tt := []interface{}{t1, t2, t3, t4}

	tree := PlantPathTree(tt, "Path")

	tree.Name = "ROOT"
	res := TreeIt([]TreeAble{tree})

	for _, item := range res {
		fmt.Println(item)
	}
}
