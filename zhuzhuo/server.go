package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/echo", res)

	r.Run()

}

var m = make([]string, 0)

func res(c *gin.Context) {
	a := c.Query("input")
	if a == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error:": "请不要输入空字符",
		})
		return
	} else {
		m = append(m, a)
	}
	c.JSON(http.StatusOK, m)

}
