package lang_ext

import "strings"

var CommonLanguageReverseExt = make(map[string]string)

var CommonLanguageLowerReverseExt = make(map[string]string)

var KnownFileExt = make(map[string]string)

var knowMapList = []map[string]string{CommonLanguageExt,
	CommonExcludeFileExt,
	ExcludeLineCount}

func init() {
	for k, v := range CommonLanguageExt {
		CommonLanguageReverseExt[v] = k
		CommonLanguageLowerReverseExt[strings.ToLower(v)] = k
	}

	for _, list := range knowMapList {
		for k, v := range list {
			KnownFileExt[k] = v
		}
	}

}

var CommonFileName = map[string]string{
	"Makefile":   "Makefile",
	"makefile":   "makefile",
	"Dockerfile": "Dockerfile",
	"README.md":  "ReadMe",
	"Readme.md":  "ReadMe",
}
