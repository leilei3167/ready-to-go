package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type stu struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	r := gin.Default()
	r.GET("/json", func(c *gin.Context) {
		//方法1:使用map
		/*data := map[string]interface{}{
			"name": "哈哈哈",
			"age":  18,
		}*/
		data := gin.H{"name": "leilei", "age": 18}
		//方法二,用结构体
		data2 := stu{
			Name: "啊哈哈",
			Age:  19,
		}
		c.JSON(http.StatusOK, data2)
		c.JSON(http.StatusOK, data)

	})

	err := r.Run(":9090")
	if err != nil {
		return
	}
}
