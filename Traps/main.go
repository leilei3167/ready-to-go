package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	//rangeTrap()
	//	rangeTrap2()
	rangeTrap3()
	//	rangeTrap4()
	//commaMap()
	//commaChan()
	commaInterface()
	fmt.Println(os.Getwd())
}

/*for range陷阱:
在一个range循环中,其中的k始终的地址都是不变的(可以循环打印k的地址验证),如果在迭代中使用k,将会使其每一次都被新的值覆盖
*/
//for range陷阱1:只会输出四个4
func rangeTrap() {
	wg := sync.WaitGroup{}
	si := []int{1, 2, 3, 5, 6}
	for k, v := range si { //i为0到4
		wg.Add(1)
		//fmt.Println("out:", k, v)
		//time.Sleep(time.Second) //此处是由于range瞬间结束,最终值被更新为4,所以协程全部打印4(i的地址不变,上面的值会更新为最新值)
		//添加等待后就能发现能够按照顺序输出123
		go func() {
			println(k, v)
			wg.Done()
		}()
	}
	wg.Wait()

}

//输出123123 而不是成为永动机,因为每一个range中都会将原切片拷贝,并获取到了原切片的长度,所以在循环中增加新元素也不会改变循环次数
func rangeTrap2() {
	si := []int{1, 2, 3}
	for _, v := range si { //因为被range 的si此刻是一个副本,其长度是固定的
		si = append(si, v)
	}
	fmt.Println(si)
}

//会输出333,对于需要取得元素的range,会额外创建一个新的变量用于存储获得的元素,在循环中使用这个变量将会是的其在每一次迭代被新值覆盖
func rangeTrap3() {
	arr := []int{1, 2, 3}
	newArr := []*int{}
	//v的值就像是:1,22,333即每次被更新,之前添加的v值会随着改变
	for _, v := range arr {
		newArr = append(newArr, &v) //正确的方法应该是按照i来取值,即&arr[i]
	}
	for k, v := range newArr {
		fmt.Println(k, *v)
	}
}
func rangeTrap4() {
	funcs := []func(){}       //以匿名函数为values的切片
	for i := 0; i < 10; i++ { //添加10个打印临时变量i的匿名函数进去
		funcs = append(funcs, func() {
			fmt.Println(i)
		})
	} //重点陷阱:出循环时i=10,之前添加的i全部为10
	for i := 0; i < 10; i++ {
		f := func(index int) func() { //构成闭包,引用了外部变量i,i被分配到堆上
			return func() {
				fmt.Println(index)
			}

		}(i)
		funcs = append(funcs, f) //添加10个闭包,与闭包绑定的值为0-9(虽然出循环
	}
	for _, v := range funcs {
		v() //每一个储存的匿名函数执行

	}
	//会输出10个10,0-9
}

/*comma,ok表达式*/
func commaMap() {
	//1.用于判断map中某个键是否存在值,因为查找一个不存在的key,只会返回零值
	m := make(map[string]string)
	v, ok := m["nigger"]
	if !ok {
		println("不存在这个值")
	} else {
		println(v)
	}

}
func commaChan() {
	//2.读取一个已关闭的chan时不会阻塞,而是返回零值,此表达式可判断chan是否关闭
	c := make(chan int)
	go func() {
		c <- 1
		time.Sleep(time.Second)
		c <- 2
		time.Sleep(time.Second)
		c <- 3
		time.Sleep(time.Second)

		close(c)
	}()
	/*	for {
		v, ok := <-c
		if ok {
			println(v)
		} else {//一旦关闭,就会返回false
			println("c已关闭,退出")
			break
		}
	}*/
	//直接使用range也可以,chan关闭就会立刻感知到
	for v := range c {
		println(v)
	}

}

//3.类型断言或类型选择,用于判断接口所绑定的实例,是否同时实现了另一个接口或是该实例是否是某个类型
func commaInterface() {
	var x interface{}
	x = 1
	v, ok := x.(int32) //x是否是string类型,若是,其值会被复制给v,并且ok为true
	if ok {
		fmt.Println("断言成功,他就是断言的类型", v)
	} else {
		fmt.Printf("断言失败,v的值%#v\n\n", v) //若断言失败,v为断言类型的零值,ok为false

	}
	//如果要断言多个类型的话就可以用 switch v:=x.(type)
}

//4.init函数和全局变量陷阱

var cmd string

//虽然cwd在外部已经声明过，但是:=语句还是将cwd和err重新声明为新的局部变量。因为内部声明的cwd将屏蔽外部的声明，因此上面的代码并不会正确更新包级声明的cwd变量。
func init() {
	cmd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	//var err error
	//cmd, err = os.Getwd() //正确打印
	log.Println(cmd)
}
