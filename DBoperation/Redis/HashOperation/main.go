package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

/*服务器一般会将用户登录信息保存到redis中,hash表特别适合用于存储对象,而如果用string来存储的话,每次存储都
得序列化成字符串才能存入,从redis中获取之后又需要再反序列化,开销特别大
而使用redis我们可以通过key(如用户id)+field(用户的字段) 就可以操作对应数据了
*/
type User struct {
	name   string `redis:"name"`
	id     string `redis:"id""`
	age    int    `redis:"age" `
	addr   string `redis:"addr" `
	phone  string `redis:"phone" `
	gender string `redis:"gender"`
}

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

//添加数据
func add() {
	//支持以下几种添加方式
	user4 := User{
		name:   "雷锋",
		id:     "123321",
		age:    19,
		addr:   "成都市双流区",
		phone:  "18602843167",
		gender: "男",
	}
	/*	rdb.HSet(ctx, "user1", "name", "ahah", "age", "nihnin")
		rdb.HSet(ctx, "user2", []string{"name", "value", "age", "13"})
		rdb.HSet(ctx, "user3", map[string]interface{}{"id": "111111", "age": "123", "name": "leieli"})
	*/

	err := rdb.HSet(ctx, "user4", user4).Err()
	if err != nil {
		fmt.Println("插入失败", err)
	} else if err == redis.Nil {
		fmt.Println("返回空")
	} else {
		fmt.Println("插入成功")
	}
}

//查
func find() {
	//查找所有字段
	reslut1, _ := rdb.HGetAll(ctx, "user1").Result()
	reslut2, _ := rdb.HGetAll(ctx, "user2").Result()
	reslut3, _ := rdb.HGetAll(ctx, "user3").Result()
	fmt.Println("user1:", reslut1)
	for s, s2 := range reslut1 {
		fmt.Println(s, s2)
	}
	fmt.Println("user2:", reslut2)
	for s, s2 := range reslut2 {
		fmt.Println(s, s2)
	}
	fmt.Println("user3:", reslut3)
	for s, s2 := range reslut3 {
		fmt.Println(s, s2)
	}
	//查找指定key的指定字段HGet
	addr, err := rdb.HGet(ctx, "user4", "addr").Result()
	if err != nil {
		fmt.Println("HGet出错:", err)
	}
	fmt.Println(addr)
}
func main() {
	add()
	find()
}
