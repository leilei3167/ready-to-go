package main

import "github.com/gin-gonic/gin"

func HELLO(c *gin.Context) {
	//c.String(200, "hello,leilei%s", "saozhu")
	c.JSON(200, gin.H{
		"name": "tom",
		"age":  12})
}
func main() {

	e := gin.Default()
	e.GET("/hello", HELLO)

	e.Run() //8080
}
