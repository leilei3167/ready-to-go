package dao

import (
	"github.com/leilei3167/ready-to-go/GoWeb/Blog/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//调用接口
type Manager interface {
	Register(user *model.User)
}

type manager struct {
	db *gorm.DB
}

//外部调用时都通过接口来调用
var Mgr Manager

//初始化函数只会执行一次,用来链接数据库
func init() {
	dsn := "root:8888@/golang_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Mgr = &manager{db: db}        //创建一个manager实例
	db.AutoMigrate(&model.User{}) //根据字段创建表
}

func (mgr *manager) Register(user *model.User) {
	mgr.db.Create(user)
}
