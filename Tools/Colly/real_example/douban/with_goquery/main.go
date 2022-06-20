package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

/* 一般查找匹配页面时不会用正则,因为可读性较差
goquery是jQuery的Golang版本实现。借用jQueryCSS选择器的语法可以非常方面的实现内容匹配和查找

*/
type bookInfo struct {
	name        string
	description string
}

func (b bookInfo) String() string {
	return fmt.Sprintf("书名:%s\n简介:%s\n", b.name, b.description)
}

func main() {
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	c.Limit(&colly.LimitRule{RandomDelay: time.Second * 3})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})
	/* 获取icon */
	c.OnHTML("head", func(e *colly.HTMLElement) { //head元素范围内
		//Text() 是获取该节点的文本值,如 <title>【置顶】大撒大撒</title>
		fmt.Println("title", e.DOM.Find("title").Text()) //直接获取head中title元素的值
		//找到所有link元素
		e.DOM.Find("link[rel*=icon]").Each(func(i int, s *goquery.Selection) {
			//遍历打印所有的href的属性值
			text, ok := s.Attr("href") //Attr是获取属性的值
			if ok {
				urll := e.Request.URL.String()
				//TrimRight会去除右边所有的目标字符串,而TrimeSuffix则只去除一个
				urll = strings.TrimRight(urll, "/")
				if strings.HasPrefix(text, "https://") || strings.HasPrefix(text, "http://") {
					fmt.Println("找到icon地址:", text)
				} else {
					fmt.Println("找到icon地址:", urll+text)
				}
			} else {
				fmt.Println("没有找到icon")
			}
		})
	})

	/* 获取其他信息 */
	c.OnHTML("body", func(e *colly.HTMLElement) {
		dom := e.DOM
		dom.Find("div[class=entry_box]").Each(func(i int, s *goquery.Selection) {
			//每一本书
			var book bookInfo
			name := s.Find("h3").Text()
			name = strings.Trim(name, "\n")
			book.name = name
			des := s.Find("div.entry_post_excerpt").Text()
			des = strings.Trim(des, "\n")
			book.description = des
			fmt.Printf("%v\n", book)
		})

	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished")
	})

	c.Visit("https://www.werebook.com/book-category/software")
}
