package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

var (
	Brokers = []string{"124.223.174.63:9092", "182.61.6.67:9092"}
	IPs     = []string{"8.8.8.8", "114.114.114.114"}
)

//向kafka中生产扫描任务
func main() {
	//userSyncProduc() //同步生产者,轮询
	//useAsyncProducerSelect()
	useAsyncProducerFor()
}

func useAsyncProducerSelect() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewManualPartitioner //指定分区选择模式
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Return.Successes = false
	config.Producer.Flush.Frequency = time.Millisecond * 10

	producer, err := sarama.NewAsyncProducer(Brokers, config)
	if err != nil {
		panic(err)
	}
	defer func() { //关闭
		if err := producer.Close(); err != nil {
			log.Println(err)
		}
	}()
	//监听一个系统退出信号
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	var (
		hasSent int

		errCount int
	)

LABLE:
	for {
		select {
		case producer.Input() <- &sarama.ProducerMessage{Topic: "test_10", Value: &Task{
			Ip:    "No." + strconv.Itoa(hasSent),
			Ports: TOP100,
		},
			Partition: 9,
		}:
			time.Sleep(time.Second)
			hasSent++
		case err := <-producer.Errors():
			log.Println("An error happened:", err)
			errCount++
		case <-signals:
			log.Printf("退出...总共已发送:%d,错误%d", hasSent, errCount)
			break LABLE
		}

	}

}

func useAsyncProducerFor() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner //指定分区选择模式
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Return.Successes = true
	config.Producer.Flush.Frequency = time.Millisecond * 10

	producer, err := sarama.NewAsyncProducer(Brokers, config)
	if err != nil {
		panic(err)
	}
	defer func() { //关闭
		if err := producer.Close(); err != nil {
			log.Println(err)
		}
	}()
	//监听一个系统退出信号
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	var (
		hasSent  int
		wg       sync.WaitGroup
		errCount int
	)
	//接受结果
	wg.Add(2)
	go func() {
		defer wg.Done()
		for sucMsg := range producer.Successes() {
			fmt.Println(sucMsg)
			hasSent++
		}
	}()
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			fmt.Println(err)
			errCount++
		}
	}()
	//生产
	go func() {
		for _, ip := range IPs {
			task := &Task{
				Ip:    ip,
				Ports: TOP100,
			}
			producer.Input() <- &sarama.ProducerMessage{
				Topic: "test_10",
				Value: task,
			}
		}
	}()
	select {
	case <-signals:
		log.Printf("收到退出...已发送%d,错误%v", hasSent, errCount)
		return
	}

}

func userSyncProduc() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner //指定轮询选择器,默认的话会是根据key(如有)进行hash
	config.Producer.Return.Successes = true                       //默认是关闭的,而同步生产者必须打开
	config.Producer.Compression = sarama.CompressionSnappy        //使用压缩
	config.Producer.Flush.Frequency = 1000 * time.Millisecond     //每多少ms发送一批次,合理的值可以加大吞吐量
	synroducer, err := sarama.NewSyncProducer(Brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	defer synroducer.Close()
	for _, ip := range IPs {
		t := &Task{
			Ip:    ip,
			Ports: TOP100,
		}
		msg := &sarama.ProducerMessage{
			Topic: "test_10",
			Value: t,
		}
		message, i, err := synroducer.SendMessage(msg)
		if err != nil {
			log.Println("Send msg Err:", err)
		}
		log.Printf("partition:%v offset:%v", message, i)
	}
}

//定义任务所需的字段,要想传递结构体,必须实现Encoder接口
type Task struct {
	Ip    string `json:"ip,omitempty"`
	Ports []int  `json:"ports,omitempty"`

	encoded []byte
	err     error
}

func (t *Task) Encode() ([]byte, error) {
	t.ensureEncode()
	return t.encoded, t.err
}
func (t *Task) Length() int {
	return len(t.encoded)
}
func (t *Task) ensureEncode() {
	if t.encoded == nil && t.err == nil {
		t.encoded, t.err = json.Marshal(t)
	}
}

var TOP100 = []int{22, 80, 2181, 9092, 1900, 2000, 2049, 2121, 2717, 3000, 3128, 3306, 3389, 3986, 4899, 5000, 5009}

/*var TOP100 = []int{7, 9, 13, 21, 22, 25, 37, 53, 79, 80, 88, 106, 110, 113, 119, 135,
139, 143, 179, 199, 389, 427, 443, 444, 465, 513, 514, 543, 548, 554, 587, 631,
646, 873, 990, 993, 995, 1025, 1026, 1027, 1028, 1110, 1433, 1720, 1723, 1755,
1900, 2000, 2049, 2121, 2717, 3000, 3128, 3306, 3389, 3986, 4899, 5000, 5009,
5051, 5060, 5101, 5190, 5357, 5432, 5631, 5666, 5800, 5900, 6000, 6646, 7070, 8000,
8008, 8080, 8443, 8888, 9100, 9999, 32768, 49152, 49153, 49154, 49155, 49156}
)*/
