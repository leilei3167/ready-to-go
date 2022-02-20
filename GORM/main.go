package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//继承GORM
type Product struct {
	//Model为gorm包自带的结构体,放入此就实现了继承
	gorm.Model
	Code  string
	Price uint
}

//创建表
func create(db *gorm.DB) {
	//创建表 迁移 schema,执行将会根据结构体字段直接在数据库创建数据表
	db.AutoMigrate(&Product{})
}

//插入数据
func insert(db *gorm.DB) {
	p := Product{
		Code:  "1001",
		Price: 100,
	}
	db.Create(&p)

}

//查询

func find(db *gorm.DB) {
	var p Product
	//根据主键查询,查询主键为1的数据
	db.First(&p, 1)
	fmt.Println("根据主键查询:", p)

	//根据条件来查询
	db.First(&p, "code=?", "1001")
	fmt.Println("根据条件查询结果:", p)

}

//更新
func update(db *gorm.DB) {
	//更新要先查到数据
	var p Product
	db.First(&p, 1) //将查到的结果装到p
	//更新单个字段
	db.Model(&p).Update("Price", 1000) //执行更新

	//更新多个字段
	db.Model(&p).Updates(Product{Price: 1001, Code: "1002"})

	//db.Model(&p).Updates(map[string]interface{}{"Price":1003,"Code":"1003"})

}

//删除,此删除为软删除,只是在数据库加上一个删除的标记
func del(db *gorm.DB) {
	//同样是要先找到
	var p Product
	db.First(&p, 1)
	db.Delete(&p, 1)

}

func main() {
	//连接到数据库(mysql)
	dsn := "root:8888@/golang_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//find(db)
	//update(db)
	del(db)
}
