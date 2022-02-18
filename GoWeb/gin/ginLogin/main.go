package main

import "github.com/gin-gonic/gin"

func Login(c *gin.Context) {
	c.HTML(200, "login.html", nil)

}
func DoLogin(c *gin.Context) {
	//取出post请求Form中的username和password
	username := c.PostForm("username")
	password := c.PostForm("password")
	//打开welcome文件响应,并将值填入到{{.}}中
	c.HTML(200, "welcome.html", gin.H{
		"username": username,
		"password": password,
	})

}
func main() {
	//1.创建Engine
	e := gin.Default()
	//2.加载指定目录的HTML文件
	e.LoadHTMLGlob("templates/*")
	//3.指定GET请求的路径和对应的方法
	e.GET("/login", Login)
	e.POST("/login", DoLogin)
	//4.开启服务
	e.Run()
}
