package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		fmt.Println("hello world!")
		time.Sleep(time.Second)
	}
}
