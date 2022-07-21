package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()

	e.GET("/panic", FakeMid(time.Second*2), TureHandler)
	e.Run()
}

func FakeMid(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("fake mid begin!")

		c.Next()

		fmt.Println("fake mid end!")

	}
}

func TureHandler(c *gin.Context) {
	time.Sleep(time.Second * 2)
	go func() {
		panic("panic")

	}()
	time.Sleep(time.Second * 10)

}
