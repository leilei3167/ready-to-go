package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	for {
		var in string
		fmt.Println("输入值")
		fmt.Scanln(&in)

		params := url.Values{}
		Url, err := url.Parse("http://localhost:8080/echo")
		if err != nil {
			panic(err.Error())
		}

		params.Set("input", in)

		//如果参数中有中文参数,这个方法会进行URLEncode
		Url.RawQuery = params.Encode()
		urlPath := Url.String()

		resp, err := http.Get(urlPath)
		if err != nil {
			fmt.Println("Get出错:", err)
		}
		defer resp.Body.Close()
		s, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取响应错误", err)
		}
		fmt.Println("返回结果", string(s))
	}
}
