package slicex

import (
	"fmt"
	"strconv"
)

func splitWithLength(arr []interface{}, length int) {

}

// todo 这个方法有人调用么
func splitWithCount(arr []interface{}, count int) {
	if count <= 1 {

	}
	length := len(arr)
	sl := length / count
	fmt.Println(sl)
}

func IsIncreasingInt64(arr []int64) bool {
	length := len(arr)
	if length >= 2 {
		for i := 1; i < length; i++ {
			if arr[i] < arr[i-1] {
				return false
			}
		}
		if arr[length-1] == arr[0] {
			return false
		}
		return true
	} else {
		return false
	}
}

func IntToString(target []int) []string {
	res := make([]string, 0, len(target))
	for _, item := range target {
		res = append(res, strconv.Itoa(item))
	}
	return res
}

func SliceStringToInterface(arr []string) []interface{} {
	s := make([]interface{}, len(arr))
	for i, v := range arr {
		s[i] = v
	}
	return s
}
