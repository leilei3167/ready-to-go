package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)
//https://segmentfault.com/a/1190000039710281
var (
	Limit  = flag.Int64("l", 1, "G上限")
	Weight = flag.Int64("w", 1, "信号量的权重")
)

func main() {
	flag.Parse()
//使用NewWeighted创建一个并发访问的最大资源数,此处即是设置可同时运行
//的G上限为多少个
	sem := semaphore.NewWeighted(*Limit)
	var w sync.WaitGroup

	for i := 0; i < 1000000; i++ {
		w.Add(1)

		go func(i int) {
			//Acquire来获取这个指定个数的资源,如果没有空闲资源
			//则当前G陷入休眠,使用完成之后记得要Release
			sem.Acquire(context.Background(), *Weight)
			fmt.Println("i=", i)
			time.Sleep(time.Second * 3)
			sem.Release(*Weight)
			w.Done()
		}(i)
	}
	//虽然有1000个协程,但没有获取到资源的都是在休眠状态
	fmt.Println(runtime.NumGoroutine())
	w.Wait()
	log.Println("全部结束!")
}
