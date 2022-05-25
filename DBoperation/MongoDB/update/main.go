package main

import (
	"context"
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
	//修改一条
	coll := client.Database("result").Collection("IP")

	//过滤查询
	filter := bson.D{{"isalive", true}}
	update := bson.D{{"$set", bson.D{{"isalive", false}}}}

	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	log.Println("修改成功,id:", result)
	fmt.Println("-------------------------更新多个-----------------")
	results, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	log.Println("修改成功,id:", results)
}
