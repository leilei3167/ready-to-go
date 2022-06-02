package node

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mapping/internal/node/service"
	"mapping/internal/pkg/db"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"
)

type Node struct {
	//消费者组控制
	C *service.Consumer
	//数据库连接控制
	mgo *mongo.Client
}

func NewNode(mgourl string, size int, brokers []string, groupID string) *Node {
	n := new(Node)
	m, err := db.NewMongoDB(mgourl)
	if err != nil {
		log.Fatal("连接mongo失败:", err)
	}
	n.mgo = m
	n.C, err = service.InitConsumer(size, brokers, groupID)
	if err != nil {
		log.Println("创建消费者组失败:", err)
	}
	return n
}

func (n *Node) Close() {
	//关闭消费者组
	err := n.C.Client.Close()
	if err != nil {
		log.Println("关闭消费者组错误:", err)
	}

	//关闭所有的G
	close(n.C.TaskChan)
	time.Sleep(time.Second)
	close(n.C.ResultChan) //有一定几率触发data race

	//断开mongo连接
	err = n.mgo.Disconnect(context.TODO())
	if err != nil {
		log.Println("关闭mog连接出错:", err)
	}
	log.Println("关闭服务成功")
}

func (n *Node) Run() {
	defer n.Close()
	for i := 0; i < 100; i++ {
		go n.C.ToScan()
	}

	go n.GetResult() //获取扫描结果,存入mongo

	//调试用,后续可增加优雅退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	//开启消费
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := n.C.Client.Consume(ctx, []string{"test_10"}, n.C)
			if err != nil {
				return
			}
			if ctx.Err() != nil {
				return
			}
		}

	}()

	select {
	case s := <-c:
		log.Printf("收到信号:%v退出...", s)
		return
	}

	wg.Wait()

}

func (n *Node) GetResult() { //可考虑将数据专门交由某个部分进行存储,node只负责放入消息队列
	var total int
	var alive int
	mgo := n.mgo.Database("result").Collection("IP")
	for result := range n.C.ResultChan {
		total++

		if len(result.OpenPorts) > 0 { //只存入存活的IP
			result.IsAlive = true
			result.Uptime = time.Now()
			sort.Ints(result.OpenPorts) //排序

			//根据IP查询,如果没有 则新插入
			fliter := bson.D{{"ip", result.IP}}
			update := bson.D{{"$set", result.ScanResult}}

			options := options.Update().SetUpsert(true)
			_, err := mgo.UpdateOne(context.TODO(), fliter, update, options)
			if err != nil {
				log.Println("插入数据出错:", err)
			}
			alive++
		}
	}

	log.Printf("GetResult退出,已完成扫描IP:%d个,存活IP:%d个", total, alive)
}
