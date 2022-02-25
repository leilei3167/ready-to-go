package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func timer(c *gin.Context) {
	fmt.Println("中间件开始计时!")
	a := time.Now()
	c.Next() //开始计时后执行后续的处理器,处理器运行完后这里也执行完毕
	fmt.Println("计时完毕,耗时:", time.Since(a).Milliseconds())

}

func main() {
	r := gin.Default()

	r.GET("/", timer, func(c *gin.Context) {
		c.String(200, "hello!!!")

	})

	r.Run()

}
