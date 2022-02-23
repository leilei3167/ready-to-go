package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < 110; i++ {
			<-ch
			if i%2 == 1 {
				fmt.Println("协程1打印:", i)

			}

		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for i := 0; i < 110; i++ {
			ch <- 1
			if i%2 == 0 {
				fmt.Println("协程2打印:", i)
			}

		}
		wg.Done()
	}()
	wg.Wait()
}
