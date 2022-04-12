package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

//基于errgroup实现httpserver的启动和关闭,以及linux 信号的注册和处理

//保证一个退出 全部退出

func main() {
	//创建g
	g, ctx := errgroup.WithContext(context.TODO())

	//路由
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong!"))

	})
	//模拟单个服务器出错
	serverOut := make(chan struct{})
	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		serverOut <- struct{}{}
	})

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	//group中开启监听
	g.Go(func() error {
		return server.ListenAndServe()
	})

	//关闭服务的逻辑,收到信号关闭和其他原因关闭
	g.Go(func() error {
		select {
		case <-ctx.Done(): //因为错误,导致关闭
			log.Println("errgroup exit!")
		case <-serverOut: //因为收到信号 关闭
			log.Println("server will out!")

		}
		timeOutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		log.Println("shutting down server...")
		return server.Shutdown(timeOutCtx)

	})
	g.Go(func() error {
		quit := make(chan os.Signal, 0)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			return fmt.Errorf("get os signal:%v\n", sig)
		}

	})
	//会阻塞在此,监听收到的错误
	fmt.Printf("errgroup exiting %+v\n", g.Wait())
}
