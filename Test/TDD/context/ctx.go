package ctx

import (
	"context"
	"fmt"
	"net/http"
)

/*
//v1
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context() //取出请求中的上下文

		data := make(chan string, 1)
		go func() {
			data <- store.Fetch()
		}()

		select {
		case d := <-data:
			fmt.Fprint(w, d)
		case <-ctx.Done(): //当请求中断时,会执行这里
			store.Cancle()
		}

	}
}

type Store interface {
	Fetch() string
	//需要增加用户取消,程序停止的接口

	Cancle()
} */

type Store interface {
	Fetch(ctx context.Context) (string, error)
}

//服务器不应该负责主动取消,而是依赖于请求的上下文,传递给下游的工作函数
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := store.Fetch(r.Context())
		if err != nil {
			return
		}
		fmt.Fprint(w, data)
	}
}
