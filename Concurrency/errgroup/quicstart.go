package main

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {
	//和普通的wg声明类似
	//此处未使用Context,这个Group中的协程出现错误后将不能取消当前
	//其他的协程
	var g errgroup.Group
	var urls = []string{
		"http://www.baidu.com/",
		"http://www.baidu.com/",
		"http://www.1234567.com/", //假的
	}
	for _, url := range urls {
		// Launch a goroutine to fetch the URL.
		url := url
		//func (g *Group) Go(f func() error)
		g.Go(func() error { //内部集成了done
			// Fetch the URL.
			resp, err := http.Get(url)
			if err == nil { // 这里记得关掉
				resp.Body.Close()
			}
			return err
		})
	}
	// Wait for all HTTP fetches to complete.
	err := g.Wait()
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	fmt.Println("Successfully fetched all URLs.")
}
