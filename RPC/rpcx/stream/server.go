package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/smallnest/rpcx/share"

	"github.com/smallnest/rpcx/server"
)

var (
	addr             = flag.String("addr", "localhost:8972", "server address")
	fileTransferAddr = flag.String("transfer-addr", "localhost:8973", "data transfer address")
)

func main() {
	flag.Parse()

	s := server.NewServer()

	p := server.NewStreamService(*fileTransferAddr, streamhandler, nil, 1000)
	s.EnableStreamService(share.StreamServiceName, p)

	go func() {
		err := s.Serve("tcp", "localhost:8974")
		if err != nil {
			panic(err)
		}
	}()

	err := s.Serve("tcp", *addr)
	if err != nil {
		panic(err)
	}
}

func streamhandler(conn net.Conn, args *share.StreamServiceArgs) {
	fmt.Printf("received args, meta: %v\n", args.Meta)

	addr := conn.RemoteAddr().String()

	// still copy until the connection is closed
	io.Copy(os.Stdout, conn)

	conn.Close()
	fmt.Printf("%s closed\n", addr)
}
