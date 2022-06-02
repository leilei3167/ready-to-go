package service

import (
	"golang.org/x/net/context"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (c *Consumer) ToScan() {
	for task := range c.TaskChan {
		start := time.Now()
		var r Result

		var wg sync.WaitGroup
		r.IP = task.IP
		r.mu = new(sync.Mutex)
		wg.Add(len(task.Ports))

		for _, p := range task.Ports {
			c.Sem.Acquire(context.TODO(), 1)
			go func(h string, p int) {
				defer wg.Done()
				host := net.JoinHostPort(r.IP, strconv.Itoa(p))

				conn, err := net.DialTimeout("tcp", host, time.Second*1)
				if err != nil {
					if strings.Contains(err.Error(), "too many open files") {
						log.Fatal(err)
					}
					//log.Println(err)
					c.Sem.Release(1)
					return
				}
				log.Printf("%v is open", host)
				conn.Close()
				c.Sem.Release(1)

				r.mu.Lock()
				r.OpenPorts = append(r.OpenPorts, p)
				r.mu.Unlock()
				return
			}(r.IP, p)
		}

		wg.Wait()
		log.Printf("%v扫描完毕,共%d个端口,其中打开端口%d个,耗时:%v", task.IP, len(task.Ports), len(r.OpenPorts), time.Since(start))
		c.ResultChan <- r //data race with close(ResultChan)
	}

}
