package main

import "fmt"

func main() {
	fmt.Println(fun1([]int{1, 1, 1, 2, 3, 3, 4}))
	fmt.Println(fun2([]int{1, 1, 1, 2, 2, 2, 2, 3, 3, 4}))
	fmt.Println(func3([]int{1, 2, 3, 4, 5, 6, 9}, 8))
}

//103.给定一个有序数组，去掉重复元素
func fun1(s []int) []int {
	//利用map key不能重复的原理
	m := make(map[int]struct{})
	res := make([]int, 0)
	for _, v := range s {
		if _, ok := m[v]; !ok {
			//说明m中没有v
			res = append(res, v)
			m[v] = struct{}{}
		}
	}
	return res

}

//102.给一个有序数组，计算出现次数最多的元素

func fun2(s []int) int {
	m := make(map[int]int)

	for i := 0; i < len(s); i++ {
		m[s[i]]++
	}
	//获取value最大的对应的key
	maxv := 0
	maxk := 0
	for k, v := range m {
		if v >= maxv {
			maxv = v
			maxk = k
		}
	}
	return maxk
}

//104.给定一个严格递增的有序数组，找到两个数字其和等于给定数字，输出有多少种情况
func func3(s []int, val int) int {
	temp := 0
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i]+s[j] == val {
				temp++

			}
		}
	}
	return temp
}
