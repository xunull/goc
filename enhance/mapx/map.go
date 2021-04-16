package mapx

func GetStrRmmmap(data map[string][]string) map[string]map[string]string {
	res := make(map[string]map[string]string)
	for k, arr := range data {
		for _, name := range arr {
			if _, ok := res[name]; !ok {
				res[name] = make(map[string]string)
			}
			res[name][k] = k
		}
	}
	return res
}
