package main

import "fmt"

func main() {
	fmt.Println(revers("hello world! 你好中国!"))
}
func revers(s string) string {
	str := []rune(s)
	n := len(str)
	if n <= 1 {
		return s
	}

	for i, j := 0, n-1; i < j; {
		str[i], str[j] = str[j], str[i]
		i++
		j--
	}
	return string(str)
}
