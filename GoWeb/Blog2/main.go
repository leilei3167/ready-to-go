package main

import (
	"fmt"
	"net/http"

	"github.com/leilei3167/ready-to-go/GoWeb/Blog2/pkg/setting"
	"github.com/leilei3167/ready-to-go/GoWeb/Blog2/routers"
)

func main() {
	r := routers.InitRouter()
	//配置下server
	s := &http.Server{
		//配置文件指定的端口
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        r,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20, //左移20位
	}

	//开启
	s.ListenAndServe()

}
