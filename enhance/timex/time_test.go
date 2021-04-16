package timex

import (
	"fmt"
	"strconv"
	"testing"
)

func Test1(t *testing.T) {
	var a float64 = float64(20) / float64(60)
	fmt.Println(a)
	strconv.FormatFloat(float64(20/60), 'f', 1, 64)
	fmt.Printf("%.2f\n", a)
	fmt.Printf(strconv.FormatFloat(float64(20/60), 'f', 2, 64))
}
