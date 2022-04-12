package main

import (
	"fmt"
	"net"
	"os"
	"runtime"

	"runtime/pprof"
	"strconv"
	"sync"
	"time"
)

//----------worker 的结构 保存当前的状态,并声明方法控制开始停止
type Worker struct {
	ID       int
	Name     string
	StopChan chan bool
}

//start,传入一个任务队列的管道
func (w *Worker) Start(jobQueue chan Job) {
	w.StopChan = make(chan bool)
	successChan := make(chan bool)

	go func() {
		successChan <- true
		for { 
			//获取任务,没有任务会阻塞;当pool关闭jobQueue时会读取到nil
			job := <-jobQueue
			if job != nil {
				job.Start(w)
			} else {
				//fmt.Printf("worker %s to be stopped\n", w.Name)
				w.StopChan <- true
				break

			}

		}

	}()
	<-successChan
}

func (w *Worker) stop() {
	<-w.StopChan
	//	fmt.Printf("worker %s stopped\n", w.Name)

}

//------Job是包含单方法的Start的接口,只要实现Start方法就可以有不同类型的job
type Job interface {
	Start(worker *Worker) error
}

//-----------创建pool管理worker
type Pool struct {
	Name      string
	Size      int //pool大小
	Workers   []*Worker
	QueueSize int      //任务队列的大小
	Queue     chan Job //任务队列
}

//初始化pool
func (p *Pool) Initialize() {
	if p.Size < 1 {
		p.Size = 1
	}
	//根据指定的size来创建worker
	p.Workers = []*Worker{}
	for i := 0; i < p.Size; i++ {
		worker := &Worker{
			ID:   i,
			Name: fmt.Sprintf("%s-worker-%v", p.Name, i),
		}
		p.Workers = append(p.Workers, worker)
	}
	//创建任务队列
	if p.QueueSize < 1 {
		p.QueueSize = 1
	}
	p.Queue = make(chan Job, p.QueueSize)

}

//start
func (p *Pool) Start() {
	for _, worker := range p.Workers {
		worker.Start(p.Queue)
	}
	fmt.Println("all workers started")

}

//stop
func (p *Pool) Stop() {
	close(p.Queue) //先关闭任务队列管道
	var wg sync.WaitGroup
	//因为当worker数量非常巨大时,可能会需要经常等待worker处理完毕,因此用协程来关闭
	for _, v := range p.Workers {
		wg.Add(1)
		go func(w *Worker) {
			defer wg.Done()
			w.stop()
		}(v)
	}
	wg.Wait()
	fmt.Println("all workers stopped")
}

func main() {
	f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	pool := &Pool{
		Name: "newPool",
		Size: 10000,

		QueueSize: 10000,
	}
	pool.Initialize()
	pool.Start()
	defer pool.Stop()

	//模拟一些Job(需实现start方法)
	for i := 0; i < 65535; i++ {
		job := &JobScan{
			//	IP:   "124.223.174.63",
			IP:   "127.0.0.1",
			Port: strconv.Itoa(i),
		}

		pool.Queue <- job

	}
	fmt.Println(runtime.NumGoroutine())

}

type JobScan struct {
	IP   string
	Port string
}

func (js *JobScan) Start(worker *Worker) error {
	conn, err := net.DialTimeout("tcp", js.IP+":"+js.Port, time.Second*2)
	if err != nil {
		return err
	}
	fmt.Printf("port:%s Open!\n", js.Port)
	conn.Close()
	return nil

}
