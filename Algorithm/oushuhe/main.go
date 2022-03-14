package main

import "fmt"

func main() {
	fmt.Println(my_sum_even(1, 5))

}

func my_sum_even(x int, y int) int {
	// write code here
	temp := 0
	for i := x + 1; i < y; i++ {
		if i%2 == 0 {
			temp = temp + i
		}
	}

	return temp
}

