package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, v := range os.Args[1:] { //获取命令行参数
		if !strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
			//默认添加个https://
			v = "https://" + v
		}

		resp, err := http.Get(v)
		if err != nil {
			log.Printf("failed to get :%v", err)
			os.Exit(1)
		}
		//	data, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("status CODE:%v\n", resp.Status)
		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			log.Printf("failed to read body :%v", err)
			os.Exit(1)
		}
		//直接拷贝到标准输出,不需要缓冲区(data)

		//fmt.Printf("got result:%v", string(data))
	}

}
