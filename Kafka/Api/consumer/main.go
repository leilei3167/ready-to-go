package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	//SinglePartition("test_10") //单分区消费者
	//Partition("test_10") //offsetManager

}

//单分区消费
func SinglePartition(topic string) {
	config := sarama.NewConfig()
	//不配置 采用默认配置
	consumer, err := sarama.NewConsumer([]string{"124.223.174.63:9092"}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	p, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("消费错误:", err)
	}

	defer p.Close()

	for msg := range p.Messages() {
		log.Printf("[Consumer] pid:%d offset:%d value:%s", msg.Partition, msg.Offset, msg.Value)
	}

}

//同时消费多个分区
func Partition(topic string) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{"124.223.174.63:9092"}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close() //必须要在子分区消费者全部关闭之后才能执行,这里利用defer先进后出的原理

	//多分区消费先查询分区列表
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Println("get partitions err:", err)
	}
	//之后每一个分区开一个消费者
	var wg sync.WaitGroup
	wg.Add(len(partitions))

	for _, p := range partitions {
		go func(p int32) {
			defer wg.Done()
			partitionCosumer, err := consumer.ConsumePartition(topic, p, sarama.OffsetNewest)
			if err != nil {
				log.Fatal(err)
			}
			defer partitionCosumer.Close() //子分区消费者必须关闭
			for message := range partitionCosumer.Messages() {
				log.Printf("[Consumer] partitionid: %d; offset:%d, value: %s\n", message.Partition, message.Offset, string(message.Value))
			}
		}(p)
	}
	wg.Wait()
}

/*
TODO:测试无效,不知什么原因
*/
//offsetManager,在使用consumer开始消费时,第三个参数之前一直指定的是oldest或者newest,而如果我们像
//让这个consumer接着上次的offset开始该怎么办?这就需要创建offsetmanager
func OffsetManager(topic string) {
	config := sarama.NewConfig()
	// 配置开启自动提交 offset，这样 samara 库会定时帮我们把最新的 offset 信息提交给 kafka
	config.Consumer.Offsets.AutoCommit.Enable = true              // 开启自动 commit offset
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 自动 commit时间间隔
	//--------实际默认配置是配置了以上的设置项的!

	//1.此处是创建client而不是直接创建Consumer
	client, err := sarama.NewClient([]string{"124.223.174.63:9092"}, config)
	if err != nil {
		log.Fatal("NewClient err: ", err)
	}
	defer client.Close()
	//2.创建偏移量管理器
	offsetManager, err := sarama.NewOffsetManagerFromClient("myGroupID", client)
	if err != nil {
		log.Println("NewOffsetManagerFromClient err:", err)
	}
	defer offsetManager.Close()

	//每个分区的offset也是分别管理的,此处只是管理分区0
	partitionOffsetManager, err := offsetManager.ManagePartition(topic, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer partitionOffsetManager.Close()
	defer offsetManager.Commit() //退出前再提交一次,防止自动提交的间隔之间的信息丢失

	consumer, err := sarama.NewConsumerFromClient(client) //创建消费者
	if err != nil {
		log.Println("NewConsumerFromClient err:", err)
	}
	//根据kafka中的记录的上次消费的offset开始+1的位置接着消费
	nextOffset, _ := partitionOffsetManager.NextOffset()
	fmt.Println("nextOffset:", nextOffset)
	pc, err := consumer.ConsumePartition(topic, 0, nextOffset)
	if err != nil {
		log.Println("ConsumePartition err:", err)
	}
	defer pc.Close()

	for message := range pc.Messages() {
		value := string(message.Value)
		log.Printf("[Consumer] partitionid: %d; offset:%d, value: %s\n", message.Partition, message.Offset, value)
		//每次消费完毕之后都更新一次offset,此处是更新的本机内存中的值,只有在commit之后才会传到broker中
		partitionOffsetManager.MarkOffset(message.Offset+1, "modified metadata")
	}

}

//----------------------------------------------------------------------------------
//以上的方式是单消费者对多分区,需要手动的去开启goroutine去消费,并且还需要维护offset 维护麻烦
//sarama支持consumer group,一个消费者组可以有多个消费者,kafka会以分区为单位将消息分给各个消费者
//每条消息将只会被消费者组中的一个消费者消费!

//重点是实现sarama.ConsumerGroup接口,作为自定义的消费者组
type MyConsumerGroupHandler struct {
	name  string
	count int64
}

// Setup 执行在 获得新 session 后 的第一步, 在 ConsumeClaim() 之前
func (MyConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup 执行在 session 结束前, 当所有 ConsumeClaim goroutines 都退出时
func (MyConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

//ConsumeClaim 核心消费的逻辑
func (h MyConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		fmt.Printf("[consumer] name:%s topic:%q partition:%d offset:%d\n", h.name, msg.Topic, msg.Partition, msg.Offset)
		// 标记消息已被消费 内部会更新 consumer offset
		sess.MarkMessage(msg, "")
		h.count++
		if h.count%10000 == 0 { //每消费10000报一次
			fmt.Printf("name:%s 消费数:%v\n", h.name, h.count)
		}

	}
	return nil
}

func ConsumerGroup(topic, group, name string) {
	//默认配置
	config := sarama.NewConfig()
	//消费者组必须要有ctx
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//创建消费者组
	cg, err := sarama.NewConsumerGroup([]string{"124.223.174.63"}, group, config)
	if err != nil {
		log.Fatal(err)
	}
	defer cg.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handler := MyConsumerGroupHandler{name: name} //实现了消费者组接口的实例
		for {
			fmt.Println("running!", name)
			/*
			   应该在一个无限循环中不停的调用Consume()
			   因为每次Rebalance后需要再次执行Consume()来恢复连接


			*/
			err = cg.Consume(ctx, []string{topic}, handler)
			if err != nil {
				log.Println("Consumer Err:", err)
			}
			//如果ctx被cancle了 则退出
			if ctx.Err() != nil {
				return
			}
		}

	}()
	wg.Wait()

}
