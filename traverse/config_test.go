package traverse

import (
	"fmt"
	"path"
	"testing"
)

func TestConfig(t *testing.T) {
	a := "test.go"
	fmt.Println(path.Ext(a))
	fmt.Println(path.Dir(a))

	b := ".git"
	fmt.Println(path.Ext(b))
	fmt.Println(path.Dir(b))

}

func Test2(t *testing.T) {

}
