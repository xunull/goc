package lang_ext

var CommonLanguageExt = map[string]string{
	".cfg":   "Config",
	".go":    "Golang",
	".java":  "Java",
	".py":    "Python",
	".lua":   "Lua",
	".j2":    "Jinja",
	".c":     "C",
	".cpp":   "C++",
	".js":    "JavaScript",
	".ts":    "TypeScript",
	".vue":   "Vue",
	".json":  "Json",
	".yaml":  "Yaml",
	".yml":   "Yaml",
	".ini":   "Ini",
	".md":    "Markdown",
	".sh":    "Shell",
	".html":  "Html",
	".pl":    "Perl",
	".perl":  "Perl",
	".xml":   "XML",
	"Grpc":   "Grpc",
	".proto": "Grpc",
	".txt":   "Text",
}

var CommonLanguageReverseExt = make(map[string]string)

func init() {
	for k, v := range CommonLanguageExt {
		CommonLanguageReverseExt[v] = k
	}
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

var CommonFileName = map[string]string{
	"Makefile":   "Makefile",
	"makefile":   "makefile",
	"Dockerfile": "Dockerfile",
	"README.md":  "ReadMe",
	"Readme.md":  "ReadMe",
}
