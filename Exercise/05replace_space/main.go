package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "hello wor ld!!!"
	fmt.Println(replaceBlank(s))

}

func replaceBlank(s string) string {
	return strings.ReplaceAll(s, " ", "%20") //replace和replaceAll的用法
}
