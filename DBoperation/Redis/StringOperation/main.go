package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var ctx = context.Background() //v8版本必须带ctx

func initDB() *redis.Client {

	rdb := redis.NewClient(&redis.Options{ //定义链接
		Addr:     "192.168.67.130:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("链接成功")
	}
	return rdb
}

func main() {
	rdb := initDB() //创建链接

	//1.sget添加一条消息
	err := rdb.Set(ctx, "key3", "321", 0).Err() //最后一个参数表示过期时间
	if err != nil {
		log.Fatal(err)

	} else {
		fmt.Println("set插入单条成功!")
	}
	//2.mset批量插入
	err = rdb.MSet(ctx, "kk1", "你好", "kk2", "我是").Err()
	if err != nil {
		log.Fatal(err)

	} else {
		fmt.Println("set插入单条成功!")
	}

	//3.mget获取多条消息(Mget返回切片)
	val, err := rdb.MGet(ctx, "key1", "key2", "key3").Result()
	if err == redis.Nil {
		fmt.Println("你输入的key不存在值!")
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Mget多条消息", val)

	}

}
