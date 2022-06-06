package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	s := "hello 中国1"
	t := "hello 中国1"
	fmt.Println(HaveSameElements(s, t))
	count := "hello 中中国国人"
	fmt.Println("count:", strings.Count(count, ""), "bytelen:",
		len(count), "trueLen:", len([]rune(count)), "utf-8Len:", utf8.RuneCountInString(count))
}

//俩个字符串的元素是否完全一致,以至于重新排序之后能够相同?
func HaveSameElements(s, t string) bool {
	s1, t1 := []rune(s), []rune(t)

	if len(s1) != len(t1) {
		return false
	}
	//遍历其中一个字符串的元素,在另一个字符串中是否都存在,且数量相同
	for _, v := range s {
		if strings.Count(s, string(v)) != strings.Count(t, string(v)) {
			return false
		}

	}
	return true
}
