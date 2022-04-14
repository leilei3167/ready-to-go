package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//1.如果要创建自定义的JWT的负载,则需要创建新结构体并嵌入StandardClaims
type MyClaims struct {
	Username string `json:"username`
	jwt.StandardClaims
}

//2.设置JWT的过期时间
const TokenExpireDuration = time.Hour * 12

//3.设置Secret(可以是服务器的私钥)
var MySecret = []byte("你好!!")

//4.生成JWT
func GenToken(username string) (string, error) {
	c := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "leilei",
		},
	}
	//使用指定签名方法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	//使用指定的secret签名并获取完整的字符串TOKEN
	return token.SignedString(MySecret)

}

//5.解析JWT
// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//在gin中使用
func main() {
	r := gin.Default()
	r.POST("/auth", authHandler) //注册一个验证密码并发放TOken的处理器
	r.GET("/home", JWTAuthMiddleWare(), func(c *gin.Context) {
		username := c.MustGet("username").(string)
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "success",
			"data": gin.H{"username": username},
		})
	})
	r.Run()
}

func authHandler(c *gin.Context) {
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "200",
			"msg":  "无效的信息",
		})
	}
	//校验密码
	if user.Username == "leilei" && user.Password == "leilei123" {
		tokenString, _ := GenToken(user.Username)

		c.JSON(http.StatusOK, gin.H{
			"code": "2000",
			"msg":  "success",
			"data": gin.H{
				"token": tokenString,
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "222",
		"msg":  "账户或者密码错误",
	})
	return
}

//验证token应该通过使用中间件的形式
func JWTAuthMiddleWare() func(c *gin.Context) {
	return func(c *gin.Context) {
		//客户端携带TOKEN有三种方式:1.header处;2.body中;3.URI中
		//如何获取取决于实际业务,此处以放在Header中为例
		authHeader := c.Request.Header.Get("Authorization") //从中获取指定字段
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort() //停止执行后续的处理器
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		/* 	if !(len(parts) == 2 && parts[0] == "Bearer") {
			log.Println(parts, parts[0])
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		} */
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(parts[0])
		if err != nil { //err==TokenExpired
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}

}
