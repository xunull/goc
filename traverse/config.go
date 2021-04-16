package traverse

var CommonLanguageExt = map[string]string{
	".go":   "Golang",
	".java": "Java",
	".py":   "Python",
	".lua":  "Lua",
	".c":    "C",
	".cpp":  "C++",
	".js":   "JavaScript",
	".ts":   "TypeScript",
	".vue":  "Vue",
	".json": "Json",
	".yaml": "Yaml",
	".ini":  "ini",
	".md":   "Markdown",
	".sh":   "Shell",
	".html": "Html",
	".pl":   "Perl",
	".perl": "Perl",
	".xml":  "XML",
}

var ExcludeLineCount = map[string]string{
	".png":  "Png",
	".jpg":  "Jpg",
	".jpeg": "Jpg",
	".gif":  "Gif",
	".ico":  "Ico",
}

var CommonExcludeDir = map[string]string{
	"dist":         "dist",
	"node_modules": "node_modules",
	"public":       "public",
}

var CommonExcludeFileExt = map[string]string{
	".exe": "exe",
	".pyc": "pyc",
}

var CommonFrontLanguageExt = map[string]string{
	".html": "html",
	".vue":  "vue",
	".css":  "css",
	".js":   "js",
}

var CommonBackLanguageExt = map[string]string{
	".go":   "golang",
	".java": "java",
	".py":   "python",
}

var CommonExcludeFileName = map[string]string{
	"package.json":      "package.json",
	"package-lock.json": "package-lock.json",
}
