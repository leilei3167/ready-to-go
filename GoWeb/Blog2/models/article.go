package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Article struct {
	Model

	//声明为索引 gorm:"index"
	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"` //实现Article和Tag关联查询的功能

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())

	return nil
}

func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())

	return nil
}

//增删改查
//是否存在
func ExistArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id=?", id).First(&article)

	if article.ID > 0 {
		return true

	}
	return false

}

//获取文章的数量
func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)

	return
}

//分页获取文章
//preload是一个预加载器,他会执行2条sql

func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

	return
}

//单个文章

func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)

	//因为嵌套了Tag 所以可以Related进行关联查询
	db.Model(&article).Related(&article.Tag)

	return
}

//修改
func EditArticle(id int, data interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)

	return true
}

//添加
func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	})

	return true
}

func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(Article{})

	return true
}
