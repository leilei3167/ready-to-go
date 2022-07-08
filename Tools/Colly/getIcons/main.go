package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/sirupsen/logrus"
)

type WebData struct {
	Url           string            `json:"url,omitempty" bson:"url,omitempty"`
	Status        string            `json:"status,omitempty" bson:"status,omitempty"`
	ContentLength string            `json:"contentlength,omitempty" bson:"contentlength,omitempty"`
	Title         string            `json:"title,omitempty"bson:"title,omitempty"` //网页标题
	Icon          map[string]string `json:"icon,omitempty"bson:"icon,omitempty"`   //网页图标, map[string]string ,"fluid-icon":"https://github.com/fluidicon.png"
}

func main() {
	url := "https://jahir.dev"
	info, err := catchWebInfo(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[结果]:icon:%#v\nurl:%s\n", info.Icon, info.Url)
}

func catchWebInfo(absURI string) (WebData, error) {
	web := WebData{Icon: make(map[string]string)}
	c := colly.NewCollector()

	extensions.RandomUserAgent(c) //使用随机的Agent
	extensions.Referer(c)

	c.WithTransport(
		&http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second, //建立连接的超时
			}).DialContext,
			DisableKeepAlives: true, //告知服务端不启用长连接,避免占用端口资源
			MaxIdleConns:      100,
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, //跳过证书验证,确保自创证书的网站也能被收集
		})

	c.OnHTML("title", func(e *colly.HTMLElement) { //抓取网站的title
		web.Title = e.Text
		web.Url = e.Request.URL.String()
		web.Url = strings.TrimRight(web.Url, "/")

		web.Status = strconv.Itoa(e.Response.StatusCode)
		web.ContentLength = strconv.Itoa(len(e.Response.Body))

	})
	//选择出rel属性包含'icon'关键字的所有link元素,并提取icon的url
	c.OnHTML("head", func(e *colly.HTMLElement) {
		//TODO:抓取icon不准确,并且需要将icon下载下来而不是存url
		//一般网站将favicon.ico放在根目录,会被浏览器自动识别成icon,此类做法无法通过爬虫爬取标签获得
		//因此,为尽可能的确保抓取icon的准确性,应将 hostname/favicon.ico 作为默认的icon地址,再进行html icon标签地址的寻找

		e.ForEach("[rel*=icon]", func(i int, s *colly.HTMLElement) {

			icon := s.Attr("href")
			icon = strings.TrimSpace(icon)

			iconUrl, err := url.Parse(icon)
			if err != nil {
				logrus.Warnf("解析iconUrl出错:%s err:%v", icon, err)
				return
			}
			log.Printf("Got [%s]:%s", s.Attr("rel"), iconUrl.String())
			//如果是绝对路径的url,直接写入
			if iconUrl.IsAbs() {
				if _, ok := web.Icon[s.Attr("rel")]; !ok {
					web.Icon[s.Attr("rel")] = iconUrl.String()
					return
				} else {
					return
				}
			}

			//是相对路径,与Host进行拼接,如 /images/favicon.ico ,scheme+host+icon
			reqUrl := web.Url + iconUrl.Path

			log.Printf("拼接之后的相对icon:%s", reqUrl)

			if _, ok := web.Icon[s.Attr("rel")]; !ok {
				web.Icon[s.Attr("rel")] = reqUrl
				return
			} else {
				return
			}

		})
	})
	c.OnError(func(r *colly.Response, err error) {
		web.Status = strconv.Itoa(r.StatusCode)
		web.Url = r.Request.URL.String()
		if len(r.Body) > 0 {
			web.ContentLength = strconv.Itoa(len(r.Body))
			//	web.MetaData = string(r.Body)
		}

	})
	/*	c.RedirectHandler = func(req *http.Request, via []*http.Request) error {
		//最近的一次跳转
		lastRequest := via[len(via)-1]
		// If domain has changed, remove the Authorization-header if it exists
		if req.URL.Host != lastRequest.URL.Host {
			web.Url = req.URL.String() //设置为跳转的目标url
		}
		return nil
	}*/

	err := c.Visit(absURI)
	if err != nil {
		logrus.Warnf("web探测%s出错:%v", absURI, err)
	}
	if web.Title == "" {
		web.Title = "none"
	}

	//如果没有任何的icon被发现,则设置默认的icon
	if len(web.Icon) == 0 {
		web.Icon["default_Icon"] = web.Url + "/favicon.ico"
	}

	return web, nil
}
