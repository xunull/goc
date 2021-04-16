package commonx

import (
	"github.com/go-resty/resty/v2"
	"regexp"
)

var ipRegStr = "(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})(\\.(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})){3}"
var ipReg = regexp.MustCompile(ipRegStr)

func GetMyIp() (string, error) {
	resp, err := resty.New().R().Get("http://myip.ipip.net/")
	if err != nil {
		return "", err
	}
	return ipReg.FindString(resp.String()), nil
}
