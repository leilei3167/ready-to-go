package model

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db *gorm.DB
)

func init() {
	var err error

	Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN: "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
	}), &gorm.Config{
		SkipDefaultTransaction: true, //关闭默认事务
	})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := Db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	fmt.Println("初始化数据库成功")
}

//可用标签设置默认值
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
