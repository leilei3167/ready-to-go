package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

/*
此案例演示两种生产者的用法:
1.SyncProducer代表同步生产者(实际较少用),程序将阻塞直到其返回结果
2.AsyncProducer则是异步生产者,初始化后将提供几个channel,将需要生产的信息放入channel即可,无需等待返回

重点:
1.根据实际的需求来创建config,如压缩算法,是否需要tls等等
2.通过NewAsyncProducer或NewSyncProducer来根据配置和brokerlist来生成生产者,他们是并发安全的
3.学习其并发控制和关闭逻辑



*/
var (
	addr    = flag.String("addr", ":8080", "The address to bind to")
	brokers = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")
	verbose = flag.Bool("verbose", false, "Turn on Sarama logging")

	//验证相关
	certFile  = flag.String("certificate", "", "The optional certificate file for client authentication")
	keyFile   = flag.String("key", "", "The optional key file for client authentication")
	caFile    = flag.String("ca", "", "The optional certificate authority file for TLS client authentication")
	verifySsl = flag.Bool("verify", false, "Optional verify ssl certificates chain")
)

func main() {
	flag.Parse()
	if *verbose { //开启sarama日志记录的话
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags) //创建一个logger
	}

	if *brokers == "" { //必须指定至少一个kafka节点
		flag.PrintDefaults()
		os.Exit(1)
	}
	brokerList := strings.Split(*brokers, ",")
	log.Printf("Kafka brokers: %s", strings.Join(brokerList, ", "))

	server := &Server{ //初始化2个消息生产者
		DataCollector:     newDataCollector(brokerList),     //记录访问者的查询参数
		AccessLogProducer: newAccessLogProducer(brokerList), //记录使用日志
	}

	//优雅退出
	defer func() {
		if err := server.Close(); err != nil {
			log.Println("Failed to close server", err)
		}
	}()

	//运行
	log.Fatal(server.Run(*addr))
}

//分别时同步生产者和异步生产者,实现生产消息的核心!
type Server struct {
	DataCollector     sarama.SyncProducer
	AccessLogProducer sarama.AsyncProducer
}

//关闭逻辑,两个生产者都需要关闭
func (s *Server) Close() error {
	if err := s.DataCollector.Close(); err != nil {
		log.Println("Failed to shut down data collector cleanly", err)
	}
	if err := s.AccessLogProducer.Close(); err != nil {
		log.Println("Failed to shut down access log producer cleanly", err)
	}
	return nil
}

//此处开始为重点逻辑,演示了两种producer的配合使用以及中间件的构建方式:log作为外层的中间件
func (s *Server) Handler() http.Handler {
	return s.withAccessLog(s.collectQueryStringData())
}

func (s *Server) Run(addr string) error {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: s.Handler(), //直接接管默认的Mux
	}
	log.Printf("Listening for requests on %s...\n", addr)
	return httpServer.ListenAndServe() //开启监听
}

func (s *Server) collectQueryStringData() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" { //只接受 / 的请求
			http.NotFound(w, r)
			return
		}
		//不设置发送的key,意味着消息将随机发送到不同的partition(轮询),此处注意和异步生产者的发送消息的不同
		//会阻塞直到返回p,o,err
		partition, offset, err := s.DataCollector.SendMessage(&sarama.ProducerMessage{
			Topic: "query",
			Value: sarama.StringEncoder(r.URL.RawQuery), //值为查询字段
		})
		if err != nil { //失败
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to store your data: %s", err)
		} else { //成功返回一个组合的唯一标志
			//topic,partition,offset可以作为一个消息在kafka中的唯一标识
			fmt.Fprintf(w, "Your data is stored with unique identifier query/%d/%d", partition, offset)
		}
	})
}

//定义日志所需要的格式,以及编码的方法
type accessLogEntry struct {
	Method       string  `json:"method"`
	Host         string  `json:"host"`
	Path         string  `json:"path"`
	IP           string  `json:"ip"`
	ResponseTime float64 `json:"response_time"`

	encoded []byte
	err     error
}

func (ale *accessLogEntry) ensureEncoded() {
	if ale.encoded == nil && ale.err == nil { //转码过的不执行任何操作
		ale.encoded, ale.err = json.Marshal(ale) //把自己转码成字节切片后存入字段,必须确保是之前没有转码过的
	}
}

//实现Encode接口必须的两个方法,有此方法才能被转换成压缩数据
func (ale *accessLogEntry) Length() int {
	ale.ensureEncoded()
	return len(ale.encoded)
}

func (ale *accessLogEntry) Encode() ([]byte, error) {
	ale.ensureEncoded()
	return ale.encoded, ale.err
}

//中间件外层
func (s *Server) withAccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r) //执行内层的处理器,此处会阻塞直到其执行完毕返回,无论成功与否
		//以下为生产一条日志记录
		entry := &accessLogEntry{
			Method:       r.Method,
			Host:         r.Host,
			Path:         r.RequestURI, //和URL.Path有什么区别?
			IP:           r.RemoteAddr,
			ResponseTime: float64(time.Since(started)) / float64(time.Second), //float64
		}
		//此时使用客户端的IP来作为消息的Key,给定key的情况下,Kafka将会根据key的hash值来选择partition,意味着同一ip的访问记录都会在同一个分区上
		//直接将要发送的消息放入即可
		s.AccessLogProducer.Input() <- &sarama.ProducerMessage{Topic: "access_log",
			Key:   sarama.StringEncoder(r.RemoteAddr),
			Value: entry,
		}
	})

}

//创建配置文件 并生成一个同步的Producer
func newDataCollector(brokerlist []string) sarama.SyncProducer {
	//这里寻求强一致性(不改变flush设置),sarama会尽可能的快速生产消息以降低延迟
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll //需要等到副本全部回复收到消息
	config.Producer.Retry.Max = 5                    //生产消息的最大重试次数
	config.Producer.Return.Successes = true          //设置为true,在生产消息成功时将向SuccessChan中发送一个消息
	//可选的tls配置
	tlsConfig := createTlsConfiguration()
	if tlsConfig != nil {
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Enable = true
	}
	producer, err := sarama.NewSyncProducer(brokerlist, config)
	if err != nil {
		log.Fatalln("开启Producer失败:", err)
	}
	return producer

}

//生成一个异步的日志producer
func newAccessLogProducer(brokerlist []string) sarama.AsyncProducer {
	//对于日志的,我更关注高吞吐量,设置压缩的,可以减少网络IO造成的延迟
	//kafka会等待消息达到一定大小(默认16k)或者达到最大等待时长时才会发送消息
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll         //等待所有副本回复
	config.Producer.Compression = sarama.CompressionSnappy   //snappy压缩算法
	config.Producer.Flush.Frequency = time.Millisecond * 500 //每500ms发送一批次
	//可选的tls配置
	tlsConfig := createTlsConfiguration()
	if tlsConfig != nil {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}

	producer, err := sarama.NewAsyncProducer(brokerlist, config)
	if err != nil {
		log.Fatalln("Failed to start Log producer:", err)
	}
	//如果生产消息失败,那么就打印到标准输出上,注意只有最大重试次数到了才会执行
	go func() {
		for err := range producer.Errors() { //在执行close的函数后 这些通道都会被关闭
			log.Println("Fail to Write access log:", err)
		}
	}()

	return producer
}

func createTlsConfiguration() (t *tls.Config) {
	if *certFile != "" && *keyFile != "" && *caFile != "" {
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			log.Fatal(err)
		}

		caCert, err := os.ReadFile(*caFile)
		if err != nil {
			log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: *verifySsl,
		}
	}
	// will be nil by default if nothing is provided
	return t
}
