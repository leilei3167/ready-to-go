package main

import "fmt"

func main() {
	before := []int{1, 2, 3, 4, 5}
	after := revers(before)
	fmt.Println(after)

}

func revers(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		//i:=0,j:=len(s)-1;i<j;i=i+1,j=j-1
		s[i], s[j] = s[j], s[i]

	}
	return s
}
