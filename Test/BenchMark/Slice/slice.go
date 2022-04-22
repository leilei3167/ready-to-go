package slice

//删除切片的某个元素

func Delete(s []string, i int) []string {
	if i < len(s)-1 {
		copy(s[i:], s[i+1:])
	}
	s[len(s)-1] = "" //最后一位置零 便于GC回收
	s = s[:len(s)-1] //只保留
	return s

}

//性能陷阱,直接在原切片上截取和copy的区别(截取可能导致底层数组一直不会释放)

//截取最后2个元素,因为返回的切片和原切片共用一个底层数组(1M),函数结束后也不会释放
func lastNumsBySlice(origin []int) []int {
	return origin[len(origin)-2:]
}

//copy,返回一个全新的切片,原底层数组将不被引用,会被回收
func lastNumsByCopy(origin []int) []int {
	result := make([]int, 2)
	copy(result, origin[len(origin)-2:])
	return result
}
