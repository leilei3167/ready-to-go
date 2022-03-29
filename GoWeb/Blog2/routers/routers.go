//用于同意注册处理器和中间件
package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/leilei3167/ready-to-go/GoWeb/Blog2/pkg/setting"
	v1 "github.com/leilei3167/ready-to-go/GoWeb/Blog2/routers/api/v1"
)

/* 此包将注册处理器的逻辑从main从分离出来 */
func InitRouter() *gin.Engine {
	//手动添加中间件
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	//后续可以在此进行扩展中间件
	gin.SetMode(setting.RunMode)

	//注册路由组
	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/tags", v1.GetTags)
		apiv1.POST("/tags", v1.AddTag)
		apiv1.PUT("/tags/:id", v1.EditTag)
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		//所有文章
		apiv1.GET("/articles", v1.GetArticles)
		apiv1.GET("/ariticles/:id", v1.GetArticle)
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)

	}

	return r

}
