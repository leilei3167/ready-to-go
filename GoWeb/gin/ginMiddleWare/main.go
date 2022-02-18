package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func TestMW(c *gin.Context) {
	c.String(200, "hello,%s", "leilei")
}

func MyMiddleware1(c *gin.Context) {
	fmt.Println("我的第一个中间件")
}

func MyMiddleware2(c *gin.Context) {
	fmt.Println("我的第二个中间件")
}

func main() {
	//1.e := gin.Default() 默认会有Logger和Recovery中间件
	//2. e:=gin.New() 此方法创建引擎不会有中间件\

	//3.自定义中间件

	e := gin.Default()

	e.Use(MyMiddleware1, MyMiddleware2) //要求的参数实际就是处理器函数

	e.GET("testMW", TestMW)
	e.Run()
}
