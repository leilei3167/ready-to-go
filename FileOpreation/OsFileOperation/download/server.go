package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/download", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "successed!",
			"file":   "this is downloadfile",
		})

	})
	r.Run()
}
