package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

func NewMongoDB(url string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(url)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3) //数据库连接超时
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("成功链接到数据库")
	return client, nil
}

//全局数据库连接
var MgoClient *mongo.Client
var one sync.Once

func InitDB(url string) {
	var err error
	MgoClient, err = NewMongoDB(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("连接数据库成功")
}
