package commonx

import (
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"strings"
)

func GetWords(target string) map[string]bool {
	res := make(map[string]bool)

	tl := strings.Split(target, "")
	length := len(tl)
	nums := make([]int, 0, len(tl))
	for i := 0; i < len(tl); i++ {
		nums = append(nums, i)
	}

	visited := make(map[int]bool)

	var dfs func(path []int)
	dfs = func(path []int) {
		if len(path) == len(nums) {
			t := make([]string, 0, length)
			for _, v := range path {
				t = append(t, tl[v])
			}
			res[strings.Join(t, "")] = true
		} else {
			for _, n := range nums {
				if visited[n] {
					continue
				}
				path = append(path, n)
				visited[n] = true
				dfs(path)
				path = path[:len(path)-1]
				visited[n] = false
			}
		}
	}
	dfs([]int{})
	return res
}

func JsonString(v interface{}) string {

	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("err：%v", err)
		return "JSON Marshal Failed"
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")
	if err != nil {
		return string(b)
	} else {
		return out.String()
	}
}
