package main

import "log"

func main() {

	a := make([]int, 10)
	for i := 0; i < cap(a); i++ {
		a[i] = i
	}
	log.Printf("---填满后:---a:%#v,len(a):%d,cap(a):%d", a, len(a), cap(a))

	a = append(a, 1000)[:len(a)] //触发扩容
	log.Printf("---扩容后:---a:%#v,len(a):%d,cap(a):%d", a, len(a), cap(a))


}
