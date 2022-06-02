package service

import (
	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/sync/semaphore"
	"mapping/internal/pkg/db"
	"sync"
)

// Consumer 实现ConsumerGroupHandler接口,才能作为ConsumerGroup使用
type Consumer struct {
	Name       string
	Count      int64
	TaskChan   chan GotTask
	ResultChan chan Result
	Sem        *semaphore.Weighted //控制同时具有的连接数
	Client     sarama.ConsumerGroup
}

// Result 以IP为单位,构建结果插入Mongo,以日期分类?
type Result struct {
	db.ScanResult
	mu *sync.Mutex
}

func InitConsumer(gsize int, addrs []string, groupID string) (*Consumer, error) {
	var a Consumer
	a.Name = "test"
	a.Count = 0
	a.TaskChan = make(chan GotTask, 1000)
	a.ResultChan = make(chan Result, 100000)
	a.Sem = semaphore.NewWeighted(int64(gsize))
	//开启消费者组
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest //从最旧的,为被标记被消费的offset开始消费

	c, err := sarama.NewConsumerGroup(addrs, groupID, config) //创建新的消费组
	if err != nil {
		return nil, err
	}
	a.Client = c

	return &a, nil
}

type GotTask struct {
	IP    string `json:"IP,omitempty"`
	Ports []int  `json:"ports,omitempty"`
}

//实现sarama的ConsumerGroupHandler接口,才能创建成消费者组,核心的消费逻辑在ConsumeClaim中

// Setup 执行在 获得新 session 后 的第一步, 在 ConsumeClaim() 之前
func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup 执行在 session 结束前, 当所有 ConsumeClaim goroutines 都退出时
func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 具体的消费逻辑
func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		//构建任务 发往任务队列
		m := GotTask{}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		_ = json.Unmarshal(msg.Value, &m)
		c.TaskChan <- m
		//标记为已消费
		sess.MarkMessage(msg, "")

	}
	return nil
}
