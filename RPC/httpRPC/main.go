package main

import (
	"errors"
	"log"
	"net/http"
	"net/rpc"
)

//标准库rpc用法(使用http协议实现)
type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

//可以注册的方法必须满足:func (t *T) MethodName(argType T1, replyType *T2) error
//方法必须是可导出的;接受的2个参数必须是可导出的或者内置类型;第二个参数必须是指针;返回一个error
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

func main() {
	//1.基本rpc服务端
	arith := new(Arith)
	rpc.Register(arith) //注册Arith实例的方法
	rpc.HandleHTTP()    //注册rpc路由;:1234端口的访问交给rpc内部路由处理
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("serve error:", err)
	}
	/* rpc的HandleHTTP和ListenAndServe是如何联系起来的??
	   // src/net/rpc/server.go

	   func (server *Server) HandleHTTP(rpcPath, debugPath string) {
	   	http.Handle(rpcPath, server)
	   	http.Handle(debugPath, debugHTTP{server})
	   }

	   func HandleHTTP() {
	   	DefaultServer.HandleHTTP(DefaultRPCPath, DefaultDebugPath)
	   }

	   HandleHTTP实际上是把预定义的路径上（/_goRPC_）注册处理器,这个处理器会被加入到http服务的默认Mux上;除了/_goRPC_这个路径之外,还注册了一个/debug/rpc,服务开启后可在浏览器进入此页面
	*/

	//2.自定义注册的方法名
	//	rpc.RegisterName("math", arith) //默认方法以实例的名字作为方法名,这样自定义只会客户端调用需要改为math.Multiply

}
