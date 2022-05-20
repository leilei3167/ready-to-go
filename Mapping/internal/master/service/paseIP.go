package service

import (
	"bufio"
	"io"
	"mapping/internal/pkg/code"
	"sync"
)

type IPparser struct {
	file       io.Reader   //数据源
	IPChan     chan string //解析成功的结果
	InvalidIPs chan string //无效的IP列表
	//Ctx        context.Context
}

func NewIPparser(f io.Reader) *IPparser {
	return &IPparser{
		file:       f,
		IPChan:     make(chan string, 10000),
		InvalidIPs: make(chan string, 100),
	}
}
func (i *IPparser) Close() { //文件句柄由调用处关闭
	close(i.IPChan)
	close(i.InvalidIPs)

}
func (i *IPparser) ReadAndParse() {
	ch := make(chan string, 1)
	scanner := bufio.NewScanner(i.file)

	var wg sync.WaitGroup
	wg.Add(10)
	for b := 0; b < 10; b++ { //10个解析器
		go validAndParse(ch, i.IPChan, i.InvalidIPs, &wg)
	}

	for scanner.Scan() {
		if scanner.Text() != "" {
			ch <- scanner.Text()
		}
	}
	close(ch) //读取完毕直接关闭,使得解析器能够退出
	wg.Wait()
	i.Close()
}

func validAndParse(ch <-chan string, IPchan chan<- string, InvalidIPs chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {
		ips := code.ParseMixIP(v)
		if ips != nil { //解析成功
			for _, ip := range ips {
				IPchan <- ip
			}
		} else {
			InvalidIPs <- v //解析失败
		}
	}

}
