package main

import (
	"github.com/gin-gonic/gin"
	v1 "mapping/internal/master/api/v1"
	"mapping/internal/pkg/mid"
	"time"
)

/*var (
	brokers = flag.String("-b", "", "brokers地址")
)*/
//var addr = flag.String("s", "0.0.0.0:8080", "服务监听地址")

func main() {
	r := gin.Default()
	r.Use(mid.RateLimitMiddleware(time.Second, 5, 1)) //限流,处理请求的最大数量暂时设置为5
	r.POST("/scantask", v1.ScanFromUpload)
	r.GET("/ipinfo")
	if err := r.Run(); err != nil {
		panic(err)
	}

}
