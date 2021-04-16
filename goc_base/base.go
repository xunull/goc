package goc_base

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type BaseStruct struct {
}

func (s *BaseStruct) OutString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%+v", *s)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *s)
	}
	return out.String()
}
