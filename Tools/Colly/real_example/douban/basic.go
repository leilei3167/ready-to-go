package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func fetch(url string) string {
	fmt.Println("Fetch Url", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
		return ""
	}
	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error", err)
		return ""
	}
	return string(body)
}

func parseUrls(url string) {
	//直接将首页下载下来
	body := fetch(url)
	//使用正则表达式提取数据
	body = strings.Replace(body, "\n", "", -1)
	rp := regexp.MustCompile(`<div class="hd">(.*?)</div>`)                         //标题所在的区域
	titleRe := regexp.MustCompile(`<span class="title">(.*?)</span>`)               //提取标题
	idRe := regexp.MustCompile(`<a href="https://movie.douban.com/subject/(\d+)/"`) //提取电影编号
	items := rp.FindAllStringSubmatch(body, -1)                                     //找到标题存在的元素列表
	for _, item := range items {
		//打印对应的编号,标题
		fmt.Println(idRe.FindStringSubmatch(item[1])[1],
			titleRe.FindStringSubmatch(item[1])[1])
	}
}
func main() {
	start := time.Now()
	//10页
	for i := 0; i < 10; i++ {
		parseUrls("https://movie.douban.com/top250?start=" + strconv.Itoa(25*i))
	}
	elapsed := time.Since(start)
	fmt.Printf("Took %s", elapsed)
}
