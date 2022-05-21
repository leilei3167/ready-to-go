package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
)

//客户端中同样需要定义编解码器,实现rpc.ClientCodec接口,同一仿造一个

type JsonClientCodec struct {
	rwc    io.ReadWriteCloser
	dec    *json.Decoder
	enc    *json.Encoder
	encBuf *bufio.Writer
}

func NewJsonClientCodec(conn io.ReadWriteCloser) *JsonClientCodec {
	encBuf := bufio.NewWriter(conn)
	return &JsonClientCodec{conn, json.NewDecoder(conn), json.NewEncoder(encBuf), encBuf}
}

func (c *JsonClientCodec) WriteRequest(r *rpc.Request, body interface{}) (err error) {
	if err = c.enc.Encode(r); err != nil {
		return
	}
	if err = c.enc.Encode(body); err != nil {
		return
	}
	return c.encBuf.Flush()
}

func (c *JsonClientCodec) ReadResponseHeader(r *rpc.Response) error {
	return c.dec.Decode(r)
}

func (c *JsonClientCodec) ReadResponseBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *JsonClientCodec) Close() error {
	return c.rwc.Close()
}

func main() {
	//先使用net进行TCP连接,拿到conn后创建客户端
	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal("dial error:", err)
	}
	//Client创建就需要使用NewClientWithCodec加上自定义的编解码器
	client := rpc.NewClientWithCodec(NewJsonClientCodec(conn))

	args := &Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("Multiply error:", err)
	}
	fmt.Printf("Multiply: %d*%d=%d\n", args.A, args.B, reply)
}

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}
