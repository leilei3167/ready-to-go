package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main() {
	//创建爬虫对象,限制只爬取指定域名,避免无限跳转
	c := colly.NewCollector(colly.AllowedDomains("www.baidu.com"))
	//注册回调,对每个有href属性的a元素执行回调函数
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href") //返回当前元素的属性
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link)) //继续访问其指向的其他连接
	})
	//每次发送请求时执行该回调，这里只是简单打印请求的 URL
	c.OnRequest(func(r *colly.Request) {
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
	//访问第一个页面
	c.Visit("http://www.baidu.com/")
	/*
		colly爬取到页面后,会用goquery解析页面,查找注册了回调函数对应的元素选择器并执行

	*/

}
