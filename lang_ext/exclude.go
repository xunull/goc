package lang_ext

var CommonExcludeFileName = map[string]string{
	"package.json":      "package.json",
	"package-lock.json": "package-lock.json",
}

var CommonExcludeDir = map[string]string{
	"dist":         "dist",
	"node_modules": "node_modules",
	"public":       "public",
	"vendor":       "vendor",
}

var CommonExcludeFileExt = map[string]string{
	".exe": "exe",
	".pyc": "pyc",
	".bin": "bin",
	".dll": "dll",
	".pdb": "pdb",
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
	".pyi":   "pyi",
	".pptx":  "ppt",
	".pth":   "PTH",
	".h5":    "H5",
}
