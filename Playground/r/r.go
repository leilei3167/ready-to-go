package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Printf("Hello world")

	fmt.Printf("\rhello world\n")

	fmt.Println(rightPad("hello", "dsaijczx", 10))

}

func rightPad(s string, padStr string, overallen int) string {
	strLen := len(s)
	if overallen <= strLen {
		return s
	}
	toPad := overallen - strLen - 1
	pad := strings.Repeat(padStr, toPad)
	return fmt.Sprintf("%s%s", s, pad)
}
