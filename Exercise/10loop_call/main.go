package main

import "fmt"

func main() {
	u := &user{}
	u.String() //栈溢出
}

type user struct {
	name string
}

func (u *user) String() string {
	return fmt.Sprintf("print:%v", u) //造成循环调用,因为Print系列函数会优先调用String()方法!

}
