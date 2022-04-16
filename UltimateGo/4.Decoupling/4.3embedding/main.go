package main

import "fmt"

//嵌套

func main() {
	a := []int{1, 2, 3, 4, 5, 5}
	fmt.Printf("len %v cap%v \n", len(a), cap(a))
}
