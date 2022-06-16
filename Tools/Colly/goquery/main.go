package main

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func main() {

	c := colly.NewCollector()
	extensions.RandomUserAgent(c) //使用随机的user-agent
	extensions.Referer(c)

	c.OnHTML("body", func(h *colly.HTMLElement) {

	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("任务完成!")
	})

}
