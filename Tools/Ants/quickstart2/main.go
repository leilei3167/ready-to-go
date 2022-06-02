package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
	"sync/atomic"
	"time"
)

var sum int32

func myFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
}
func demoFunc() {
	time.Sleep(time.Second)
	fmt.Println("Hello World!")
}

func main() {
	defer ants.Release() //本身就存在一个默认的池(NewPool创建的)
	runTimes := 10
	var wg sync.WaitGroup
	syncSum := func() {
		demoFunc()
		wg.Done()
	}
	wg.Add(runTimes)
	for i := 0; i < runTimes; i++ {
		_ = ants.Submit(syncSum)          //只能提交没有任何参数的函数,需要传入参数只能使用闭包
		_ = ants.Submit(WarpFunc(11, 22)) //需要传参数要利用闭包
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all tasks.\n")
	//带函数的池
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) { //限制worker数为10
		myFunc(i)
		wg.Done()
	})
	defer p.Release()
	// Submit tasks one by one.
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(int32(i))
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks, result is %d\n", sum)
}

//如果需要提交带参数的函数,可以利用闭包函数进行包裹

func WarpFunc(a, b int) func() {
	return func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			fmt.Println(a + b)
		}
	}
}
