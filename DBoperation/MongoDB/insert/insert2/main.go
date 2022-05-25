package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func initDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("成功链接到数据库")
	return client
}

type HostInfo struct {
	IP        string
	Uptime    time.Time
	PortsInfo []Port
}

type Port struct {
	P    int
	Type string
	Info string
}

func main() {
	client := initDB()
	//先指定目标数据库和集合
	coll := client.Database("testdb").Collection("inventory")

	to := HostInfo{
		IP:     "127.0.0.1",
		Uptime: time.Now(),
		PortsInfo: []Port{
			{
				P:    22,
				Type: "ssh",
				Info: "rsa-256",
			},
		},
	}
	
	result, err := coll.InsertOne(context.TODO(), to)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("插入成功", result)

}
