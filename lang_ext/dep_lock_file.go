package lang_ext

var commonDepLockFile = map[string]string{
	".sum":       "*.sum",
	".lock":      "*.lock",
	".gitignore": ".gitignore",
}

func CommonDepLockFile() map[string]string {
	return commonDepLockFile
}
