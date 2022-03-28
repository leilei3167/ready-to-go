package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

}

/*一:将整个文件读取入内存
效率最高,但仅适用于小文件,对于大文件不合适,因为会使用过多的内存
*/

//1.1直接指定文件名读取
//第一种方法:os.ReadFile
func OsreadFile() {
	bs, err := os.ReadFile("a.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}

//第二种方法 ioutil.ReadFile
func ioutilReadFile() {
	bs, err := ioutil.ReadFile("a.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}

//1.2创建句柄再读取(先用Open(打开一个只能读取的文件)或OpenFile打开文件)
func openAndRead() {
	file, err := os.OpenFile("a.txt", os.O_RDONLY, 0) //可实现多功能的读写
	if err != nil {
		panic(err)
	}
	defer file.Close() //非常重要,open后一定要记得关闭
	bs, _ := ioutil.ReadAll(file)
	fmt.Println(string(bs))

}

/*二:每次只读取一行
一次性读取所有数据太耗费内存,因此可以指定每次指定只读取一行数据
也可以使用bufio.Readline,但官方不建议
*/
//2.1使用bufio.ReadBytes
func bufioreadbytes() {
	fi, err := os.Open("a.go")
	if err != nil {
		panic(err)
	}
	//创建reader
	r := bufio.NewReader(fi)
	for { //循环直到读取到EOF错误
		lineBytes, err := r.ReadBytes('\n')
		line := strings.TrimSpace(string(lineBytes)) //去掉单行的数据的空格
		if err != nil {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		fmt.Println(line)

	}

}

//2.2使用bufio.ReadString
func bufioreadstring() {
	fi, err := os.Open("a.go")
	if err != nil {
		panic(err)
	}
	//创建reader
	r := bufio.NewReader(fi)
	for { //循环直到读取到EOF错误
		lineBytes, err := r.ReadString('\n')
		line := strings.TrimSpace(string(lineBytes)) //去掉单行的数据的空格
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		fmt.Println(line)

	}
}

/*三:每次读取固定字节数
每次仅读取一行数据,可以解决内存占用过大的问题,但要注意的是,不是所有的文件都会有\n,
对于不换行的大文件怎么办??
*/
//3.1使用os库,创建句柄然后创建一个Reader,然后在for循环里调用Reader的Read函数
func readbyNum() {
	fi, err := os.Open("a.go")
	if err != nil {
		panic(err)
	}
	//创建reader
	r := bufio.NewReader(fi)

	//每次读取1024个字节
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf) //n是每次读取的字节数,为0时说明读取完毕
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

	}

}
