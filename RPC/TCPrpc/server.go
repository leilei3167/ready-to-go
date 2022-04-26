package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
)

//待注册的方法 实例
type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by 0")
	}

	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

/* 使用tcp服务实现rpc */
func main() {
	//创建listener
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	arith := new(Arith)
	rpc.Register(arith) //和http流程一样,注册服务
	//使用tcp的listener来监听,实际上也是接受连接然后go协程处理
	rpc.Accept(l) //默认是gob编码器

}
