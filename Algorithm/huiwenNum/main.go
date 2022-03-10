package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println(isPal(12345654321))
	fmt.Println(isPal(12345654322))

}

func isPal(x int) bool {
	//转换为字符串
	c := strconv.Itoa(x)
	for i := range c {
		if c[len(c)-1-i] != c[i] {
			return false
		}

	}
	return true
}

func isPal2(x int) bool {
	if x < 0 {
		return false
	}
	val := 0
	temp := x
	for x != 0 {
		val = val*10 + x%10
		x = x / 10
	}
	if val == temp {
		return true
	} else {
		return false
	}

}
