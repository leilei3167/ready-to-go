package main

import "fmt"

func main() {
	fmt.Println(typeswitch(user{Name: "kedsa"}))
}

func typeswitch(i interface{}) string {

	switch i := i.(type) {
	case int, uint:
		return fmt.Sprintf("i is int:%#v 类型: %T", i, i)
	case bool, float64:
		return fmt.Sprintf("i is bool and float64:%#v 类型: %T", i, i)
	case string:
		return fmt.Sprintf("i is string %#v", i)
		//在多个类型列表的情况下,此时的i的类型是和其操作数一致(这里是interface{})
	case user, float32:
		return fmt.Sprintf("i is user %#v", i.Name)
	default:
		panic("not implemented")
	}

}

type user struct {
	Name string
}
