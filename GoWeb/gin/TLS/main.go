package main

import (
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

//利用errgroup阻塞的特性,来实现单进程多端口服务!
func main() {
	var eg errgroup.Group //收集协程的错误
	//一进程多端口,实际开发中80端口供内部使用,443端口使用https协议 供外部使用!
	inServer := &http.Server{
		Addr:         ":8080",
		Handler:      router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	outServer := &http.Server{
		Addr:         ":8443",
		Handler:      router(), //两个服务的handler相同
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	eg.Go(func() error {
		err := inServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	eg.Go(func() error {
		err := outServer.ListenAndServeTLS("server.pem", "server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	//eg阻塞
	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

}
