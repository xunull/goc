package commonx

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
)

func SaveMapInCurrentWorkDir(name string, m interface{}) error {
	r, err := json.Marshal(m)
	if err != nil {
		return err
	} else {
		return ioutil.WriteFile(name, r, fs.ModePerm)
	}
}
