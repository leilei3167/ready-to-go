package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 100; i++ {
		fmt.Println("hello world!!!!!!")
		time.Sleep(time.Second * 2)
	}
}
