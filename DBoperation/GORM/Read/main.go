package main

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	var err error

	Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN: "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal("error to connect DB:::", err)
	}

}

type User struct {
	gorm.Model
	Name string
	Age  int
	Job  string `gorm:"default:None"`
	Pages
}
type Pages struct {
	Title string

	User_id int
}

func Read() {
	//一,适用于了解结构的情况下
	fmt.Println("---------普通查询----------------")
	//1.1检索单个对象SELECT * FROM users ORDER BY id LIMIT 1;
	var user User
	Db.First(&user)
	fmt.Println("First", user)
	//1.2不指定排序,不指定Orderby
	var user2 User
	Db.Take(&user2)
	fmt.Println("Take:", user2)
	//1.3获取最后一条数据,主键降序
	var user3 User
	Db.Last(&user3)
	fmt.Println("Last:", user3)
	//1.4错误检查,是否为找不到记录的错误
	res := Db.First(&user)
	errors.Is(res.Error, gorm.ErrRecordNotFound)
	//1.5根据不知道将返回的结构,用map来装
	result := make(map[string]interface{})
	Db.Model(&User{}).First(&result)
	fmt.Println("用map:", result)
	fmt.Println(result["name"])
	//1.6如果没有主键,将会按照第一个字段排序

	//二.可以使用内联条件检索对象。传入字符串参数时注意避免SQL注入问题
	//2.1// SELECT * FROM users WHERE id = 10;
	var user4 User
	Db.First(&user4, 2) //也可以传"2"
	fmt.Println("根据主键查找", user4)
	//2.2检索全部对象;相当于select * from,需切片容纳
	var user5 []User
	res1 := Db.Find(&user5)
	fmt.Println("查询所有条目", res1.RowsAffected)
	for _, v := range user5 {
		fmt.Println(v)
	}

	//三.条件
	fmt.Println("-------条件查询-----------------")
	//3.1获取第一条匹配的记录
	var user6 User
	Db.Where("name=?", "wangli").First(&user6)
	fmt.Println(user6)
	//3.2获取全部匹配的记录
	var user7 []User
	Db.Where("user_id=?", "0").Find(&user7)
	for _, v := range user7 {
		fmt.Println(v)
	}
	//3.3 IN条件
	var user8 []User
	Db.Where("name in ?", []string{"leilei", "wangli"}).Find(&user8)
	for _, v := range user8 {
		fmt.Println(v)
	}
	//还支持类似Like , And ,Between等限制语句
	//3.4使用结构体或map做条件
	var user9 []User
	Db.Where(map[string]interface{}{
		"job": "GO",
	}).Find(&user9)
	for _, v := range user9 {
		fmt.Println("用map做条件", v)
	}

	var user10 []User
	Db.Where(&User{Job: "GO", Name: "leilei"}).Find(&user10)
	for _, v := range user10 {
		fmt.Println("用结构体做条件", v)
	}//select from users where name ='leilei' and job='GO';

//limit
//TODO 分页器实现




}

func main() {
	Read()
}
