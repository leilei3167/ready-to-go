package main

import "fmt"

/* Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.

You may assume that each input would have exactly one solution, and you may not use the same element twice.

You can return the answer in any order.

Example 1:

Input: nums = [2,7,11,15], target = 9 Output: [0,1] Explanation: Because nums[0] + nums[1] == 9, we return [0, 1]. Example 2:

Input: nums = [3,2,4], target = 6 Output: [1,2] Example 3:

Input: nums = [3,3], target = 6 Output: [0,1] */
func main() {
	fmt.Println(leilei([]int{2, 7, 11, 15}, 26))

}

/*输入一个数字数组,以及一个目标值,输出相加能够得到目标值对应的下标*/
func leilei(input []int, target int) (res string) {
	n := len(input) //获取数组长度
	if n < 2 {
		fmt.Println("请输入至少2个值")
	}
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if input[i]+input[j] == target {
				res = fmt.Sprintf("[%v,%v]\n", i, j)
				return
			}

		}

	}
	return "没有结果!"
}
