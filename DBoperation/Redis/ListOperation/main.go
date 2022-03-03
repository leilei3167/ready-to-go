package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var ctx context.Context //声明全局ctx
var rdb *redis.Client

//1.初始化链接
func init() {
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.67.130:6379",
		Password: "",
		DB:       0,
	})

}

//2.Lpush/Rpush/lpop/rpop
func inAndout() {
	//左边推入数据
	for i := 0; i < 5; i++ {
		value := fmt.Sprintf("no.%vfrom left", i)
		err := rdb.LPush(ctx, "list", value).Err()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("左边推入数据成功")
		}
	}
	//右边推入数据
	/*	for i := 0; i < 5; i++ {
		value := fmt.Sprintf("no.%vfrom right", i)
		err := rdb.LPush(ctx, "list", value).Err()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("右边推入数据成功")
		}
	}*/
	//右边pop一个数据出来
	val3, err := rdb.RPop(ctx, "list").Result()
	if err == redis.Nil {
		fmt.Println("没有任何值了")
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("从右边取出的一个值为:", val3)
	}

}

func findlist() {
	//获取list的长度
	fmt.Println(rdb.LLen(ctx, "list").Result())
	//从0下标开始,15结束获取其list中的值
	val, err := rdb.LRange(ctx, "list", 0, -1).Result()
	if err == redis.Nil {
		fmt.Println("没有任何值")
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("遍历结果:", val)
	}
	//根据下标取值
	val2, err := rdb.LIndex(ctx, "list", 0).Result()
	if err == redis.Nil {
		fmt.Println("没有任何值")
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("下标0的值是:", val2)
	}

}

func main() {
	inAndout()
	time.Sleep(time.Second * 3)
	findlist()

}
