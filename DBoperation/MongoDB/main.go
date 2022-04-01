package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Student struct {
	Name string
	Age  int
}

func main() {
	/*连接到mongo*/
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

	/*MongoDB中的JSON文档存储在名为BSON(二进制编码的JSON)的二进制表示中。与其他将JSON数据存储为简单字符串和数字的数据库不同，
	BSON编码扩展了JSON表示，使其包含额外的类型，如int、long、date、浮点数和decimal128。
	这使得应用程序更容易可靠地处理、排序和比较数据。
		连接MongoDB的Go驱动程序中有两大类型表示BSON数据：D和Raw
	类型D家族被用来简洁地构建使用本地Go类型的BSON对象。这对于构造传递给MongoDB的命令特别有用。D家族包括四类:

	D：对应切片,这种类型应该在顺序重要的情况下使用，比如MongoDB命令。
	M：对应map
	A：一个BSON数组。
	E：D里面的一个元素。
	*/
	/*bson.D{{
		"name", bson.D{{
			"$in", bson.A{"张三", "李四"},
		}},
	}}*/
	//Raw类型家族用于验证字节切片。你还可以使用Lookup()从原始类型检索单个元素。
	//当要查找BSON字节而不将其解编为另一种类型时，此类型最有用。

	/*CRUD*/
	//需先指定要操作的数据集
	collection := client.Database("test").Collection("student")
	s1 := Student{"小红", 12}
	s2 := Student{"小兰", 10}
	s3 := Student{"小黄", 11}
	//插入:InsertOne 和InsertMany
	one, err := collection.InsertOne(context.TODO(), s1)
	if err != nil {
		return
	}
	fmt.Println("insertOne成功:", one, "ID", one.InsertedID)

	many, err := collection.InsertMany(context.TODO(), []interface{}{s2, s3})
	if err != nil {
		log.Fatal("批量插入失败!", err)
	}
	fmt.Println("insertMany成功!", many, "ID", many.InsertedIDs)

	//更新Update
	//它需要一个筛选器文档来匹配数据库中的文档，并需要一个更新文档来描述更新操作。你可以使用bson.D类型来构建筛选文档和更新文档
	//TODO:创建filter这步需要再看看
	filter := bson.D{{"name", "小兰"}}
	updata := bson.D{{
		"$inc", bson.D{
			{"age", 1},
		},
	}}

	updateRes, err := collection.UpdateOne(context.TODO(), filter, updata)
	if err != nil {
		log.Fatal("更新失败", err)
	}
	fmt.Println("更新成功", updateRes)

	//查找FindOne Find
	//查找同样需要一个filter文档，以及一个指向可以将结果解码为其值的指针
	var res Student                                               //容纳结果
	err = collection.FindOne(context.TODO(), filter).Decode(&res) //必须解码
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("res的结果:", res)
	//TODO:查询多个值

	/*	//删除
		//单个
		deleteRes, err := collection.DeleteOne(context.TODO(), filter)
		fmt.Println("删除成功!", deleteRes)
		//删除所有(注意filter是bson.D{}
		deleteRes2, err := collection.DeleteMany(context.TODO(), bson.D{})
		fmt.Println("删除所有!", deleteRes2)*/
}
