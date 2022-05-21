package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

//发起rpc调用所需的参数
type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func main() {
	//创建rpc的客户端
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("------------同步调用---------------")
	time.Sleep(time.Second)
	//构建参数
	args := &Args{7, 8}
	var reply int
	//调用,指明实例的名字.方法
	/* Call是同步调用,会阻塞等待服务端的响应 */
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Got resp:%d*%d=%d\n", args.A, args.B, reply)

	args = &Args{15, 6}
	var quo Quotient
	err = client.Call("Arith.Divide", args, &quo)
	if err != nil {
		log.Fatal("Divide error:", err)
	}
	fmt.Printf("Divide: %d/%d=%d...%d\n", args.A, args.B, quo.Quo, quo.Rem)

	/* 2.异步调用
		Go方法会启动异步调用,会返回一个rpc.Call结构体

		type Call struct {
	  ServiceMethod string
	  Args          interface{}
	  Reply         interface{}
	  Error         error
	  Done          chan *Call
	}

	我们应该通过这一个对象来获取此次调用的信息,如对象参数 返回值和错误等;我们通过监听其Done通道来判断该调用是否完成




	*/
	log.Println("------------异步调用--------------")
	time.Sleep(time.Second)
	args1 := &Args{7, 8}
	var reply1 int
	mulResp := client.Go("Arith.Multiply", args1, &reply1, nil)

	args2 := &Args{15, 6}
	var quo1 Quotient
	divideReply := client.Go("Arith.Divide", args2, &quo1, nil)

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	//select监听各个调用的Done字段来获取结果
	var multiplyReplied, divideReplied bool
	for !multiplyReplied || !divideReplied {
		select {
		case resplyCall := <-mulResp.Done:
			if err := resplyCall.Error; err != nil { //先判断错误字段是否未nil
				fmt.Println("Multiply error:", err)
			} else {
				fmt.Printf("Multiply异步调用结果: %d*%d=%d\n", args1.A, args1.B, reply)

			}
			multiplyReplied = true
		case replyCall := <-divideReply.Done:
			if err := replyCall.Error; err != nil {
				fmt.Println("Divide error:", err)
			} else {
				fmt.Printf("Divide异步调用结果: %d/%d=%d...%d\n", args2.A, args2.B, quo.Quo, quo.Rem)
			}
			divideReplied = true
		case <-ticker.C:
			fmt.Println("tick")
		}

	}

}
