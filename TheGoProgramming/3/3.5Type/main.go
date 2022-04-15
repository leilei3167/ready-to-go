package main

import "fmt"

//type 类型名字 底层类型
//类型声明一般出现在包一级
type Celsius float64    // 摄氏温度
type Fahrenheit float64 // 华氏温度

const (
	AbsoluteZeroC Celsius = -273.15 // 绝对零度
	FreezingC     Celsius = 0       // 结冰点温度
	BoilingC      Celsius = 100     // 沸水温度
)

//当底层类型相同时,可以进行类型转换,必须显式转换
//浮点数转为整数将丢失精度,字符串可以转化为字节切片等;这类转换可能改变值的表现,将一个字符串转为[]byte类型的slice将拷贝一个字符串数据的副本
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }

//底层数据类型决定了内部结构和表达方式，也决定是否可以像底层类型一样对内置运算符的支持。这意味着，Celsius和Fahrenheit类型的算术运算行为和底层的float64类型是一样的，正如我们所期望的那样。

func main() {
	var c Celsius
	var f Fahrenheit
	fmt.Println(c == 0) // "true"
	fmt.Println(f >= 0) // "true"
	//	fmt.Println(c == f)          // compile error: type mismatch
	fmt.Println(c == Celsius(f)) // "true"!
}

//命名类型可以提供书写方便,尤其是在处理复杂类型时(结构体等)
//命名类型还可以为该类型的值定义新的行为。这些行为表示为一组关联到该类型的函数集合，我们称为类型的方法集。我们将在第六章中讨论方法的细节，这里只说些简单用法。