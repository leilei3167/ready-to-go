package main

import "fmt"

func main() {
	c := a() //c此刻为闭包
	d := a()
	//同一函数构建出来的多个闭包所引用的外部变量是多个副本,彼此独立(但如果引用的是全局变量,则始终都为同一个)
	//调用同一闭包多次,将对其引用的外部变量造成影响
	//编译器检测到闭包时会将其引用的外部变量分配到堆上
	c()
	c()
	c()
	c() //输出4
	d() //输出1
}

//构建一个闭包
func a() func() int {
	i := 0
	b := func() int {
		i++
		fmt.Println(i)
		return i

	}
	return b
}
