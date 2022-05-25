package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

func main() {
	client := initDB()
	//先指定目标数据库和集合
	coll := client.Database("testdb").Collection("inventory")
	//1.插入一条
	result, err := coll.InsertOne(
		context.TODO(),
		bson.D{
			{"item", "leilei"},
			{"qty", 100},
			{"tags", bson.A{"cotton", "paper"}}, //bson.A代表的是一个切片
			{"size", bson.D{ //bson.D代表的是一个结构体
				{"h", 28},
				{"w", 35.5},
				{"uom", "cm"},
			}},
		})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("插入成功ID:", result)
	//1.1查询数据库中的结构
	res, err := coll.Find(context.Background(), bson.D{{"qty", 100}})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("查询成功:", res)

	//2.批量插入
	result2, err := coll.InsertMany(
		context.TODO(),
		[]interface{}{
			bson.D{
				{"item", "journal"},
				{"qty", int32(25)},
				{"tags", bson.A{"blank", "red"}},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
			},
			bson.D{
				{"item", "mat"},
				{"qty", int32(25)},
				{"tags", bson.A{"gray"}},
				{"size", bson.D{
					{"h", 27.9},
					{"w", 35.5},
					{"uom", "cm"},
				}},
			},
			bson.D{
				{"item", "mousepad"},
				{"qty", 25},
				{"tags", bson.A{"gel", "blue"}},
				{"size", bson.D{
					{"h", 19},
					{"w", 22.85},
					{"uom", "cm"},
				}},
			},
		})
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range result2.InsertedIDs {
		log.Println("批量插入成功:", d)
	}

}
