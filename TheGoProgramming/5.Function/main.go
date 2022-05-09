package main

import (
	"fmt"
	"log"
)

//函数的类型被称为函数的签名。如果两个函数形式参数列表和返回值列表中的变量类型一一对应，那么这两个函数被认为有相同的类型或签名,

//在Go中，函数被看作第一类值（first-class values）,函数像其他值一样,拥有类型,可以赋值给变量,作为函数参数,返回值等;

func f1(x, y int) int {
	return x * y
}
func f2(z, w int) (s int) {
	return z * w
}
func f3(fn func(a int, b int) int) { //x,y int 和x int y int相同
	log.Println("Yeah!")
}
func f4(a, b int) int {
	return a * b
}

func main() {
	a := f1 //函数赋值给a,a为函数值
	b := f2
	f3(a)
	f3(b)
	f4(a(1, 2), b(2, 4)) //函数值调用和调用对应的函数一致
	//函数类型零值是nil
	var c func() int
	fmt.Printf("type:%T,value %v\n", c, c) //type:func() int,value <nil>,调用会panic
	//函数值可以和nil比较,但函数值之间不能比较,因此函数值不能作为map的key

	//--------形参和返回值的变量名不影响函数签名，也不影响它们是否可以以省略参数类型的形式表示

}
