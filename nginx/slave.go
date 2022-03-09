package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/slave", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"我是cong机 ":   "测试成功",
			"我cong机1 ":   "测试成功",
			"我congji 11": "测试成功",
		})

	})

	r.Run()
}
