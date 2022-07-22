package main

import (
	"fmt"
)

func main() {

	a := make(chan int, 10)

	src := a

	fmt.Printf("src addr: %p\n", &src)
	fmt.Printf("a addr: %p\n", a)

	go func() {
		for i := 0; i < 10; i++ {
			a <- i
		}
		close(a)
	}()

	for s := range src {
		fmt.Printf("got %d from src\n", s)
	}

	fmt.Println("a 关闭后 src成功感知到")

}
