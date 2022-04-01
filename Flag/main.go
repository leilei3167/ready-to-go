package main

import (
	"flag"
	"fmt"
)

//编写命令行程序（工具、server）时，我们有时需要对命令参数进行解析

/* 命令行参数（或参数）：是指运行程序时提供的参数；
已定义命令行参数：是指程序中通过 flag.Type 这种形式定义了的参数；
非 flag（non-flag）命令行参数（或保留的命令行参数）：可以简单理解为 flag 包不能解析的参数。 */

var (
	name = flag.String("name", "root", "输入姓名")
)

func main() {
	flag.Parse()
	//主要2种定义flag的方法 分别是直接定义或者放入到某个变量
	//1.flag.Type(flag的名字,默认值,帮助信息)*Type
	//如 name=flag.String("name","张三","姓名")
	//如 age=flag.Int("age","18","年龄")
	//以上age 和 name均为指针

	//2.flag.TypeVar(Type指针,flag名,默认值,帮助信息)
	//var name string
	//flag.StringVar(&name,"name","leilei","名字")

	//通过以上2种方式定义好flag参数后,必须调用flag.Parse来进行解析
	fmt.Println(*name)
	//输入-name=leilei 则会打印leilei
	//输入 -h则会显示帮助信息
}
