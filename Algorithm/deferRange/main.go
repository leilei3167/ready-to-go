package main

import "fmt"

//深入理解defer和return的执行过程,具名返回值和匿名返回值有什么不同?
//for循环时,迭代的值始终都是同一块地址

func main() {
	fmt.Println(f1()) // 5
	fmt.Println(f2()) // 10
	f3()              //按序输出
	f4()              //全部输出5
	f5()              //输出9个10!!!!!!!,然后0-9
}

func f1() (r int) {
	t := 5
	defer func() {

		t = t + 5

	}()

	return t
}

func f2() (t int) {
	t = 5
	defer func() {

		t = t + 5

	}()

	return t
}

func f3() {
	c := []int{1, 2, 3, 4, 5}

	for _, v := range c {
		fmt.Println("f3:", v)
	}

}
func f4() {
	c := []int{1, 2, 3, 4, 5}

	for _, v := range c {
		defer func() {
			fmt.Println("f4:", v)
		}()
	}

}

func f5() {
	funcs := []func(){}

	for i := 0; i < 10; i++ {
		funcs = append(funcs, func() {
			fmt.Println(i)
		})
		//超级巨坑,出循环时i=10!
	}

	for i := 0; i < 10; i++ {
		f := func(index int) func() {
			return func() {
				fmt.Println(index)
			}
		}(i) //index=i,产生了拷贝,index不是同一个i
		funcs = append(funcs, f) //funcs中9个打印9的函数,9个打印0-9的函数
	}

	//全部执行
	for _, f := range funcs {
		f()
	}

}
