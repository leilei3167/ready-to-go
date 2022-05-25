package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	//查询一条
	coll := client.Database("result").Collection("IP")

	var result bson.M //map[string]interface{}
	//查找并解码
	err = coll.FindOne(context.TODO(), bson.D{{"isalive", true}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("没有相关的信息!")
		}
		panic(err)

	}
	//编码成json格式
	outpt, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(outpt))

	fmt.Println("--------------------查询多条-------------------")
	//后续复杂的过滤条件单独赋值
	//fliter := bson.D{{"isalive", false}}         //过滤
	fliter3 := bson.D{{"ip", "114.114.114.114"}} //过滤
	//fliter2 := bson.D{{"openports", 22}} //指定值而不是bson.A代表只要含有指定元素

	cursor, err := coll.Find(context.Background(), fliter3)
	if err != nil {

		log.Fatal(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil { //批量查询时使用All解析cursor
		panic(err)
	}
	if len(results) == 0 {
		log.Println("没有符合条件的结果!")
		return
	}
	for i, result := range results {
		outpt, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("结果%d:%s\n", i, string(outpt))
	}

}
