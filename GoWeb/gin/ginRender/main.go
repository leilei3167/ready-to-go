package main

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Static("/xxx", "./statics")
	//自定义一个函数并加载到模板里,funcMap是装函数的容器
	router.SetFuncMap(template.FuncMap{
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
	})
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
			"title": "<a href='https://www.baidu.com'>雷磊你好</a>",
		})

	})
	router.Run(":9090")
}
