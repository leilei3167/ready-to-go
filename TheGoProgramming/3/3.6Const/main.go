package main

import "fmt"

//常量的值在编译期间就会确定,常量不可取地址
//常量的声明可以包含类型也可以省去,省去类型将会根据右边的类型进行推断
type Weekday int

//在其他语言中被称之为枚举类型
const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type Currency int //指定类型

const (
	USD Currency = iota // 美元 0
	EUR                 // 欧元 1
	GBP                 // 英镑 2
	RMB                 // 人民币 3
)

//无类型常量
//Go语言的常量有个不同寻常之处。虽然一个常量可以有任意一个确定的基础类型，例如int或float64，或者是类似time.Duration这样命名的基础类型，但是许多常量并没有一个明确的基础类型。编译器为这些没有明确基础类型的数字常量提供比基础类型更高精度的算术运算；你可以认为至少有256bit的运算精度

func main() {
	fmt.Println(RMB) //a[1] = 12
}
