package main

import (
	"math/rand"
	"os"
	"runtime/trace"
	"time"
)
//https://zhuanlan.zhihu.com/p/332501892
func generate(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}
func bubbleSort(nums []int) {
	for i := 0; i < len(nums); i++ {
		for j := 1; j < len(nums)-i; j++ {
			if nums[j] < nums[j-1] {
				nums[j], nums[j-1] = nums[j-1], nums[j]
			}
		}
	}
}

func main() {
	//创建文件记录数据
	f, _ := os.OpenFile("trace.out", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()
	n := 10
	for i := 0; i < 5; i++ {
		nums := generate(n)
		bubbleSort(nums)
		n *= 10
	}
}
