package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func main() {
	url := "https://httpbin.org/delay/2"
	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(2),
	)
	//使用随机的User-Agent和随机的延迟来伪装
	extensions.RandomUserAgent(c)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 3,
		RandomDelay: 5 * time.Second, //随机延迟,最长5秒
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("准备发送请求:", r.Headers.Get("User-Agent"))
	})

	for i := 0; i < 4; i++ {
		c.Visit(fmt.Sprintf("%s?n=%d", url, i))
	}
	c.Visit(url)
	c.Wait()
}
