package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

/*
实现优雅关机:
服务端发起关机命令后不是立刻关机,而是等待当前还在处理的请求全部处理完毕后再退出;

*/

func main() {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		time.Sleep(time.Second * 15)
		ctx.JSON(http.StatusOK, gin.H{
			"code": "200",
			"msg":  "after down",
		})
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	//开启协程进行监听
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	//主协程创建接收系统信号的通道
	quit := make(chan os.Signal, 1)
	//KILL会发送SIGTERM信号
	//kill -2 发送SIGINT信号
	//kill -9 发送SIGKILL信号(但不能被捕获)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt) //Notify 会将后面参数指定的信号放入指定chan
	log.Printf("收到信号:%v", <-quit)                                      //阻塞在此 直到收到以上指定的两种信号

	log.Println("Shutting down server...")
	//设置一个超时关闭,等待请求处理完,8秒还没处理完就结束
	ctx, cancle := context.WithTimeout(context.TODO(), time.Second*8)
	defer cancle()
	if err := server.Shutdown(ctx); err != nil && err != ctx.Err() {
		log.Fatalf("关机失败!%v", err)

	}
	log.Println("关闭成功")
}
