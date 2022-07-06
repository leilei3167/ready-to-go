package main

import "fmt"

//concul zookeeper rpcx

func main() {
	// Initialize a map for the integer values
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}

	// Initialize a map for the float values
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	fmt.Printf("Non-Generic Sums: %v and %v\n",
		SumInts(ints),
		SumFloats(floats))

	//SumIntsOrFloats[string, int64](ints),[]中内容可省略,编译器可自动推断
	//如果要调用一个没有参数的泛型函数,则需要带上类型参数
	fmt.Printf("Generic Sums: %v and %v\n",
		SumIntsOrFloats(ints),
		SumIntsOrFloats(floats))

}

/* 没有泛型，以下两个函数区别仅仅在于一个是int一个是float，都需要写2遍 */

// SumInts adds together the values of m.
func SumInts(m map[string]int64) int64 {
	var s int64
	for _, v := range m {
		s += v
	}
	return s
}

// SumFloats adds together the values of m.
func SumFloats(m map[string]float64) float64 {
	var s float64
	for _, v := range m {
		s += v
	}
	return s
}

/* 泛型函数 */
//[]内为类型参数,K 为可比较的参数（因为map的key必须为可比的，V可以是int64或float64，|表示合集）
//输入值为map[k]v,返回值为v的函数
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

//可以将类型约束写入接口来简化代码
type Number interface {
	int64 | float64
}

func SumNumbers[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}
