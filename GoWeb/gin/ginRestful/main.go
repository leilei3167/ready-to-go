package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

//用gin实现restful风格
//user的信息类型
type User struct {
	UId  int    `json:"uid"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

//用切片模拟一个数据库
var users = make([]User, 3)

//加入一些初始数据
func init() {
	u1 := User{1, "tom", 20}
	u2 := User{2, "kite", 30}
	u3 := User{3, "rose", 40}
	users = append(users, u1)
	users = append(users, u2)
	users = append(users, u3)
	fmt.Println(users)
}

//先写一个根据id查询的方法,返回实例和下标
func find(uid int) (*User, int) {
	for k, v := range users { //变量数据库,现实可能是查询语句
		if v.UId == uid {
			return &v, k //找到则返回
		}
	}
	return nil, -1
}

//Get请求的处理方法,根据id查询
func Finduser(c *gin.Context) {
	uid := c.Param("uid")      //获取路径变量的值
	id, _ := strconv.Atoi(uid) //将string转化为int
	u, _ := find(id)           //将获得的id调用查询
	c.JSON(200, u)

}

//POST请求,添加用户(根据传入的json)
func AddUser(c *gin.Context) {
	var newUser User
	c.ShouldBindJSON(&newUser) //传入的json与newUser绑定
	//u4 := User{4, "Joe", 50}
	users = append(users, newUser)
	c.JSON(200, users)
}

//根据id删除
func DelUser(c *gin.Context) {
	uid := c.Param("uid")
	id, _ := strconv.Atoi(uid)
	_, i := find(id)
	users = append(users[:i], users[i+1:]...) //取该下标前的切片,将它之后的切片加入,实现删除
	c.JSON(200, users)
}

//PUT,修改某个用户
func UpdateUser(c *gin.Context) {
	uid := c.Param("uid")
	id, _ := strconv.Atoi(uid)
	u, _ := find(id)
	u.Name = "修改的Name"
	c.JSON(200, u)
}

func main() {
	e := gin.Default()
	e.GET("/user/:uid", Finduser)   //根据id查询
	e.POST("/user/", AddUser)       //根据传入的json添加用户
	e.DELETE("/user/:uid", DelUser) //根据传入的id删除
	e.PUT("/user/:uid", UpdateUser) //根据id修改名字

	e.Run()
}
