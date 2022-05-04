package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

//简单consumer
func main() {
	consumer, err := sarama.NewConsumer([]string{"124.223.174.63:9092", "182.61.6.67:9092"}, nil) //默认配置
	if err != nil {
		log.Fatal("连接broker失败:", err)
	}
	partitionList, err := consumer.Partitions("test_10")

	if err != nil {
		log.Fatal("fail to get partitionList:", err)
	}
	fmt.Println("partitionList:", partitionList)

	//遍历所有的分区
	for p := range partitionList {
		//针对每个分区创建分区消费者
		pc, err := consumer.ConsumePartition("quickstart-events", int32(p), sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("failed to start consumer for partition %d,err:%v\n", p, err)
		}
		defer pc.Close()

		//异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("分区[%d] offset:%d Got: Key:%s Value:%s\n", msg.Partition,
					msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}
	time.Sleep(time.Minute)
}
