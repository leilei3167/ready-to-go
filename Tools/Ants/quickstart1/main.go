package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"math/rand"
	"sync"
)

const (
	DataSize    = 10000
	DataPerTask = 100
)

func main() {
	//创建工作池,taskFunc是其中的工作内容
	p, _ := ants.NewPoolWithFunc(10, taskFunc)
	defer p.Release() //延迟关闭协程池
	nums := make([]int, DataSize, DataSize)
	for i := range nums {
		nums[i] = rand.Intn(1000)
	}

	var wg sync.WaitGroup
	wg.Add(DataSize / DataPerTask) //100个任务
	tasks := make([]*Task, 0, DataSize/DataPerTask)
	for i := 0; i < DataSize/DataPerTask; i++ {
		task := &Task{
			index: i + 1,
			nums:  nums[i*DataPerTask : (i+1)*DataPerTask],
			wg:    &wg,
		}
		tasks = append(tasks, task)
		p.Invoke(task) //每一个任务获取协程池来执行,池中任务执行顺序是随机的
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())
	//可以使用NewPool创建池
}

type Task struct {
	index int
	nums  []int
	sum   int
	wg    *sync.WaitGroup
}

func (t *Task) Do() {
	for _, num := range t.nums {
		t.sum += num
	}

	t.wg.Done()
}
func taskFunc(data interface{}) {
	task := data.(*Task)
	task.Do()
	fmt.Printf("task:%d sum:%d\n", task.index, task.sum)
}
