package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(basename("/usr/local/go.txt"))  //"go"
	fmt.Println(basename2("/usr/local/go.txt")) //"go"
	fmt.Println(comma("helloworld"))
	fmt.Println(comma("helloworld你好"))

}

func basename(s string) string {
	// Discard last '/' and everything before.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}
	// Preserve everything before last '.'.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	return s
}

func basename2(s string) string {
	index := strings.LastIndex(s, "/")
	s = s[index+1:] //去除/以及前面的
	index2 := strings.LastIndex(s, ".")
	s = s[:index2]
	return s

}

//每隔3个字节以,分隔
func comma(s string) string {
	//	n := utf8.RuneCountInString(s)
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]

}
