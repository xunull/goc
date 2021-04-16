package outputx

import (
	"fmt"
	"github.com/xunull/goc/commonx"
)

func OutputJsonString(v interface{}) {
	t, _ := commonx.JsonStringError(v)
	fmt.Println(t)
}
