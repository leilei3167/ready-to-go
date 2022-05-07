package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unsafe"
)


func main() {
	//ReadAt()
	//WriteAt()
	//ReadFrom()

	a := make([]int, 3, 3)
	a = append(a, 1000)[:len(a)] //任意加入一个元素触发扩容,舍弃元素
	fmt.Printf("%#v,len:%v,cap:%v\n", a, len(a), cap(a))

	ip := "127.0.0.1"
	fmt.Println(unsafe.Sizeof(ip) * 655360)
	iprange := strings.Split(ip, ".")[0]
	fmt.Printf("ip:%s\n", iprange)
	//LimitRead()
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

//标准库的ReadAll
type Reader interface {
	Read(p []byte) (n int, err error)
}

func ReadAll(r Reader) ([]byte, error) {
	b := make([]byte, 0, 512) //创建缓冲区,512字节
	for {
		if len(b) == cap(b) { //长度和容量相等时,主动触发扩容
			b = append(b, 0)[:len(b)]
		}
		//Read是读取len(p)个字节!!,len(p):cap(p)刚好就是空闲的空间
		n, err := r.Read(b[len(b):cap(b)]) //截取空闲的缓冲空间,并Read数据进去
		b = b[:len(b)+n]                   //去除多余的部分
		if err != nil {
			if err == io.EOF { //不把EOF视作异常
				err = nil
			}
		}
		return b, err
	}
}

func LimitRead() {
	a := "hello world!!"
	lreader := io.LimitReader(strings.NewReader(a), 3)
	b := make([]byte, 512)
	lreader.Read(b)
	fmt.Printf("%v\n", string(b))

}
