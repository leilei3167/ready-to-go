package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

//1.一次性读取,适用于一般大小的文件
func read1() {
	bs, err := ioutil.ReadFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ioutil读取结果:", string(bs)) //打印字符切片的话全是字符码
}

//2.os包里也有readfile
func read2() {
	bs, err := os.ReadFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("os读取结果:", string(bs))

}

//3.也可以创建句柄实现后读取
func read3() {
	//也可使用openfile函数,可指定权限
	file, err := os.Open("test.txt") //file实现了io.reader接口
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("open创建句柄之后的结果:", string(bs))

}

//对于较大文件,建议使用有缓存的读取
func read4() {
	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	//创建bufio.Reader
	r := bufio.NewReader(file)
	//创建容器,每次读取多少字节
	buf := make([]byte, 1024)
	//循环读取
	for {
		//调用bufio.reader的方法
		n, err := r.Read(buf)
		if err != nil && err != io.EOF { //有错误
			log.Fatal(err)
		}
		if n == 0 { //每次循环读取到的字节为0时
			break
		}
	}
	fmt.Println(string(buf))

}

func main() {
	//read1()
	//read2()
	//read3()
	read4()
}
