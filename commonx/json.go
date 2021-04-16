package commonx

import (
	"bytes"
	"encoding/json"
)

func FormatJson(b []byte) (string, error) {
	var str bytes.Buffer
	err := json.Indent(&str, b, "", "    ")
	return str.String(), err
}

func FormatJsonStr(src string) (string, error) {
	var str bytes.Buffer
	err := json.Indent(&str, []byte(src), "", "    ")
	return str.String(), err
}
