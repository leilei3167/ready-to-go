package main

//https://www.liwenzhou.com/posts/Go/validator_usages/
import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

//gin内部的validator使用的是第三方包
type SignUpParam struct {
	Age  uint8  `json:"age" binding:"gte=1,lte=130"`
	Name string `json:"name" binding:"required"`
	//要求输入标准邮件格式
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	//eqfield要求和某个字段相等
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

//可以定义翻译器,将报错信息翻译成中文
var trans ut.Translator

// InitTrans 初始化翻译器
func InitTrans(locale string) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

func main() {
	//创建翻译
	if err := InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed,err %v", err)
		return
	}

	r := gin.Default()

	r.POST("/signup", func(c *gin.Context) {
		var u SignUpParam
		if err := c.ShouldBind(&u); err != nil {

			//判断是否是验证的错误
			errs, ok := err.(validator.ValidationErrors)
			if !ok { //不是直接返回
				c.JSON(http.StatusOK, gin.H{
					"msg": err.Error(),
				})
				return

			}
			//否则对其进行翻译
			c.JSON(http.StatusOK, gin.H{
				"msg": errs.Translate(trans),
			})
			return
		}
		// 保存入库等业务逻辑代码...

		c.JSON(http.StatusOK, "success")
	})

	_ = r.Run(":8999")
}
