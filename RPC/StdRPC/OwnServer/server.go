package main

import (
	"log"
	"net/http"
	"net/rpc"
)

//直接调用Register,ServeConn ServeCodec都是默认调用的DefaultServer的相关的方法
//DefaultServer是全局共享的,如果有第三方库也使用了相关方法,则有可能造成错误;
//所以推荐使用自定义的NewServer
type Arith struct{}

func main() {
	arith := new(Arith)
	server := rpc.NewServer()
	server.RegisterName("math", arith)
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("serve error:", err)
	}

	//一般标准库中都会提供一个默认的实现便于直接使用,如log,http这些都会有默认的实例供调用

	//但是实际开发中建议自己创建实例,避免不必要的干扰
}
