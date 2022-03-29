package util

import (
	"github.com/gin-gonic/gin"
	"github.com/leilei3167/ready-to-go/GoWeb/Blog2/pkg/setting"
	"github.com/unknwon/com"
)

/* 分页页码获取 */
//go get -u github.com/unknwon/com

func GetPage(c *gin.Context) int {
	result := 0
	//gin从url的查询字段获取page的值,并转换为int
	//先强转为StrTo类型,在调用其Int()方法
	page, _ := com.StrTo(c.Query("page")).Int()

	if page > 0 {
		//从读取的配置中获取
		result = (page - 1) * setting.PageSize

	}
	//返回页码
	return result
}
