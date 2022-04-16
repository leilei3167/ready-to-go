package main

import "fmt"

//map在零值状态无法使用,必须make,或者直接用字面量构建

func main() {
	//遍历map是随机的无序的

	//如果用map做缓存,一定要记得在某个时机删除key

	//在map中尝试获取一个不存在的key时将返回零值

	a := map[string]bool{
		"leilei": true,
		"aha":    true,
	}

	v, ok := a["leilei"]
	if !ok {
		fmt.Printf("没有这个key:%v\n", ok)
	} else {
		fmt.Printf("成功:%v\n", v)
	}
	//对map进行排序,用sort包对key排序,之后遍历key切片取值
}
