package sample

import (
	"fmt"
	"testing"
)

//1.基准测试以Benchmark为开头,否则编译器将无法识别出
//2.参数必须是b *testing.B
//例:func BenchmarkDownload(b *testing.B) {}

//测试pirnt哪种更快

var gs string
var a []int

func BenchmarkPrint(b *testing.B) {
	var s string
	//	a = append(a, b.N)
	for i := 0; i < b.N; i++ {
		s = fmt.Sprint("hello")
	}
	/* 	if len(a) > 4 {
		fmt.Println(a) //[1 100 10000 1000000 63441712]
	} */

	gs = s
}

func BenchmarkPrintf(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = fmt.Sprintf("hello")
	}
	gs = s
}

// go test -bench . -benchtime 3s -benchmem
//-bench . 运行所有的benchmark; -benchtime 3s 将测试时间延长到3s; -benchmem 显示内存分配信息
//Sprintf比不格式化的版本更快!

//基准测试的核心在于循环,b.N,默认的基准测试时间是1s,而循环是不能以时间作为标的的,会有一个算法不断的修正b.N以接近设置的时间
//非常重要的一点就是 循环内的代码就是实际被测试的代码,编译器会将其中的代码编译,里面的代码和实际生产中保持一致是非常重要的

/* 二,子基准测试*/
func BenchmarkSprint(b *testing.B) {
	b.Run("none", benchSprint)
	b.Run("format", benchSprintf)
}

func benchSprintf(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = fmt.Sprintf("hello")
	}
	gs = s
}

func benchSprint(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = fmt.Sprint("hello")
	}
	gs = s
}

/* 三,验证benchmark
基准测试的结果可能不一定会是准确的;

*/

var n []int

func init() {
	for i := 0; i < 1_000_000; i++ {
		n = append(n, i)
	}
}
func BenchmarkSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		single(n)
	}
}
func BenchmarkUnlimited(b *testing.B) {
	for i := 0; i < b.N; i++ {
		unlimited(n)
	}
}
func BenchmarkNumCPU(b *testing.B) {
	for i := 0; i < b.N; i++ {
		numCPU(n, 0)
	}
}

//多个基准测试直接会互相干扰,一定要确保进行benchmark时 主机要尽可能的空闲
