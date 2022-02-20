package model

import "gorm.io/gorm"

type User struct {
	gorm.Model        //集成Gorm
	Username   string `json:"username,omitempty" `
	Password   string `json:"password,omitempty" `
}
