package relib

func RemoveSpace(str string) string {
	return ReSpace.ReplaceAllString(str, "")
}
