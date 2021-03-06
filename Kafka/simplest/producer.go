package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	//配置文件,还有非常多 非常复杂
	config := sarama.NewConfig()                              //默认配置
	config.Producer.RequiredAcks = sarama.WaitForAll          //发送数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //选出一个partition
	config.Producer.Return.Successes = true                   //成功交付的消息在success chan返回
	//Errors默认设置的是true
	//构造消息
	msg := &sarama.ProducerMessage{} //非常多的配置
	msg.Topic = "test_10"
	msg.Value = sarama.StringEncoder("go to scan!")

	//连接kafka
	client, err := sarama.NewSyncProducer([]string{"124.223.174.63:9092"}, config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("连接成功")
	defer client.Close() //延迟关闭释放资源

	for i := 0; i < 15; i++ {
		time.Sleep(time.Second * 2)
		pid, offset, err := client.SendMessage(msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("pid:%v offset:%v\n", pid, offset)
	}

}
