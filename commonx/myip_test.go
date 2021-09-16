package commonx

import (
	"fmt"
	"testing"
)

func TestGetMyIp(t *testing.T) {
	ip, err := GetMyIp()
	CheckErrOrFatal(err)
	fmt.Println(ip)
}
