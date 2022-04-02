package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

type User struct {
	Uid      string
	Password string
}

//https://segmentfault.com/a/1190000039953487
func main() {

	r := gin.Default()
	//下载gin的redis驱动
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("firstKey"))
	//使用中间件
	r.Use(sessions.Sessions("test_session", store))

	//访问此路径时检查其有没有登录状态
	r.POST("/in", func(ctx *gin.Context) {
		//创建操作句柄
		
		session := sessions.Default(ctx)

		//获取某个sessionid对应的值
		var a User
		ctx.ShouldBind(&a)
		cookie, err := ctx.Cookie("test_session")
		if err != nil && err == http.ErrNoCookie {
			//没有session的话就跳转至登录逻辑处理
			//假设登录后,设置session
			session.Set(a.Uid, a.Password)
			session.Save()
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "你还未登录!",
			})

		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"ok cookie:test_session :": cookie,
			})

		}

	})
	r.Run(":8088")
}
