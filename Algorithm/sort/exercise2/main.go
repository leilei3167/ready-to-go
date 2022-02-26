package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	bubblesort(randNum())
	selectsort(randNum())
	insertsort(randNum())
	testquicksort()
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
	fmt.Println("冒泡排序前:", s)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if s[j] > s[j+1] {
				s[j], s[j+1] = s[j+1], s[j]

			}

		}

	}
	fmt.Println("冒泡排序后:", s)

}

//选择排序
func selectsort(s []int) {
	n := len(s)

	fmt.Println("选择排序前:", s)
	for i := 0; i < n-1; i++ {
		//每轮开始时,最小值下标就是未排序数列的第一个值
		minNum := s[i]
		minIndex := i
		//找出每一轮的最小值
		for j := i + 1; j < n; j++ {
			if s[j] < minNum {
				minNum = s[j]
				minIndex = j

			}

		}
		//判断是否交换
		if minIndex != i {
			s[i], s[minIndex] = s[minIndex], s[i]
		}

	}
	fmt.Println("选择排序后:", s)
}

//插入排序
func insertsort(s []int) {
	n := len(s)
	fmt.Println("插入排序前:", s)
	for i := 1; i < n; i++ {
		//待排数
		deal := s[i]
		j := i - 1                         //待排数前一个下标
		for ; j >= 0 && s[j] > deal; j-- { //如果说前面的比待排数大,则需要后移
			s[j+1] = s[j]

		}
		//出循环时j=-1
		s[j+1] = deal

	}
	fmt.Println("插入排序后的数:", s)

}

//快速排序
func testquicksort() {
	s := randNum()
	fmt.Println("快速排序前的数:", s)
	fmt.Println("快速排序后的数:", quicksort(s))

}

func quicksort(s []int) []int {
	if len(s) <= 1 {
		return s
	}
	//选第一个为基准数
	mid := s[0]
	bigger := make([]int, 0, 0)
	smaller := make([]int, 0, 0)
	equ := make([]int, 0, 0)
	equ = append(equ, mid)
	for i := 1; i < len(s); i++ { //根据和基准数的对比放入不同切片
		if s[i] < mid {
			smaller = append(smaller, s[i])
		} else if s[i] > mid {
			bigger = append(bigger, s[i])

		} else {
			equ = append(equ, s[i])
		}

	}
	//递归
	smaller, bigger = quicksort(smaller), quicksort(bigger)
	res := append(append(smaller, equ...), bigger...)
	return res

}
