package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

/*
gin常用中间件:
	gin-jwt:实现JWT认证
	cors:实现跨域相关
	sessions:session管理
	authz:基于casbin的授权中间件
	gin-limit:并发控制
	requestid:给每个request生成uuid,并添加返回的X-Request-ID

*/
//初始化,加载设置中间件(跨域,请求ID,认证)
func router() http.Handler {
	router := gin.Default()
	productHandler := newProductHandler()

	//跨域设置
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	})) //限制跨域请求

	// 给每个请求附加RequestID
	router.Use(requestid.New())

	//全局认证(postman中选择basicauth),也可以在不同的分组中使用,实现特定的鉴权
	router.Use(gin.BasicAuth(gin.Accounts{"foo": "bar", "leilei": "123"}))

	//分组注册路由
	v1 := router.Group("v1") //分组
	{
		productv1 := v1.Group("/products")
		{
			productv1.POST("", productHandler.Create)  //v1/products
			productv1.GET(":name", productHandler.Get) //v1/products/:name
		}
		testv1 := v1.Group("/test")
		{
			testv1.GET("", gin.HandlerFunc(func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"msg": "hello"})
			}))
		}

	}
	return router
}
