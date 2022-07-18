package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch1 := make(chan bool)
	ch2 := make(chan bool)
	ch3 := make(chan bool)
	ch4 := make(chan bool)

	var wg = sync.WaitGroup{}

	wg.Add(4)

	go func() {
		for {
			fmt.Println(1)
			time.Sleep(time.Second)
			ch2 <- true
			//发送完信号后阻塞 直到ch4
			<-ch1
		}
		wg.Done()
	}()
	go func() {
		for {
			<-ch2
			fmt.Println(2)
			time.Sleep(time.Second)
			ch3 <- true
		}
		wg.Done()
	}()
	go func() {
		for {
			<-ch3
			fmt.Println(3)
			time.Sleep(time.Second)
			ch4 <- true
		}
		wg.Done()
	}()
	go func() {
		for {
			<-ch4
			fmt.Println(4)
			time.Sleep(time.Second)
			ch1 <- true
		}
		wg.Done()
	}()
	wg.Wait()

}
