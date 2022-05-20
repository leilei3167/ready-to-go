package service

import (
	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

// GetDefaultConsumer 以分区创建消费者(弃用)
func GetDefaultConsumer(topic string) {
	config := sarama.NewConfig()

	consumer, err := sarama.NewConsumer([]string{"124.223.174.63:9092", "182.61.6.67:9092"}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	//查询该topic有多少个partition,每个分区对应一个消费者
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(partitions))

	for _, partitionID := range partitions {
		go consumeByPartition(consumer, topic, partitionID, &wg)
	}
	wg.Wait()
}

type Res struct {
	IP    string
	Ports []int
}

func consumeByPartition(consumer sarama.Consumer, topic string, partitionId int32, wg *sync.WaitGroup) {
	defer wg.Done()
	pConsumer, err := consumer.ConsumePartition(topic, partitionId, sarama.OffsetNewest)
	if err != nil {
		log.Fatal(err)
	}
	defer pConsumer.Close()

	for msg := range pConsumer.Messages() {
		m := Res{}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		_ = json.Unmarshal(msg.Value, &m)
		log.Printf("get msg from partition[%d]:%#v", partitionId, m)
		for _, port := range m.Ports {
			host := net.JoinHostPort(m.IP, strconv.Itoa(port))
			conn, err := net.DialTimeout("tcp", host, time.Second)
			if err != nil {
				log.Printf("HOST:%v连接失败:%v", host, err)
				continue
			}
			conn.Close()
			log.Printf("HOST:%v连接成功!!", host)
		}

	}
}
