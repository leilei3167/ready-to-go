package main

import "github.com/gin-gonic/gin"

type User struct {
	Username string   `form:"username"`
	Password string   `form:"password"`
	Hobby    []string `form:"hobby"` //因为hobby是多值的 所以要用切片装
	Gender   string   `form:"gender"`
	City     string   `form:"city"`
}
type Query struct {
	name string `form:"name"`
	text string `form:"text"`
}

//
//以html页面响应
func GoRegister(c *gin.Context) {
	c.HTML(200, "testform.html", nil)
}

//高效的办法,实现创建好结构体,直接装form中的值
func formbind(c *gin.Context) {
	var user User
	c.ShouldBind(&user)
	c.String(200, "form data:%s", user)

}
func testQueryBind(c *gin.Context) {
	//还可以绑定查询参数
	var query Query
	c.ShouldBind(&query)
	c.String(200, "query 获取到了查询字段:%s", query)

}

//低效的方法,挨个取出form值
func Register(c *gin.Context) {
	//获取表单信息
	username := c.PostForm("username")
	password := c.PostForm("password")
	hoby := c.PostFormArray("hobby") //因为hoby是多选项的
	gender := c.PostForm("gender")
	city := c.PostForm("city")

	c.String(200, "名字:%s,密码:%s,爱好:%s,性别:%s,城市:%s", username, password,
		hoby, gender, city)

}

func main() {
	e := gin.Default()
	e.LoadHTMLGlob("templates/*")
	e.GET("/testForm", GoRegister) //执行展示
	//获取页面提交的表单信息
	e.POST("/register", formbind) //表单中设置的action跳转到/register
	e.GET("/testQueryBind", testQueryBind)
	e.Run(":9090")
}
