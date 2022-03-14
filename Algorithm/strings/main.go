package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	n := bufio.NewReader(os.Stdin)
	x, _ := n.ReadString('\n')
	fmt.Println(fun1(x))
}

func fun1(s string) int {
	a := strings.Fields(s)
	n := len(a)
	return len(a[n-1])

}
