package common

func IsTargetIn(target string, list []string) bool {
	for _, item := range list {
		if target == item {
			return true
		}
	}
	return false
}
