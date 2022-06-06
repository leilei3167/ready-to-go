package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

	wg.Add(2)
	ch := make(chan bool, 1)
	letter := make(chan bool, 1)

	go func() {
		var i = 0
		for {
			select {
			case s := <-ch:
				if !s {
					wg.Done()
					return
				}
				fmt.Println(i)
				i++
			case letter <- true:

			}
		}

	}()

	go func() {
		var i = 'A'

		for {
			select {
			case <-letter:
				if i > 'Z' {
					wg.Done()
					ch <- false
					return
				}
			}
			fmt.Println(string(i))
			i++
			ch <- true
		}
	}()

	wg.Wait()
}
