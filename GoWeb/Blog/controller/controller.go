package controller

/* 此包包含的是各种处理函数,即handler */
import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/leilei3167/ready-to-go/GoWeb/Blog/dao"
	"github.com/leilei3167/ready-to-go/GoWeb/Blog/model"
)

//注册用户
//前往注册页面
func GoRegister(c *gin.Context) {
	c.HTML(200, "register.html", nil)
}

//执行注册逻辑
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

//登录
//前往登录页面
func GoLogin(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

//执行登录
func Login(c *gin.Context) {
	username := c.PostForm("username") //从表单中获取用户名和密码
	password := c.PostForm("password")

	fmt.Println(username)
	u := dao.Mgr.Login(username) //执行数据库查询
	if u.Username == "" {        //查不到
		c.HTML(200, "login.html", "用户名不存在!") //第三个参数填入模板{{.}}处
		fmt.Println("用户名不存在")
	} else { //用户名存在则判断密码
		if u.Password != password { //数据库密码和输入的密码不一致
			fmt.Println("密码错误")
			c.HTML(200, "login.html", "密码错误")
		} else {
			fmt.Println("登录成功")
			c.Redirect(301, "/") //跳转回主页
		}

	}

}

//前往主页
func Index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

//博客操作
//博客列表
func GetPostIndex(c *gin.Context) {
	posts := dao.Mgr.GetAllPost()
	c.HTML(200, "postindex.html", posts)

}

//添加博客
func AddPost(c *gin.Context) {
	title := c.PostForm("title")
	tag := c.PostForm("tag")
	content := c.PostForm("content")

	post := model.Post{
		Title:   title,
		Tag:     tag,
		Content: content,
	} //此处应该可以使用shouldbind
	dao.Mgr.AddPost(&post)
	c.Redirect(302, "/post_index")

}

func GoAddPost(c *gin.Context) {

	c.HTML(200, "post.html", nil)

}
