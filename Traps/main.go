package main

import "sync"

func main() {
	rangeclosure()
}

//for range中开启协程造成闭包的问题
func rangeclosure() {
	wg := sync.WaitGroup{}
	si := []int{1, 2, 3, 5, 6}
	for i := range si {
		wg.Add(1)
		go func() {
			println(i) //形成闭包,因为引用了外部变量i
			wg.Done()
		}()
	}
	wg.Wait()

}
