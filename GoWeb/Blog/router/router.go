package router

//开启服务,路由分发请求
import (
	"github.com/gin-gonic/gin"
	"github.com/leilei3167/ready-to-go/GoWeb/Blog/controller"
)

func Start() {
	e := gin.Default()
	e.LoadHTMLGlob("templates/*")   //模板路径
	e.Static("/assets", "./assets") //静态资源路径

	e.GET("/", controller.Index)              //主页
	e.GET("/register", controller.GoRegister) //去到注册页面

	e.POST("/register", controller.Register)

	e.Run()
}
