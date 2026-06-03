package lang_ext

import "strings"

var commonLanguageReverseExt = make(map[string]string)

var commonLanguageLowerReverseExt = make(map[string]string)

var knownFileExt = make(map[string]string)

var knowMapList = []map[string]string{commonLanguageExt,
	commonExcludeFileExt,
	excludeLineCount}

func init() {
	for k, v := range commonLanguageExt {
		commonLanguageReverseExt[v] = k
		commonLanguageLowerReverseExt[strings.ToLower(v)] = k
	}

	for _, list := range knowMapList {
		for k, v := range list {
			knownFileExt[k] = v
		}
	}

}

var commonFileName = map[string]string{
	"Makefile":   "Makefile",
	"makefile":   "makefile",
	"Dockerfile": "Dockerfile",
	"README.md":  "ReadMe",
	"Readme.md":  "ReadMe",
}
