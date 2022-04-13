package main

import (
	"io"
	"log"
	"os"
)

//标准的log库
//默认输出到标准错误,自动加上日期和时间,会自动换行
type User struct {
	Age  int
	Name string
}

func main() {
	WithPrefixLog()
	WithFlagLog()
	OwnLogger()
}

//可设置日志前面的前缀
func WithPrefixLog() {
	u := &User{
		Name: "leilei",
		Age:  124,
	}
	log.SetPrefix("Login12321321312:")
	log.Println(u.Name, u.Age)

}

//log库提供了6个选项,可通过log.SetFlag设置
/* const (
  Ldate         = 1 << iota //当地的时间
  Ltime                    //当地时区的时间
  Lmicroseconds            //时间精确到微秒
  Llongfile                //输出长文件名+行号
  Lshortfile               //短文件名+行号,不含包名
  LUTC                     //UTC时间
) */

func WithFlagLog() {
	u := &User{
		Name: "nihao",
		Age:  12,
	}
	log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
	log.Println(u.Age, u.Name)

}

//自定义log:标准库事先定义了一个默认的logger,叫std(标准日志),直接调用
//log库的方法其实都是使用的std的方法 var std = New(os.Stderr, "", LstdFlags)

func OwnLogger() {
	u := &User{
		Name: "Jack",
		Age:  22,
	}
	file, err := os.Create("./Log/log.txt") //在当前目录下创建log文件
	if err != nil {
		panic("err open file")
	}
	//可以设置多出输出位置(也可以发送到网络)
	logger := log.New(io.MultiWriter(file, os.Stderr), "", log.Lshortfile|log.LstdFlags)
	logger.Println(u.Age, u.Name)

}

//标准库缺少日志级别 Info/debug/error等
//缺少日志结构化输出的功能
//性能一般

