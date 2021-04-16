package rand_fake

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	for i := 0; i < 20; i++ {
		a := FakeGameName()
		fmt.Println(a)
	}
}
