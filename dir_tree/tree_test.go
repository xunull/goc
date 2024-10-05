package dir_tree

import (
	"fmt"
	"os"
	"testing"
)

func TestTree(t *testing.T) {
	dir := os.Getenv("goc_dir_tree_test_target_dir")
	if dir == "" {
		fmt.Printf("please set env goc_dir_tree_test_target_dir\n")
		return
	}
	tree := NewTree(dir, HandlerFunc{
		DirFunc: func(obj *TreeItem) {
			fmt.Printf("Dir: %s %s \n", obj.Name(), obj.Parent)

		},
		FileFunc: func(obj *TreeItem) {
			fmt.Printf("File: %s %s \n", obj.Name(), obj.Parent)
		},
		FinishFunc: func() {
			fmt.Printf("finish\n")
		},
	})

	tree.Tree()
}

func TestTreeWithDefaultExclude(t *testing.T) {
	dir := os.Getenv("goc_dir_tree_test_target_dir")
	if dir == "" {
		fmt.Printf("please set env goc_dir_tree_test_target_dir\n")
		return
	}
	tree := NewTree(dir, HandlerFunc{
		DirFunc: func(obj *TreeItem) {
			fmt.Printf("Dir: %s %s \n", obj.Name(), obj.Parent)

		},
		FileFunc: func(obj *TreeItem) {
			fmt.Printf("File: %s %s \n", obj.Name(), obj.Parent)
		},
		FinishFunc: func() {
			fmt.Printf("finish\n")
		},
	}, WithDepth(2))

	tree.Tree()
}
