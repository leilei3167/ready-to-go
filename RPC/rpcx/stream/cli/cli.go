package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
)

var addr = flag.String("addr", "localhost:8972", "server address")

func main() {
	flag.Parse()

	d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr}, {Key: "localhost:8974"}})
	//需要额外创建一个客户端专门用于流式传输
	xclient := client.NewXClient(share.StreamServiceName, client.Failtry, client.RoundRobin, d, client.DefaultOption)
	defer xclient.Close()

	// get a connection for streaming
	/* 	conn, err := xclient.Stream(context.Background(), map[string]string{"开始流式传输": ""})
	   	if err != nil {
	   		panic(err)
	   	} */

	conns := []net.Conn{}

	for i := 0; i < 2; i++ {
		conn, err := xclient.Stream(context.Background(), map[string]string{"开始流式传输": ""})
		if err != nil {
			panic(err)
		}
		log.Printf("成功连接:%s", conn.RemoteAddr().String())
		time.Sleep(time.Second * 4)
		conns = append(conns, conn)
	}

	/* 	go func() {
		io.Copy(os.Stdout, conn)
		conn.Close()
	}() */

	file, err := os.Open("hello.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var index int
	for scanner.Scan() {
		if scanner.Text() != "" {
			x := index % 2
			fmt.Fprintln(conns[x], scanner.Text())
			fmt.Printf("已发送:%s To index:%d\n", scanner.Text(), x)
			index++
		}
		time.Sleep(time.Second)
	}
	/* reader := bufio.NewReader(file)
	io.Copy(conn, reader) */
	for _, v := range conns {
		v.Close()
	}
}
