package file_path

import (
	"fmt"
	"github.com/xunull/goc/commonx"
	"github.com/xunull/goc/enhance/timex"
	"os"
	"path/filepath"
)

func MakeCurTimeDir(parent string, opts ...Option) (string, error) {
	op := getOption(opts...)
	var name string
	if op.Suffix != "" {
		name = fmt.Sprintf("%s_%s", timex.GetYYYYMMDDHHMMSS(), op.Suffix)
	} else {
		name = timex.GetYYYYMMDDHHMMSS()
	}

	p := filepath.Join(parent, name)
	err := os.MkdirAll(p, os.ModePerm)
	return p, err
}

func MakeDirOrFatal(p string) {
	exist, err := PathExists(p)
	commonx.CheckErrOrFatal(err)
	if !exist {
		err = os.Mkdir(p, 0700)
		commonx.CheckErrOrFatal(err)
	}
}
