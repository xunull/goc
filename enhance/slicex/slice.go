package slicex

import (
	"fmt"
)

func splitWithLength(arr []interface{}, length int) {

}

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
