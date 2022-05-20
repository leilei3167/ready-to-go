package service

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

// Consumer 实现ConsumerGroupHandler接口,才能作为ConsumerGroup使用
type Consumer struct {
	Name       string
	Count      int64
	TaskChan   chan GotTask
	ResultChan chan Result   //构建一个结果结构体
	WorkerPool chan struct{} //并发限制
}

// Setup 执行在 获得新 session 后 的第一步, 在 ConsumeClaim() 之前
func (*Consumer) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup 执行在 session 结束前, 当所有 ConsumeClaim goroutines 都退出时
func (*Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim 具体的消费逻辑
func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		//构建任务 发往任务队列
		m := GotTask{}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		_ = json.Unmarshal(msg.Value, &m)
		//	log.Printf("get msg from partition[%d]:%v", msg.Partition, m)
		c.TaskChan <- m
		//标记为已消费
		sess.MarkMessage(msg, "")
		c.Count++
		if c.Count%100 == 0 {
			fmt.Printf("Name:%s 消费数:%v\n", c.Name, c.Count)
		}
	}
	return nil
}

func (c *Consumer) ToScan() {
	for task := range c.TaskChan {
		var r Result
		var wg sync.WaitGroup
		r.IP = task.IP
		r.Mu = new(sync.Mutex)
		wg.Add(len(task.Ports))
		for _, p := range task.Ports {
			c.WorkerPool <- struct{}{}
			go func(h string, p int) {
				defer wg.Done()
				host := net.JoinHostPort(r.IP, strconv.Itoa(p))
				//	log.Println(host)
				conn, err := net.DialTimeout("tcp", host, time.Second)
				if err != nil {
					//	log.Println(err)
					return
				}
				log.Printf("%v is open", host)
				conn.Close()
				r.Mu.Lock()
				r.OpenPorts = append(r.OpenPorts, p)
				r.Mu.Unlock()
				<-c.WorkerPool
			}(r.IP, p)
		}
		/*if len(r.OpenPorts) > 0 {
			c.ResultChan <- r
		}*/
		wg.Wait()
		c.ResultChan <- r
	}
}

func (c *Consumer) GetResult() {
	for result := range c.ResultChan {
		if len(result.OpenPorts) > 0 {
			log.Printf("%s:[%v] is open!", result.IP, result.OpenPorts)
			continue
		}
		log.Printf("%v没有打开的端口!", result.IP)

	}
}

type GotTask struct {
	IP    string `json:"IP,omitempty"`
	Ports []int  `json:"ports,omitempty"`
}

type Result struct { //一个IP一个结果
	IP string
	//结果
	IsAlive   bool //TODO:扫描开始前只对存活的ip扫描
	OpenPorts []int
	Mu        *sync.Mutex //多端口并发扫描,一个IP存活时对端口进行并发扫描
}

func InitConsumer(gsize int) *Consumer {
	var a Consumer
	a.Name = "test"
	a.Count = 0
	a.TaskChan = make(chan GotTask, 1000)
	a.ResultChan = make(chan Result, 1000)
	a.WorkerPool = make(chan struct{}, gsize)
	return &a
}

func GetConsumerGroup(gsize int) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest //从中断的offset开始消费
	//TODO:消费者组可以有更多配置
	//创建一个consumerGroup
	consumer := InitConsumer(gsize)
	for i := 0; i < 10; i++ {
		go consumer.ToScan()
	}
	go consumer.GetResult()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := sarama.NewConsumerGroup([]string{"124.223.174.63:9092", "182.61.6.67:9092"}, "testGroup", config) //创建新的消费组
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}
	defer client.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := client.Consume(ctx, []string{"test_10"}, consumer)
			if err != nil {
				return
			}
			if ctx.Err() != nil {
				return
			}
		}

	}()
	wg.Wait()
}
