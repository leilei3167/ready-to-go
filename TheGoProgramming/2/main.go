package main

import (
	"fmt"
)

func main() {
	//2.1变量声明
	var i, j, k int
	fmt.Println(i, j, k) //零值,var 适用于需要显示的指定变量的地方
	l := 2               //简短声明,应限制只在提高代码可读性的地方使用,如for
	fmt.Println(l)
	//创建指针变量
	var m *int
	m = &i
	fmt.Printf("%p\n", m)
	m = new(int)          //new函数=声明一个类型的零值,并返回该值的指针到调用处
	fmt.Printf("%p\n", m) //地址不相同
	//局部变量的生命周期是动态的,当某个函数中创建了变量
	//并返回了该变量的地址,则该变量将逃逸到堆上,

	//2.2赋值:元组赋值,同时更新多个变量的值
	a := func(n int) int {
		x, y := 0, 1
		for i := 0; i < n; i++ {
			x, y = y, x+y
		}
		return x
	}
	fmt.Println("斐波那契数列第几个数:", a(12))
	//隐式赋值:如函数调用时,实参隐式的赋值给形参,复合类型的字面量等
	models := []int{1, 2, 3, 4}
	fmt.Printf("%#v", models)
	/* 对于两个值是否可以用==或!=进行相等比较的能力也和可赋值能力有关系：
	对于任何类型的值的相等比较，
	第二个值必须是对第一个值类型对应的变量是可赋值的，反之亦然 */
}
