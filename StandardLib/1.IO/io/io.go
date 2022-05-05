package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
1.基本接口
在 io 包中最重要的是两个接口：Reader 和 Writer 接口。本章所提到的各种 IO 包，都跟这两个接口有关，也就是说，只要满足这两个接口，它就可以使用 IO 包的功能。

type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

*/

//os.File同时实现了以上2个接口!
//os.Stdin和os.Stdout是特殊的文件类型
/*
var (
    Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
    Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
    Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
)
其都是File类型,所以也都实现了RD两个接口

*/

/*
2. ReaderAt和WriterAt

type ReaderAt interface {
    ReadAt(p []byte, off int64) (n int, err error)
}

ReaderAt从指定的偏移量来读取数据,WriterAt在指定的偏移量位置写入
*/

/*
3. ReaderFrom 和 WriterTo接口

type ReaderFrom interface {
    ReadFrom(r Reader) (n int64, err error)
}

注意:ReadFrom不会返回err=EOF错误


*/

func main() {
	//ReadAt()
	//WriteAt()
	ReadFrom()
}

func ReadAt() {
	reader := strings.NewReader("hello world!")
	p := make([]byte, 16)
	n, err := reader.ReadAt(p, 2) //llo world!
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%s, %d\n", p, n)

}
func WriteAt() {
	file, err := os.Create("writeAt.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString("hello go-后续多余")

	n, err := file.WriteAt([]byte("hahaha"), 24) //在文件的offset=24处写入hahaha
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
}
func ReadFrom() {
	file, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	wirter := bufio.NewWriter(os.Stdout) //默认的size是4096,NewWriterSize可自定义大小
	wirter.ReadFrom(file)
	wirter.Flush()

}
