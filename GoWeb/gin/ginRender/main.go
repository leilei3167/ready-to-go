package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	//模板解析
	//router.LoadHTMLFiles("templates/index.html")
	//使用不同目录下名称相同的模板
	router.LoadHTMLGlob("templates/**/*") //**表示templates下面的一级目录

	router.GET("users/index", func(c *gin.Context) {
		//通过文件名渲染模板
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "users目录下的模板我是雷磊",
		})

	})
	router.GET("posts/index", func(c *gin.Context) {
		//通过文件名渲染模板
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "posts目录下的模板我是雷磊",
		})

	})
	router.Run(":9090")
}
