package main

import "fmt"

func main() {
	//len=0,cap=3
	base := make([]int, 0, 3)
	toAdd := []int{1, 2, 3, 4, 5}
	fmt.Printf("base:%#v len:%d,cap:%d\n", base, len(base), cap(base))
	fmt.Printf("toAdd:%#v len:%d,cap:%d\n", toAdd, len(toAdd), cap(toAdd))

	//超出base容量会重新分配一个新底层数组,不影响原base
	fmt.Println("====超出base的Cap,重新分配底层数组,原base不受影响=====")
	newInts := append(base, toAdd...)
	fmt.Printf("append newInts:%#v len:%d,cap:%d\n", newInts, len(newInts), cap(newInts))
	fmt.Printf("before append base:%#v len:%d,cap:%d\n", base, len(base), cap(base))

	//如果增加的没有超过容量?
	fmt.Println("====没有超出base的Cap=====")
	newNew := append(base, 1, 2, 3)
	fmt.Printf("append newInts:%#v len:%d,cap:%d\n", newNew, len(newNew), cap(newNew))
	fmt.Printf("before append base:%#v len:%d,cap:%d\n", base, len(base), cap(base))

	//append 赋值回去 将替换base为已扩容的底层数组
	fmt.Println("====赋值回base=====")
	base = append(base, toAdd...)
	fmt.Printf("append base:%#v len:%d,cap:%d\n", base, len(base), cap(base))

}
