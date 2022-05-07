package main

import (
	"context"
	"log"
	"time"

	"github.com/smallnest/rpcx/client"
)

type Arg struct {
	Name string
}
type Res struct {
	Msg string
}

func main() {

	d, _ := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:8080", "")
	xclient := client.NewXClient("test", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	a := &Arg{
		Name: "just for test",
	}
	b := &Res{}
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
	defer cancle()

	err := xclient.Call(ctx, "Try", a, b)
	if err != nil {
		log.Println("调用出错!", err)
		return
	}
	log.Println("调用成功", b)

}
