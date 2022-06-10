package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
)

func main() {
	ips := []string{"http://182.61.25.124:80", "182.61.25.124:443"}

	c := colly.NewCollector(colly.MaxDepth(2), colly.Debugger(&debug.LogDebugger{}))
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: time.Second * 3,
	})
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	c.WithTransport(
		&http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns: 100,

			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.Request.URL.String())
		fmt.Println(r.StatusCode)
	})

	c.Visit(ips[0])

}
