package service

import "github.com/Shopify/sarama"

// GetDefaultProducer 创建一个默认配置的生产者,暂时硬编码
func GetDefaultProducer() (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner //轮询发送至分区
	config.Producer.Compression = sarama.CompressionSnappy
	//config.Producer.Flush.Frequency = time.Millisecond * 10
	producer, err := sarama.NewAsyncProducer([]string{"124.223.174.63:9092", "182.61.6.67:9092"}, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
