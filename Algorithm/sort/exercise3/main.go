package main

import "fmt"

func main() {
	a := []int{1, 3, 5, 2, 4}
	bubblesort(a)
	b := []int{3, 2, 4, 5}
	selectsort(b)
	c := []int{5, 4, 3, 2, 1}
	insertsort(c)
	d := []int{6, 5, 4, 2, 4, 9, 1}
	res := quicksort(d)
	fmt.Println("快排前:", d, "快排后:", res)
}

//冒泡排序
func bubblesort(s []int) {
	n := len(s)
	fmt.Println("冒泡排序前:", s)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if s[j] > s[j+1] {
				s[j], s[j+1] = s[j+1], s[j]

			}

		}

	}
	fmt.Println("冒泡排序后:", s)
}

//选择排序,假定每次循环最左边为最小值,后一位开始遍历剩余数组每一个找比他更小的,找完之后若有则交换
func selectsort(s []int) {
	n := len(s)
	fmt.Println("选择排序前:", s)
	for i := 0; i < n-1; i++ {
		minNum := s[i]
		minIndex := i
		for j := i + 1; j < n; j++ {
			if s[j] < minNum {
				minNum = s[j]
				minIndex = j
			}
		}
		if minIndex != i {
			s[i], s[minIndex] = s[minIndex], s[i]
		}
	}
	fmt.Println("选择排序后:", s)
}

//插入排序,假定一开始第二位是待排数
func insertsort(s []int) {
	n := len(s)
	fmt.Println("插入排序前:", s)
	//从第二位开始
	for i := 1; i < n; i++ {
		//待排数
		deal := s[i]
		j := i - 1                         //待排数前一位
		for ; j >= 0 && s[j] > deal; j-- { //如果前一位比他还大的话,执行后移
			s[j+1] = s[j]
		}
		//出循环是j=-1
		s[j+1] = deal
	}
	fmt.Println("插入排序后:", s)
}

//快速排序
func quicksort(s []int) []int {
	if len(s) <= 1 {
		return s
	}
	//选取第一个作为基准数
	mid := s[0]
	small := make([]int, 0, 0)
	big := make([]int, 0, 0)
	equal := make([]int, 0, 0)
	equal = append(equal, mid)
	for i := 1; i < len(s); i++ {
		if s[i] < mid {
			small = append(small, s[i])
		} else if s[i] == mid {
			equal = append(equal, s[i])

		} else {
			big = append(big, s[i])

		}

	}
	//递归
	small, big = quicksort(small), quicksort(big)
	s = append(append(small, equal...), big...)
	return s

}
