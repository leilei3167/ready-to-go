package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func main() {
	host := "http://182.61.25.124"
	urls, err := url.Parse(host)
	if err != nil {
		log.Fatal(err)
	}

	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 3 * time.Second, //建立连接的超时
		}).DialContext,
		DisableKeepAlives: true, //告知服务端不启用长连接,避免占用端口资源
		MaxIdleConns:      100,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, //跳过证书验证,确保自创证书的网站也能被收集
	}

	cli := http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Printf("[即将发出的请求]:%#v\n", req.URL.String())
			for _, via1 := range via {
				fmt.Printf("[之前已执行的请求]:%#v\n", via1.URL.String())
			}
			return nil
		},
	}

	req, err := http.NewRequest("GET", urls.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[Req] URL:%#v\n", req.URL)
	resp, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[Resp] :%#v\n", resp.StatusCode)

}
