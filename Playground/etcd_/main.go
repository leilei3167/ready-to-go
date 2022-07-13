package main

import (
	"fmt"
	"time"

	"github.com/rpcxio/libkv/store"
	etcdV3 "github.com/rpcxio/rpcx-etcd/client"
	"github.com/sirupsen/logrus"
)

func main() {
	etcdConfigOption := &store.Config{
		ClientTLS:         nil,
		TLS:               nil,
		ConnectionTimeout: 5 * time.Second,
		Bucket:            "",
		PersistConnection: true,
	}
	d, err := etcdV3.NewEtcdV3Discovery(
		"/mapping",
		"LogicRpc",
		[]string{"127.0.0.1:2379"},
		true,
		etcdConfigOption,
	)
	if err != nil {
		logrus.Fatal("创建etcd服务发现错误:", err)
	}
	kv := d.GetServices()

	for _, v := range kv {
		fmt.Printf("%s:%s\n", v.Key, v.Value)
	}

	//	fmt.Printf("%#v,Len:%d\n", d.GetServices(), len(d.GetServices()))

}
