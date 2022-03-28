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
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		//从读取的配置中获取
		result = (page - 1) * setting.PageSize

	}
	return result
}
