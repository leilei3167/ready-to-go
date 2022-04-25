package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

//使用原生rpc时需要使用除gob之外的编码器需要实现接口,对于json标准库也单独实现了一个包便于使用
type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func main() {
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	arith := new(Arith)
	rpc.Register(arith)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		// 注意这一行,其实就是单独实现了一个解码器
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
