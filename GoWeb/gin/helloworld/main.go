package main

import "github.com/gin-gonic/gin"

func sayHello(c *gin.Context) {
	// H is a shortcut for map[string]interface{}
	c.JSON(200, gin.H{"message": "hahaha",
		"location": "成都"})

}
func main() {
	//创建默认路由
	r := gin.Default()
	//注册路由
	r.GET("/hello", sayHello) //指定handler的方法

	//-------------------------------------------------------------
	//restful架构
	r.GET("/book", func(c *gin.Context) { //传匿名函数的方法
		c.JSON(200, gin.H{"method": "GET"})

	})
	r.POST("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{"method": "POST"})

	})
	r.PUT("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{"method": "PUT"})

	})

	r.DELETE("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{"method": "DELETE"})

	})

	//启动服务
	r.Run("0.0.0.0:9090") //不填的话默认是8080端口
}
