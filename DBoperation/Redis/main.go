package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func example() {
	//创建链接
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.67.130:6379",
		Password: "",
		DB:       0,
	})
	//插入一个string
	err := rdb.Set(ctx, "key1", "v100", 0).Err()
	if err != nil {
		panic(err)
	}
	//获取一个string
	val, err := rdb.Get(ctx, "key1").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key1", val)
	//获取一个不存在的key
	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2不存在!!")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

}

func main() {
	example()

}