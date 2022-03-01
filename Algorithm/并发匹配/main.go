package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

/*模拟从海量数据中并发的去找寻匹配的值*/
func faces() []int {
	face := make([]int, 0)
	rand.Seed(time.Now().Unix())
	for i := 0; i < 10000; i++ {

		n := rand.Intn(101000)

		face = append(face, n)
	}

	//模拟产生10000个数据
	return face
}

var wg sync.WaitGroup

func main() {
	x := faces()
	for i, i2 := range x {
		fmt.Println(i, i2)
	}
	signal := make(chan int)
	signal2 := make(chan int)

	y := 21456 //找某一个数据

	ctx, cancle := context.WithCancel(context.Background())
	for i := 0; i < 100; i++ { //并发读并不产生竞态
		wg.Add(1)
		go certif(ctx, x, y, signal)
	}
	//开启一个协程监听
	go func() {
		for {
			select {
			case <-signal:
				cancle()
				log.Fatalln("已找到值,位置:", <-signal)
				return
			case <-signal2: //说明所有协程执行完毕
				fmt.Println("没有该值")
				signal2 <- 1
				return

			}
		}
	}()
	wg.Wait()    //等待处理任务的协程全部执行完毕
	signal2 <- 1 //其他协程处理完毕发送信号
	<-signal2    //等待消息打印完毕
	fmt.Println(runtime.NumGoroutine())
}

func certif(ctx context.Context, n []int, x int, c chan int) {
	for k, v := range n {
		if x == v {
			c <- k
			wg.Done()
			return
		}
	}
	wg.Done()
	return
}
