package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	//生成一个随机数切片
	s := Sort()
	fmt.Println("排序前的数组:", s)
	bubblesort(s)
	selectsort(s)
	insertsort(s)
	m := quicksort(s)
	fmt.Println("快速排序后的数组:", m)
}

//生成随机数排序前

func Sort() []int {
	s := make([]int, 50)
	seedNum := time.Now().UnixNano()
	rand.Seed(seedNum)
	for i := 0; i < 50; i++ {

		s[i] = rand.Intn(10001)

	}

	return s

}

//冒泡排序
func bubblesort(s []int) {
	//得到长度
	n := len(s)

	for i := 0; i < n-1; i++ { //需要进行n-1轮
		for j := 0; j < n-1-i; j++ { //每轮需要的比较次数,每一轮都会确定一个,因此要-i
			if s[j] > s[j+1] {
				s[j], s[j+1] = s[j+1], s[j]
			}

		}

	}

	fmt.Println("冒泡排序后的数组:", s)
}

//选择排序,从左往右扫描,将最小的放到最左,第二次从第二个开始扫描,一次类推,效率同样糟糕,比较次数相同
func selectsort(s []int) {

	n := len(s)
	for i := 0; i < n-1; i++ {
		//每次从第i位开始,找最小元素
		min := s[i]   //最小数
		minIndex := i //最小数的下标
		//假设最左边是最小数,从他下一位开始找,如果找到更小的,将它设为最小值
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

	fmt.Println("选择排序后的数组:", s)

}

//插入排序,从第二位开始比较,数列越有序越快,从第二个数开始,向后找更小的数
func insertsort(s []int) {
	n := len(s)
	for i := 1; i <= n-1; i++ { //第二位开始比
		deal := s[i]                       //待排序的数
		j := i - 1                         //待排数前一个数
		for ; j >= 0 && deal < s[j]; j-- { //判断待排数是否比前一位更小
			s[j+1] = s[j] //后移,赋值给他后一位

		}
		s[j+1] = deal

	}
	fmt.Println("插入排序后的数组:", s)

}

//快速排序,假设一位数是中位数,所有比他小的放在左边,比他大的放在右边,分割成两部分,整个排序递归进行
func quicksort(s []int) []int {
	if len(s) <= 1 { //递归判断
		return s
	}
	splitdata := s[0]         //基准数,取第一个
	low := make([]int, 0, 0)  //比基准数小的
	high := make([]int, 0, 0) //比基准数大的
	mid := make([]int, 0, 0)  //和基准数一样大的数
	mid = append(mid, splitdata)

	for i := 1; i < len(s); i++ {
		if s[i] < splitdata { //如果小于基准数,加入low切片
			low = append(low, s[i])

		} else if s[i] > splitdata { //如果大于基准数,加入high切片
			high = append(high, s[i])

		} else {
			mid = append(mid, s[i])
		}

	}
	//递归执行
	low, high = quicksort(low), quicksort(high)
	myarr := append(append(low, mid...), high...)
	return myarr

}
