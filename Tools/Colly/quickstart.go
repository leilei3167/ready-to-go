package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
)

func main() {
	//创建爬虫对象,并进行配置(限制跳转域名,使用随机agent,跳转时的referer,maxdepth为1则不进行跳转)
	c := colly.NewCollector(colly.AllowedDomains("www.baidu.com"), colly.MaxDepth(1),
		colly.Debugger(&debug.LogDebugger{}), //加入一个debugger,可实现其接口来自定义
	)
	extensions.RandomUserAgent(c) //使用随机的user-agent
	extensions.Referer(c)         //从callback中的visit才有效,访问子页面时带上refer,意思是从那里点击过来的
	/* 	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2, //每一个域名,最大2个并发
	}) */

	//默认使用的go标准库的默认client,可以实现自定义配置的client(RoundTriper是一个接口,transport就实现了它)
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	//注册回调,对每个有href属性的a元素执行回调函数
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href") //返回当前元素的属性
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		e.Request.Visit(link) //继续访问其指向的其他连接
	})
	c.OnHTML("title", func(e *colly.HTMLElement) {
		link := e.Attr("title") //返回当前元素的属性
		fmt.Printf("title found: %q -> %s\n", e.Text, link)

	})
	//每次发送请求时执行该回调，这里只是简单打印请求的 URL
	c.OnRequest(func(r *colly.Request) {
		//比如每个请求前做一些事
		fmt.Println("Visiting", r.URL.String())
	})
	//每次收到响应时执行该回调，这里也只是简单的打印 URL 和响应大小
	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Response %s: %d bytes\n", r.Request.URL, len(r.Body))
	})
	//错误时执行
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error %s: %v\n", r.Request.URL, err)
	})
	//在抓取完一个HTML页面时执行
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("抓取完毕,do something...")
	})

	//访问第一个页面
	c.Visit("http://www.baidu.com/")
	/*
	   关于分布式爬虫:
	   	大多数时候 网络层面的分布式就足够了(加代理),colly中可以很方便的实现代理的切换


	*/

}
