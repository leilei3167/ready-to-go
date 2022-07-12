package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var wg sync.WaitGroup

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go doSomething(ctx, &wg)
	}

	//阻塞
	select {
	case <-s:
		//主动退出
		cancel()
	case <-ctx.Done():
		//超时退出
	}

	wg.Wait()

}

func doSomething(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("err:", ctx.Err())
			return
		default:
			fmt.Println("working...")
			time.Sleep(time.Second)
		}
	}
}
