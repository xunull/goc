package lang_ext

import "strings"

var CommonLanguageExt = map[string]string{
	".cfg":   "Config",
	".go":    "Golang",
	".java":  "Java",
	".py":    "Python",
	".lua":   "Lua",
	".j2":    "Jinja",
	".c":     "C",
	".cpp":   "C++",
	".h":     "C/C++ Header",
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
	".css":   "CSS",
	".pl":    "Perl",
	".perl":  "Perl",
	".xml":   "XML",
	"Grpc":   "Grpc",
	".proto": "Grpc",
	".txt":   "Text",
	".ipynb": "Jupyter",
	".m4":    "M4",
	".am":    "AutoMake",
	".texi":  "TEXINFO",
	".po":    "PO File",
	".awk":   "awk",
	".rc":    "Windows Resource File",
	".bmp":   "Bitmap",
	".tex":   "TeX",
	".php":   "PHP",
	".conf":  "Configuration File",
	".ac":    "Autoconf Script",
}

var CommonLanguageReverseExt = make(map[string]string)

var CommonLanguageLowerReverseExt = make(map[string]string)

func init() {
	for k, v := range CommonLanguageExt {
		CommonLanguageReverseExt[v] = k
		CommonLanguageLowerReverseExt[strings.ToLower(v)] = k
	}
}

var ExcludeLineCount = map[string]string{
	".png":   "Png",
	".PNG":   "Png",
	".jpg":   "Jpg",
	".jpeg":  "Jpg",
	".gif":   "GIF",
	".ico":   "Ico",
	".pdf":   "PDF",
	".ttf":   "TTF",
	".TTF":   "TTF",
	".csv":   "CSV",
	".zip":   "Zip",
	".data":  "Data",
	".model": "Model",
	".pkl":   "pkl",
	".svg":   "SVG",
	".gz":    "Tar File",
	".tgz":   "Tar File",
	".sig":   "Signature File",
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
