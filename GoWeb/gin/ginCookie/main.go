package main

import (
	"github.com/gin-gonic/gin"
)

func hand(c *gin.Context) {
	s, err := c.Cookie("username")
	if err != nil {
		s = "leilei"
		c.SetCookie("username", s, 60*60, "/", "localhost", false, true)

	}
	c.String(200, "测试cookie")
}

func main() {

	e := gin.Default()
	e.GET("/testcookie", hand)
	e.Run()

}
