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

	e.GET("/", controller.Index) //主页
	//注册
	e.GET("/register", controller.GoRegister)
	e.POST("/register", controller.Register)
	//登录
	e.GET("/login", controller.GoLogin)
	e.POST("login", controller.Login)

	//博客操作
	e.GET("/post_index", controller.GetPostIndex) //列表
	e.POST("/post", controller.AddPost)           //添加博客
	e.GET("/post", controller.GoAddPost)          //跳转到添加博客页面

	e.Run()

}
