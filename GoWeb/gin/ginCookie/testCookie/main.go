package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//中间件
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端cookie并校验
		if cookie, err := c.Cookie("abc"); err == nil {
			if cookie == "123" {
				//Next只能用于中间件,此处意义:满足条件的才执行后续函数
				c.Next()
				return
			}
		}
		// 否则说明还没有被设置cookie,返回错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "err"})
		// 若验证不通过，不再调用后续的函数处理(Abort用于验证中间件)
		c.Abort()

	}
}

func main() {
	// 1.创建路由
	r := gin.Default()
	r.GET("/login", func(c *gin.Context) {
		// 设置cookie
		c.SetCookie("abc", "123", 60, "/",
			"localhost", false, true)
		// 返回信息
		c.String(200, "Login success!")
	})
	//登录home时先执行AuthMiddleWare(返回值必须是满足handler格式的函数)
	r.GET("/home", AuthMiddleWare(), func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "home"})
	})
	r.Run(":8000")
}
