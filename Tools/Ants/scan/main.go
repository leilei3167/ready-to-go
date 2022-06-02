package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"net"
	"runtime"
	"time"
)

func main() {
	//设置提前分配内存,将使得其内部的数据结构选择循环队列
	pool, _ := ants.NewPool(100, ants.WithPreAlloc(true))
	defer pool.Release()
	result := make(chan string, 1)

	go func() {
		for i := 1; i < 256; i++ {
			ip := "182.61.6.67"
			port := i * 2
			hosts := fmt.Sprintf("%v"+":"+"%v", ip, port)
			_ = pool.Submit(ScanWithWarp(result, hosts))
		}
	}()

	fmt.Println("阻塞的任务:", pool.Waiting(), "总G数量:", runtime.NumGoroutine())
	for i := 0; i < 255; i++ {
		//	fmt.Println(<-result)
		<-result
	}
	fmt.Println(pool.Running())
	time.Sleep(time.Millisecond * 2000)
	fmt.Println("过期时间后:", pool.Running()) //默认1s间隔会清除G
	fmt.Println("任务结束")
}

func ScanWithWarp(resultChan chan<- string, host string) func() {
	return func() {
		conn, err := net.DialTimeout("tcp", host, time.Second)
		if err != nil {
			resultChan <- err.Error()
			return
		}
		conn.Close()
		resultChan <- fmt.Sprintf("%s is OPEN!", host)
	}
}
