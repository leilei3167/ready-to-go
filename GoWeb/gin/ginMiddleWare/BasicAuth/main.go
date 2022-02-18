package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//模拟一些需要权限才能访问的数据
var secrets = gin.H{
	"leilei":   gin.H{"gender": "men", "num": "123"},
	"yangzhen": gin.H{"gender": "men", "num": "567"},
	"haoyun":   gin.H{"gender": "men", "num": "890"},
}

func main() {
	r := gin.Default()

	//路由组使用gin.BasicAuth中间件,设定/admin路径只有对Accounts内的用户开放
	//gin.Accounts是map[string]string类型的快捷方式
	authorize := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"leilei":   "123",
		"yangzhen": "123",
		"haoyun":   "123",
	}))

	//
	authorize.GET("/secrets", func(c *gin.Context) {
		//获取用户,他是由中间件设置的
		user := c.MustGet(gin.AuthUserKey).(string) //类型断言,将返回的interface断言string
		if secret, ok := secrets[user]; ok {        //如果secrets中有user对应的名字
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})

		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}

	})
	r.Run(":8080")
}
