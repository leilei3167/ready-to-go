package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	bubblesort(randNum())
	testquicksort()
	selectsort(randNum())
	insertsort(randNum())

}

//生成随机切片
func randNum() []int {
	rand.Seed(time.Now().UnixNano())
	arr := make([]int, 50)
	for i := 0; i < 50; i++ {
		arr[i] = rand.Intn(101)
	}
	return arr
}

//冒泡排序
func bubblesort(s []int) {
	n := len(s)
	fmt.Println("冒泡排序前的数组:", s)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if s[j] > s[j+1] {
				s[j], s[j+1] = s[j+1], s[j]
			}

		}

	}
	fmt.Println("冒泡排序结果:", s)
	fmt.Println()
	fmt.Println()

}

//选择排序,从左往右
func selectsort(s []int) {
	fmt.Println("选择排序前的数:", s)
	n := len(s)

	for i := 0; i < n-1; i++ {
		min := s[i]   //最小数
		minIndex := i //最小数的下标
		for j := i + 1; j < n; j++ {
			if s[j] < min {
				min = s[j]
				minIndex = j

			}
		}
		//交换
		if i != minIndex {
			s[i], s[minIndex] = s[minIndex], s[i]
		}

	}
	fmt.Println("选择排序后的数:", s)
	fmt.Println()
	fmt.Println()
}

//插入排序,从第二个开始找,比他小的插入到前面,他后移
func insertsort(s []int) {
	fmt.Println("插入排序前的数:", s)

	n := len(s)
	for i := 1; i < n; i++ { //从第二位开始
		deal := s[i] //待排数
		j := i - 1   //待排数前一位元素
		for ; j >= 0 && deal < s[j]; j-- {
			s[j+1] = s[j] //后移,赋值给他后一位
		}
		s[j+1] = deal //出循环是j=-1
	}
	fmt.Println("插入排序后的数:", s)
	fmt.Println()
	fmt.Println()
}

//快速排序
func testquicksort() {
	s := randNum()
	fmt.Println("快速排序前的数:", s)
	fmt.Println("快速排序后的数:", quicksort(s))
	fmt.Println()
	fmt.Println()

}

func quicksort(s []int) []int {
	if len(s) <= 1 {
		return s
	}
	//选第一个数作为基准
	mid := s[0]
	smaller := make([]int, 0, 0)
	bigger := make([]int, 0, 0)
	equal := make([]int, 0, 0)
	equal = append(equal, mid)
	//根据大小放入切片
	for i := 1; i < len(s); i++ { //0已经是基准数了
		if s[i] < mid {
			smaller = append(smaller, s[i])
		} else if s[i] > mid {
			bigger = append(bigger, s[i])

		} else {
			equal = append(equal, s[i])
		}

	}
	//递归
	smaller, bigger = quicksort(smaller), quicksort(bigger)
	//拼接
	reslut := append(append(smaller, equal...), bigger...)
	return reslut
}
