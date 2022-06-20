package fuzz

import (
	"errors"
	"unicode/utf8"
)

func isEqual(arr1, arr2 []byte) bool {

	for i, _ := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}

	}
	return true
}

//字符串本身是只读的字节切片,而将其转换为[]rune时,go将字节编码为utf-8编码的rune,用替换过的utf-8和原始输入的相比较肯定是不相等的
func Reverse(s string) (string, error) {
	if !utf8.ValidString(s) { //避免输入字符串不合法
		return s, errors.New("input is not valid UTF-8")
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}

//数组求和
func Sum(vals []int64) int64 {
	var total int64

	for _, val := range vals {

		total += val
	}

	return total
}
