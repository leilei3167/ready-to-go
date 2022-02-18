package main

import (
	"github.com/gin-gonic/gin"
)

//获取Get中的查询字段
func testGet(c *gin.Context) {
	//4.It is shortcut for `c.Request.URL.Query().Get(key)`
	//获取查询字段中的"username" 如果没有返回空字符串
	s := c.Query("username")
	//获取查询字段中的password字段,没有的话设置默认值为123
	pwd := c.DefaultQuery("password", "123")

	c.String(200, "username:%s,password:%s", s, pwd)

}

//获取Post中的表单信息
func testPost(c *gin.Context) {
	nam := c.PostForm("name")
	pas := c.PostForm("password")
	c.String(200, "你的用户名是%s,密码是%s", nam, pas)

}

//获取路径参数
func testPath(c *gin.Context) {
	s1 := c.Param("name")
	s2 := c.Param("age")
	c.String(200, "你的用户名是%s,密码是%s", s1, s2)
}

//post请求中同时获取查询字段
func testGetAndPost(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	key := c.PostForm("key")
	c.String(200, "Page:%s, Key:%s", page, key)

}
func main() {
	//1.设置默认引擎
	e := gin.Default()
	//2.注册方法
	e.GET("/testGet", testGet)
	e.POST("/testGet", testPost)
	e.GET("/testPath/:name/:age", testPath)
	//post请求中同时获取查询字段
	e.POST("/testGetPost", testGetAndPost)
	//3.在默认端口开启监听
	e.Run()
}
