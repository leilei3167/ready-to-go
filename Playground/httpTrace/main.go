package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
)

func main() {
	req, _ := http.NewRequest(http.MethodGet, "http://www.baidu.com", nil)

	//创建追踪器
	clientTrace := &httptrace.ClientTrace{
		GetConn:      func(hostPort string) { fmt.Println("starting to create conn ", hostPort) },
		DNSStart:     func(info httptrace.DNSStartInfo) { fmt.Println("starting to look up dns", info) },
		DNSDone:      func(info httptrace.DNSDoneInfo) { fmt.Println("done looking up dns", info) },
		ConnectStart: func(network, addr string) { fmt.Println("starting tcp connection", network, addr) },
		ConnectDone:  func(network, addr string, err error) { fmt.Println("tcp connection created", network, addr, err) },
		GotConn:      func(info httptrace.GotConnInfo) { fmt.Printf("connection established:%+v", info) },
	}

	//必须将追踪器设置到req的请求contex中,这样才能被传递,否则不会被追踪
	clientTraceCtx := httptrace.WithClientTrace(req.Context(), clientTrace)
	req = req.WithContext(clientTraceCtx)

	_, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}
