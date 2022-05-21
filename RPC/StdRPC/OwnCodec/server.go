package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/rpc"
)

//rpc默认使用gob编解码,基本只能在go程序间实现调用,要使用自定义的编解码器,需实现 ServerCodec接口
/* // src/net/rpc/server.go
type ServerCodec interface {
  ReadRequestHeader(*Request) error
  ReadRequestBody(interface{}) error
  WriteResponse(*Response, interface{}) error

  Close() error
}
直接仿造其中gob编解码器实现一个即可


*/
//和gob解码器一模一样几乎
type JsonServerCodec struct {
	rwc    io.ReadWriteCloser
	dec    *json.Decoder
	enc    *json.Encoder
	encBuf *bufio.Writer
	closed bool
}

func NewJsonServerCodec(conn io.ReadWriteCloser) *JsonServerCodec {
	buf := bufio.NewWriter(conn)
	return &JsonServerCodec{conn, json.NewDecoder(conn), json.NewEncoder(buf), buf, false}
}

func (c *JsonServerCodec) ReadRequestHeader(r *rpc.Request) error {
	return c.dec.Decode(r)
}

func (c *JsonServerCodec) ReadRequestBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *JsonServerCodec) WriteResponse(r *rpc.Response, body interface{}) (err error) {
	if err = c.enc.Encode(r); err != nil {
		if c.encBuf.Flush() == nil {
			log.Println("rpc: json error encoding response:", err)
			c.Close()
		}
		return
	}
	if err = c.enc.Encode(body); err != nil {
		if c.encBuf.Flush() == nil {
			log.Println("rpc: json error encoding body:", err)
			c.Close()
		}
		return
	}
	return c.encBuf.Flush()
}

func (c *JsonServerCodec) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}

func main() {
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	arith := new(Arith)
	rpc.Register(arith)
	//要使用自定义的Codec就不能直接Accept(l)了,需自己获取连接,并调用ServeCodec
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		go rpc.ServeCodec(NewJsonServerCodec(conn))
	}
}

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
