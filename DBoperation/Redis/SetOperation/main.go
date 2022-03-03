package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
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
func add() {
	//向set中添加多个值
	err := rdb.SAdd(ctx, "set", 1, "nihao", 4, "woshi leilei").Err()
	if err != nil {
		log.Fatalln("sadd失败:", err)
	} else {
		fmt.Println("sadd成功!")
	}
	//返回set中的数量
	val, err := rdb.SCard(ctx, "set").Result()
	if err == redis.Nil {
		fmt.Println("没有值!")
	} else if err != nil {
		log.Fatalln("scard失败:", err)
	} else {
		fmt.Println("set中此时的长度:", val)

	}
}
func find() {
	//查找set中所有的值,sismember是根据value看是否是成员
	val, err := rdb.SMembers(ctx, "set").Result()
	if err == redis.Nil {
		fmt.Println("没有值!")
	} else if err != nil {
		log.Fatalln("smembers失败:", err)
	} else {
		fmt.Println("set中的成员:", val)

	}
	//随机取几个值(不删除)
	val2, err := rdb.SRandMemberN(ctx, "set", 3).Result()
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("随机获得的3个结果:", val2)
	}

}

//删除,因为删除不存在的值会返回0,所以增加一个判断
func del() {
	val, err := rdb.SRem(ctx, "set", 5).Result()
	if err != nil {
		log.Fatalln("删除失败:", err)
	} else if val == 0 {
		fmt.Println("没有该值!")
	} else {
		fmt.Println("删除成功")
	}
}

func main() {
	add()
	find()
	del()

}
