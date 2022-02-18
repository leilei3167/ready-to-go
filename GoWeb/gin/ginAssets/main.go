package main

import "github.com/gin-gonic/gin"

func Go(c *gin.Context) {
	c.HTML(200, "index.html", nil)

}

func main() {
	e := gin.Default()
	e.Static("/assets", "./assets")
	e.LoadHTMLFiles("index.html")
	e.GET("/login", Go)
	e.Run()
}

//GoWeb/gin/ginAssets/assets
