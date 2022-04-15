package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	//收集结果
	ch := make(chan string)
	//并发get
	for _, v := range os.Args[1:] {
		go fetch(v, ch)

	}
	//阻塞获取传入url数量那么多的结果
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) //将错误发送到chan
		return
	}
	//返回拷贝了多少字节
	//discard可以把这个变量看作一个垃圾桶，可以向里面写一些不需要的数据
	//此处不想要他的结果,只想要他读取的字节数
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	//since从开始到现在
	secs := time.Since(start).Seconds()
	//保留2位小数,将结果拼成字符串传回
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)

}
