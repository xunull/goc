package git_cmd

import (
	"github.com/xunull/goc/commonx"
	"regexp"
)

const (
	ReLogGraphBracket = `\([\s\S]*\)`
	ReOnlyTag         = `\(tag:[^,]*\)`
)

var (
	reLogGraphBracket *regexp.Regexp
	reOnlyTag         *regexp.Regexp
)

func init() {
	var err error
	reLogGraphBracket, err = regexp.Compile(ReLogGraphBracket)
	commonx.CheckErrOrFatal(err)
	reOnlyTag, err = regexp.Compile(ReOnlyTag)
}
