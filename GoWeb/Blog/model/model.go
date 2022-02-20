package model

import "gorm.io/gorm"

type User struct {
	gorm.Model        //集成Gorm
	Username   string `json:"username,omitempty" `
	Password   string `json:"password,omitempty" `
}
//博客模型
type Post struct {
	gorm.Model
	Title   string
	Content string `gorm:"type:text"`
	Tag     string
}