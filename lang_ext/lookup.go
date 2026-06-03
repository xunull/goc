package lang_ext

// 这些 lookup 函数是 lang_ext 包对外的只读访问入口。
// 底层 map 在 init 阶段一次性构建完毕，之后通过本文件函数暴露，
// 既阻止外部写入（编译期保证），也让运行期 race 风险消失。

// LanguageOf 根据文件扩展名（含点，如 ".go"）返回对应的语言名。
func LanguageOf(ext string) (string, bool) {
	v, ok := commonLanguageExt[ext]
	return v, ok
}

// IsExcludeDir 判断目录名是否在通用排除列表中（如 node_modules、vendor）。
func IsExcludeDir(name string) bool {
	_, ok := commonExcludeDir[name]
	return ok
}

// IsExcludeFileExt 判断扩展名是否属于通用排除类型（如 .exe、.so、.pyc）。
func IsExcludeFileExt(ext string) bool {
	_, ok := commonExcludeFileExt[ext]
	return ok
}

// IsExcludeFileName 判断文件名是否属于通用排除清单（如 package-lock.json）。
func IsExcludeFileName(name string) bool {
	_, ok := commonExcludeFileName[name]
	return ok
}

// IsSkipLineCountExt 判断扩展名是否属于"跳过行数统计"的二进制/媒体类型。
func IsSkipLineCountExt(ext string) bool {
	_, ok := excludeLineCount[ext]
	return ok
}

// ExtOfLanguage 用语言名（如 "Golang"）反查扩展名。
func ExtOfLanguage(language string) (string, bool) {
	v, ok := commonLanguageReverseExt[language]
	return v, ok
}

// ExtOfLanguageLower 用小写语言名（如 "golang"）反查扩展名。
func ExtOfLanguageLower(language string) (string, bool) {
	v, ok := commonLanguageLowerReverseExt[language]
	return v, ok
}

// KnownFileExtOf 查询是否为合并过的"已知扩展名"集合中的项。
func KnownFileExtOf(ext string) (string, bool) {
	v, ok := knownFileExt[ext]
	return v, ok
}

// CommonFileNameOf 判断特殊文件名（Makefile、Dockerfile 等）。
func CommonFileNameOf(name string) (string, bool) {
	v, ok := commonFileName[name]
	return v, ok
}

// FrontLanguageExts 返回前端语言扩展名映射的一份拷贝；
// 拷贝避免调用方修改影响包内不变性。
func FrontLanguageExts() map[string]string {
	return cloneStrMap(commonFrontLanguageExt)
}

// BackLanguageExts 返回后端语言扩展名映射的一份拷贝。
func BackLanguageExts() map[string]string {
	return cloneStrMap(commonBackLanguageExt)
}

func cloneStrMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
