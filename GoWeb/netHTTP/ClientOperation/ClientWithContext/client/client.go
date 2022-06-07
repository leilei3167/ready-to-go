package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type respData struct {
	Response *http.Response
	err      error
}

func doCall(ctx context.Context) {

	//配置transport
	transport := http.Transport{

		DisableKeepAlives: true,
	}
	client := http.Client{
		Transport: &transport,
	}

	//2.创建数据通道,用协程处理来发请求数据,以及定义请求类型(会返回*Request)
	respChan := make(chan *respData, 1)
	//NewRequest使用指定的方法、网址和可选的主题创建并返回一个新的*Request。
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		fmt.Printf("new requestg failed, err:%v\n", err)
		return

	}

	//3.更新成为创建带上下文的请求
	//使用带超时的ctx创建一个新的client Request
	//func (r *Request) WithContext(ctx context.Context) *Request
	req = req.WithContext(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait() //等协程处理完才会退出函数

	//4.执行请求(会得到Response 和err)
	go func() {
		//func (c *Client) Do(req *Request) (*Response, error)
		//Do方法发送请求，返回HTTP回复。
		//它会遵守客户端c设置的策略（如重定向、cookie、认证）。
		resp, err := client.Do(req)
		fmt.Printf("client.do resp:%v, err:%v\n", resp, err)
		//将response的内容取出并用通道传递
		rd := &respData{
			Response: resp,
			err:      err,
		}
		respChan <- rd
		wg.Done()

	}()

	//5.select阻塞直到超时退出或者成功
	select {

	case <-ctx.Done():
		fmt.Println("链接超时!!!")
	case result := <-respChan: //监听是否成功的到response数据
		fmt.Println("链接成功!")
		if result.err != nil {
			fmt.Printf("call server api failed, err:%v\n", result.err)
			return
		}

		defer result.Response.Body.Close() //延迟关闭
		data, _ := ioutil.ReadAll(result.Response.Body)
		fmt.Println(string(data))

	}

}

func main() {
	//定义一个超时关闭
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*1)
	defer cancle()
	doCall(ctx)

}
