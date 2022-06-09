package main

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func main() {

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(2),
	)
	extensions.RandomUserAgent(c)

	//使用代理切换
	/* 	rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:10808", "http://127.0.0.1:10809")
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	c.SetProxyFunc(rp) */

	c.Limit(&colly.LimitRule{
		Parallelism: 100, //并发量
		DomainGlob:  "*",
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println(link)
		e.Request.Visit(link)

	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("code:", r.StatusCode)
	})

	c.Visit("https://www.bilibili.com/")
	c.Wait() //并发模式下必须使用wait

}
