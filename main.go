package main

import "fmt"

func main() {
	//var i []int
	var i = make([]int, 1)

	var p map[string]string

	fmt.Printf("i is %T,%#v", i, i)
	fmt.Printf("p is %T,%#v", p, p)

}
