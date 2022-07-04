package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	Greet(os.Stdout, "Elodie")
}

//依赖某个接口,可以使得我们按需来注入依赖,使得我们可以控制输入位置
func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}

//此版本的无法控制数据的写入目的地,非常不利于测试,应该使用更广泛的接口
func Greet1(name string) {
	fmt.Printf("Hello, %s", name)
}

/*
依赖注入的作用:

1.一个难以编写单元测试的函数往往是因为依赖被硬连接到某个函数或者全局的状态上,如你的函数依赖全局的数据库连接,
那么其将会非常难以测试;正确的做法应该是依赖接口,这样就使得我们能够进行依赖注入,从而脱离真实的数据库连接进行测试
2.分离关注点,将如何生成数据和数据从何而来解耦;如果你觉得某个函数具有太多的职责(比如生产存入db的数据,又要处理http请求),
那么关注依赖注入
3.使得我们的代码在不同的环境下更具复用性


*/
