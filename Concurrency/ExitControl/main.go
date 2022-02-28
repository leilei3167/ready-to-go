package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func worker1(done chan struct{}) chan int { //通过返回的chan来进行链接
	//done用来接收退出信号
	ch := make(chan int, 5)
	go func() {
	Lable:
		for {
			select {
			case ch <- rand.Int():
			case <-done: //如果收到退出信号
				break Lable
			}
		}
		close(ch)
	}()
	return ch
}
func worker2(done chan struct{}) chan int {
	//done用来接收退出信号
	ch := make(chan int, 5)
	go func() {
	Lable:
		for {
			select {
			case ch <- rand.Int():
			case <-done: //如果收到退出信号
				break Lable
			}
		}
		close(ch)
	}()
	return ch
}

//用select来扇入
func fanIn(done chan struct{}) chan int {
	ch := make(chan int)
	send := make(chan struct{}) //给聚合的chan发信息
	go func() {
	Lable:
		for {
			select {
			case ch <- <-worker1(send):
				time.Sleep(time.Second)
			case ch <- <-worker2(send):
				time.Sleep(time.Second) //将两个worker聚合在此
			case <-done: //如果收到退出信号,则发出两个send(对于worker是done)
				send <- struct{}{}
				send <- struct{}{}
				break Lable

			}

		}
		close(ch)
	}()
	return ch
}

func main() {
	//创建一个控制信号的chan
	done := make(chan struct{})
	//启动生成器
	ch := fanIn(done)
	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}

	//取到需要的数量后再发出退出信号
	done <- struct{}{}
	fmt.Println(runtime.NumGoroutine())
	fmt.Println("stoped")

}
