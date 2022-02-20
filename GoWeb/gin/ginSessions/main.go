package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	//注入中间件
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/hello", func(c *gin.Context) {
		session := sessions.Default(c) //实例化sessions

		if session.Get("hello") != "world" { //获得sessions并判断
			//不等于则设置
			session.Set("hello", "world")
			//保存
			session.Save()
		}
		//输出一个key为hello,value为sessions
		c.JSON(200, gin.H{"hello": session.Get("hello")})
	})
	r.Run(":8000")
}
