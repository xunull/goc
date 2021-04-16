package cachex

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	a := func(w, e string) (string, string) {
		return w, e
	}

	fmt.Println(a)
}
