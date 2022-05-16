package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

//构建消费者 获取消息任务并执行扫描

func main() {
	config := sarama.NewConfig()
	client, err := sarama.NewConsumer([]string{"124.223.174.63:9092", "182.61.6.67:9092"}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	partitionConsumer, err := client.ConsumePartition("test_10", 1, sarama.OffsetNewest)
	if err != nil {
		log.Fatal(err)
	}
	defer partitionConsumer.Close()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	var wg sync.WaitGroup
	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var res Task
			err := json.Unmarshal(msg.Value, &res)
			if err != nil {
				log.Println(err)
			}
			log.Printf("Consumed message offset %d\n,Value:%v", msg.Offset, res)
			wg.Add(1)
			go func() {
				defer wg.Done()
				for _, port := range res.Ports {
					conn, err := net.DialTimeout("tcp", net.JoinHostPort(res.Ip, strconv.Itoa(port)), time.Second*2)
					if err != nil {
						log.Printf("%v:%v is CLOSED!", res.Ip, port)
						continue
					}
					conn.Close()
					log.Printf("%v:%v is OPEN!", res.Ip, port)
				}

			}()
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}
	wg.Wait()
	log.Printf("Consumed: %d\n", consumed)
}

type Task struct {
	Ip    string `json:"ip,omitempty"`
	Ports []int  `json:"ports,omitempty"`

	encoded []byte
	err     error
}
