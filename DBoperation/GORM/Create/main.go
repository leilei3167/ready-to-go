package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//要注意,声明全局Db后,在连接数据库时要用=,而不是:=,否则全局Db将依然是初始nil值
var (
	Db *gorm.DB
)

func init() {
	var err error
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	//基本模式,使用默认配置,但实际基本都是自己配置
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	Db, err = gorm.Open(mysql.New(mysql.Config{
		//数据库相关配置
		DSN: dsn,
	}), &gorm.Config{
		SkipDefaultTransaction: true, //禁用默认的事务(默认所有写都会加事务)

	})
	if err != nil {
		log.Fatal("数据库链接出错!", err)

	}

	//设置连接池,获取底层db
	sqlDB, err := Db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	//自动迁移表格
	fmt.Println("链接数据库成功")
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

func migrate() {
	//自动迁移,会自动创建外键约束,可以设置禁用
	err := Db.AutoMigrate(&User{}, &Pages{})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	migrate()
	create()
}

//创建记录
func create() {
	//1.普通创建
	user := User{Name: "leilei", Age: 13, Job: "go"}
	res := Db.Create(&user)
	if res.Error != nil {
		log.Fatal("插入数据出错:", res.Error)
	}
	fmt.Println("影响行数", res.RowsAffected)
	fmt.Println(user)

	//2.给定的字段创建,虽然有Age但选择不插入
	user2 := User{Name: "lei", Age: 11, Job: "Rust1"}
	res = Db.Select("Name", "Job").Create(&user2)
	if res.Error != nil {
		log.Fatal("插入数据出错:", res.Error)
	}
	fmt.Println("影响行数", res.RowsAffected)
	fmt.Println(user2)

	//3.批量插入,传入实例的切片
	var users = []User{
		{Name: "yanzghen", Age: 77, Job: "huashui"}, {Name: "haoytun", Age: 32, Job: "boss"},
	}

	Db.Create(&users)
	//检查看是否生成了id
	for _, v := range users {
		fmt.Println(v.ID)
	}
	//4.根据map[string]interface{}来创建,关联的部分不会自己填充
	Db.Model(&User{}).Create(map[string]interface{}{
		"Name": "wangli", "Age": 18, "Job": "GO",
	})

	//创建关联数据时,如果关联值是非零值，这些关联会被 upsert，且它们的 Hook 方法也会被调用
	Db.Create(&User{
		Name:  "zhangsan ",
		Job:   "GOOOOO",
		Age:   15,
		Pages: Pages{Title: "daa,xiaya"},
	})

}

//钩子,将会在指定的实例执行某个动作前执行
//此函数在Create User实例时执行一次检查 Job是否为Rust
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Job == "Rust" {

		return errors.New("no Rust")
	}
	return

}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.ID == 0 {
		return errors.New("失败")
	}
	fmt.Println("AfertCreate钩子:", u)
	return

}
