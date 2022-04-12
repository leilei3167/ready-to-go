package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"net"
	"time"
)

type Worker struct {
	Num int
	//	Err chan error
}

//Worker method
//由workerpool直接控制,workerpool传入任务队列
func (w *Worker) Start(ctx context.Context, tasks chan Task) {

	//阻塞获取任务,或监听停止信号
	for {
		select {
		case t := <-tasks:
			t.Start(w)

		case <-ctx.Done():
			log.Println("Cancled")
			return
		}
	}
}

type WorkerPool struct {
	Type      string
	Size      int //worker的数量(协程)
	Workers   []*Worker
	QueueSize int //任务队列的大小
	TaskQueue chan Task
	Ctx       context.Context //控制worker关闭
}

//返回一个Pool和一个控制Worker关闭的Cancle
func NewWorkerPool(typeName string, workerSize, queueSize int) (*WorkerPool, context.CancelFunc) {
	//TODO:server和cli RPC交互的模型,可加入字段type,针对不同的type创建工作池(size由自己写死还是交由指定?)
	/*
		switch task.Type
			case:...根据不同的type创建不同的pool
	*/
	ctx, cancle := context.WithCancel(context.Background())

	return &WorkerPool{
		Type:      typeName,
		Size:      workerSize,
		QueueSize: queueSize,
		Ctx:       ctx,
	}, cancle

}

//创建工作池后,根据参数初始化worker
func (p *WorkerPool) Initlize() {
	//	errChan := make(chan error)
	for i := 0; i < p.Size; i++ {
		worker := &Worker{
			Num: i,
			//Err: errChan,
		}
		p.Workers = append(p.Workers, worker)
	}
	p.TaskQueue = make(chan Task, p.QueueSize)
}

func (p *WorkerPool) Start() {
	for _, worker := range p.Workers {
		go worker.Start(p.Ctx, p.TaskQueue)
	}
}

//WorkerPool  method

//task接口,便于后期拓展task的类别
type Task interface {
	Start(*Worker) error
}

//扫描任务相关的
type Scan struct {
	IP   string
	Port string
}

func (s *Scan) Start(w *Worker) error {

	conn, err := net.DialTimeout("tcp", s.IP+":"+s.Port, time.Second*2)
	if err != nil {
		//	w.Err <- err
		return nil
	}
	fmt.Printf("port:%s Open!\n", s.Port)
	conn.Close()
	return nil
}

func main() {
	p, cancle := NewWorkerPool("scanport", 20000, 20000)

	p.Initlize()
	p.Start()
	//需扫描一个ip
	for i := 0; i < 65535; i++ {
		t := &Scan{
			IP:   "127.0.0.1",
			Port: strconv.Itoa(i),
		}
		p.TaskQueue <- t
	}
	time.Sleep(time.Second * 10)
	cancle()
}
