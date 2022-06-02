package main

import (
	"golang.org/x/net/context"
	"mapping/internal/master"
	"mapping/internal/pkg/db"
)

/*var (
	brokers = flag.String("-b", "", "brokers地址")
)*/
//var addr = flag.String("s", "0.0.0.0:8080", "服务监听地址")

func main() {
	/*	r := gin.Default()
		r.Use(mid.RateLimitMiddleware(time.Second, 5, 1)) //限流,处理请求的最大数量暂时设置为5
		r.POST("/scantask", v1.ScanFromUpload)
		r.GET("/ipinfo")
		if err := r.Run(); err != nil {
			panic(err)
		}*/

	master.NewMaster(":8080", "mongodb://localhost:27017").Run()
	err := db.MgoClient.Disconnect(context.TODO())
	if err != nil {
		return
	}
}
