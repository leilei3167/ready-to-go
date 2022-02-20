package dao

import (
	"github.com/leilei3167/ready-to-go/GoWeb/Blog/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//调用接口
type Manager interface {
	//用户操作
	Register(user *model.User)
	Login(username string) model.User
	// 博客操作
	AddPost(post *model.Post)
	GetAllPost() []model.Post
	GetPost(pid int) model.Post
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
	db.AutoMigrate(&model.Post{}) //根据字段创建表
}

//注册
func (mgr *manager) Register(user *model.User) {
	mgr.db.Create(user)
}

func (mgr *manager) Login(username string) model.User {
	var user model.User
	//查找,根据传入的用户名来查找
	mgr.db.Where("username=?", username).First(&user)
	return user
}

//博客操作
//添加博客
func (mgr *manager) AddPost(post *model.Post) {
	mgr.db.Create(post) //插入数据
}

//找到所有的博客
func (mgr *manager) GetAllPost() []model.Post {
	posts := make([]model.Post, 10)
	mgr.db.Find(&posts) //找到所有的post
	return posts
}

//查找博客
func (mgr *manager) GetPost(pid int) model.Post {
	var post model.Post
	mgr.db.First(&post, pid)
	return post
}
