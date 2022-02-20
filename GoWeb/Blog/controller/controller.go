package controller

/* 此包包含的是各种处理函数,即handler */
import (
	"github.com/gin-gonic/gin"
	"github.com/leilei3167/ready-to-go/GoWeb/Blog/dao"
	"github.com/leilei3167/ready-to-go/GoWeb/Blog/model"
)

//添加用户,用于POst请求
func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user := model.User{
		Username: username,
		Password: password,
	}
	//调用dao写入数据库
	dao.Mgr.Register(&user)
	c.Redirect(301, "/") //写入数据库后执行跳转,跳转到首页
}

//列出所有的用户
func GoRegister(c *gin.Context) {
	c.HTML(200, "register.html", nil)
}

func Index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}
