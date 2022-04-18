package main

import (
	"bytes"
	"fmt"
)

/* 1.切片 */
//切片不能比较,无法使用==判断两个slice是否具有相同的元素,
//但是对于byte切片来讲,有bytes.Equal函数来判断两个字节切片是否相等
//make创建slice,不指定容量时默认等于长度;slice扩容将发生底层数组间的拷贝和内存分配以及GC,尽量make时设置长度,避免频繁复制
func main() {
	//	a := []byte{'1', '2'}
	b := []byte{'1', '2', '1'}
	var a []byte
	fmt.Println(bytes.Equal(a, b))

	//----------
	data := []string{"one", "", "three"}
	//data = nonempty(data)
	fmt.Printf("%q\n", nonempty(data)) // `["one" "three"]`
	fmt.Printf("%q\n", data)           // `["one" "three" "three"]`

	//-----------------
	s := []int{5, 6, 7, 8, 9}
	fmt.Println(remove(s, 2)) // "[5 6 8 9]"

	//---------
	s1 := []int{5, 6, 7, 8, 9}
	s2 := s1[2:] //7.8.9
	s3 := s1[3:] //8.9
	copy(s2, s3)
	fmt.Println("s2", s2) //8.9.9
	//test copy:copy会按照索引来复制
	s4 := []int{5, 6, 7, 8, 9}
	s5 := []int{0, 0}
	copy(s4, s5)
	//---------不要对引用类型的底层数据取值(包括map),因为发生扩容后其地址会发生变更
	fmt.Println("s4", s4) //0.0.7.8.9
	fmt.Printf("扩容前的地址%v\n", &s4[1])
	s4 = append(s4, 1, 2, 3, 1, 4, 45, 5)
	fmt.Printf("扩容后的地址%v\n", &s4[1])

}

func nonempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	return strings[:i] //共享了一个底层数组,导致原数组可能被覆盖,因此用data=noneempyt(data)来将新值赋值
}

//要删除slice中间的某个元素并保存原有的元素顺序，可以通过内置的copy函数将后面的子slice向前依次移动一位完成
func remove(slice []int, i int) []int {
	//7,8,9copy8,9 ->8,9,9
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}
