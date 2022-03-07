package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(count2("Hello there, world! What’s going on there? "))
	fmt.Println(count2("我,我,你,他"))
	fmt.Println(count2("Package bufio implements buffered I/O. " +
		"It wraps an io.Reader or io.Writer object, creating " +
		"another object Reader or Writer that also implements " +
		"the interface but provides buffering and some help for " +
		"textual I/O,Reader.Reader" +
		".Reader"))
}

func count2(s string) (string, int) {
	//返回结果

	m := make(map[string]int)
	v := make([]string, 0)
	//自定义分隔方式(用空格分隔用fields函数,指定单个分隔用split
	f := func(r rune) bool {
		if r == ',' || r == ' ' || r == '!' || r == '\n' || r == '?' || r == '.' || r == '`' {
			return true
		} else {
			return false
		}
	}
	v = strings.FieldsFunc(s, f)
	for _, v := range v {
		data := v
		m[data] = m[data] + 1
	}
	//统计map中value最大的key-value并返回

	var tem = 0 //临时变量
	var rk string
	var rv int
	for s2, i := range m {
		if i >= tem {
			tem = i
			rk = s2
			rv = i
		}
	}
	return rk, rv

}
