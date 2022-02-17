package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

/*演示如何以'自定义的Client发出GET请求'*/

func main() {
	//1.配置transport
	transport := http.Transport{
		DisableKeepAlives: true,
		// 如果DisableKeepAlives为真，会禁止不同HTTP请求之间TCP连接的重用。

	}

	//2.配置client
	client := http.Client{Transport: &transport}

	//创建,配置Request
	req, err := http.NewRequest("GET", "http://httpbin.org/get", nil)
	if err != nil {
		panic(err)
	}
	//添加header信息,主要用于爬虫
	req.Header.Set("nimade", "www.leilei.com")

	//执行Request
	res, err := client.Do(req)
	fmt.Printf("res is :%v,err:%v,header is %v \n\n\n", res, err, req.Header)

	//解析body的内容
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	fmt.Println("res中的body为:", string(data))
	fmt.Println("------------------------------------------\n\n")

}
