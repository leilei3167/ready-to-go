package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	var pool = make(chan bool, 10000)
	var begin = time.Now()
	//wg
	var wg sync.WaitGroup
	//ip
	var ip = "127.0.0.1"
	//var ip = "192.168.43.34"
	//循环
	for j := 21; j <= 65535; j++ {
		//添加wg
		wg.Add(1)
		go func(i int) {
			pool <- true

			//释放wg
			defer wg.Done()
			var address = fmt.Sprintf("%s:%d", ip, i)
			conn, err := net.DialTimeout("tcp", address, time.Second*5)
			//conn, err := net.Dial("tcp", address)
			if err != nil {
				//fmt.Println(address, "是关闭的", err)
				<-pool
				return
			}
			conn.Close()
			fmt.Println(address, "打开")
			<-pool
		}(j)
	}
	//等待wg
	wg.Wait()

	var elapseTime = time.Now().Sub(begin)
	fmt.Println("耗时:", elapseTime)
}
