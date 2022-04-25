package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"

	"github.com/panjf2000/ants/v2"
)

//Ants一个高性能协程池
/*
   //过多的G直接导致程序崩溃
   var wg sync.WaitGroup
   	wg.Add(10000000)
   	for i := 0; i < 10000000; i++ {
   		go func() {
   			time.Sleep(1 * time.Minute)
   		}()
   	}
   	wg.Wait()
*/
/* 对于G的产生不加以限制,会造成资源耗尽导致程序崩溃;另一方面G的管理也是一个问题,G只能自己结束,管理不当会造成内存泄漏,并且频繁的开启创建G也是额外的开销
因此就有了协程池这样的需求
*/

//计算大量整数的程序
type Task struct {
	index int
	nums  []int
	sum   int
	wg    *sync.WaitGroup
}

func (t *Task) Do() { //将切片所有元素相加
	for _, num := range t.nums {
		t.sum += num
	}
	t.wg.Done()
}

//执行任务
func taskFunc(data interface{}) {
	task := data.(*Task)
	task.Do()
	fmt.Printf("task:%d sum:%d\n", task.index, task.sum)
}

const (
	DataSize    = 10000
	DataPerTask = 100
)

func main() {
	//1.创建G池,第一个参数是容量,代表池中最多有10个G,第二个参数是执行任务的函数
	//掉哦那个p.Invoke(data)的时候,会在池中找出一个空闲的G,让他执行该函数,并以data为参数
	p, _ := ants.NewPoolWithFunc(10, taskFunc)
	defer p.Release() //池用完需要关闭

	//生成10000个整数,分为1000个任务
	nums := make([]int, DataSize, DataSize)
	for i := range nums {
		nums[i] = rand.Intn(1000)
	}

	var wg sync.WaitGroup
	wg.Add(DataSize / DataPerTask)
	tasks := make([]*Task, 0, DataSize/DataPerTask)
	for i := 0; i < DataSize/DataPerTask; i++ {
		task := &Task{
			index: i + 1,
			nums:  nums[i*DataPerTask : (i+1)*DataPerTask], //任务拆分
			wg:    &wg,
		}

		tasks = append(tasks, task)
		p.Invoke(task) //执行单个任务

	}
	fmt.Println(runtime.NumGoroutine())
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())

	var sum int
	for _, task := range tasks {
		sum += task.sum //单个任务的结果相加
	}

	var expect int
	for _, num := range nums {
		expect += num //总任务的相加
	}

	fmt.Printf("finish all tasks, result is %d expect:%d\n", sum, expect)

}
