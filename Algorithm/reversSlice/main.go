package main

import "fmt"

func main() {
	before := []int{1, 2, 3, 4, 5}
	after := revers(before)
	fmt.Println("reverse1", after)
	fmt.Println("test:", revers2("abcdef"))
}

func revers(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		//i:=0,j:=len(s)-1;i<j;i=i+1,j=j-1
		s[i], s[j] = s[j], s[i]

	}
	return s
}

//翻转字符串,转化为byte切片后,通过交换下标的形式
func revers2(s string) string {
	s1 := []byte(s)
	for i, j := 0, len(s1)-1; i < j; i, j = i+1, j-1 {
		s1[i], s1[j] = s1[j], s1[i]

	}
	return string(s1)
}
