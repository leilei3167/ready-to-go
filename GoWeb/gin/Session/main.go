package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//设置储存的介质,此处使用自带存储引擎
	store := cookie.NewStore([]byte("seceree12112"))
	//设置全局中间件
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/hello", func(c *gin.Context) {
		//初始化session
		session := sessions.Default(c)

		if session.Get("hello") != "world" {
			//设置
			session.Set("hello", "world")
			//删除
			session.Delete("leilei")
			//保存session数据
			session.Save()
			//session.Clear 删除所有session

		}
		c.JSON(200, gin.H{
			"hello": session.Get("hello"),
		})

	})
	r.Run(":8088")

}
