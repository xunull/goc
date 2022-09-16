package lang_ext

var CommonExcludeFileName = map[string]string{
	"package.json":      "package.json",
	"package-lock.json": "package-lock.json",
}

var CommonExcludeDir = map[string]string{
	"dist":              "dist",
	"node_modules":      "node_modules",
	"public":            "public",
	"vendor":            "vendor",
	"cmake-build-debug": "cmake-build=debug",
}

var CommonExcludeFileExt = map[string]string{
	".exe": "exe",
	".pyc": "pyc",
	".bin": "bin",
	".dll": "dll",
	".pdb": "pdb",
	".pt":  "pt",
	".log": "log",
	".o":   "c/c++",
	".so":  "c/c++",
	".sum": "golang go.sum",
}

var ExcludeLineCount = map[string]string{
	".mp4":   "MP4",
	".png":   "Png",
	".tif":   "TIF",
	".jpg":   "Jpg",
	".jpeg":  "Jpg",
	".gif":   "GIF",
	".ico":   "Ico",
	".pdf":   "PDF",
	".ttf":   "TTF",
	".csv":   "CSV",
	".zip":   "Zip",
	".data":  "Data",
	".model": "Model",
	".svg":   "SVG",
	".gz":    "Tar File",
	".tgz":   "Tar File",
	".tar":   "Tar File",
	".sig":   "Signature File",
	".pyi":   "pyi",
	".pptx":  "ppt",
	".pth":   "PTH",
	".pkl":   "PKL",
	".h5":    "H5",
	".npy":   "npy",
	".npz":   "npz",
	".db":    "Database Files",
	".obj":   "OBJ",
	".ply":   "PLY",
	".pb":    "PB",
}
