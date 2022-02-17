package main

import (
	"context"
	"fmt"
	"net/http"
)

type fucker struct {
}

/*服务器端核心就是配置Server和各种Handler*/
func (fucker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background() //创建根上下文,便于子处理器进行包装
	fmt.Fprintln(w, "hello world!!!!")
	switch r.Method {
	case http.MethodGet:
		//执行处理方法
		handlerGet(ctx, w, r)
	case http.MethodPost:
		//执行处理方法
		//...
	}
}
func handlerGet(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello here is get")
}
func handlerput(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello here is get")
}

func main() {
	//1.创建路由
	router := http.NewServeMux()

	//2.创建实例并实现ServeHTTP方法,并注册到路由
	router.Handle("/hello", &fucker{})
	router.HandleFunc("/hi", handlerput) //此方法也能够实现注册
	//3.配置server
	server := http.Server{Addr: ":8080",
		Handler: router}
	//4.用自定义的配置运行监听
	server.ListenAndServe()
}
