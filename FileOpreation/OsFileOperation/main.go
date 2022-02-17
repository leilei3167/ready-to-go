package main

import (
	"fmt"
	"os"
)

func main() {
	//createFile()创建一个a.txt
	//createDir() //创建文件夹a/b/c
	//remove() //删除a.txt或文件夹
	writeFile()
	readFile()
}

//创建文件
func createFile() {
	f, err := os.Create("a.txt") //会创建到项目根目录
	if err != nil {
		fmt.Println(err)

	} else {
		fmt.Println("f name is:", f.Name())

	}

}

//创建目录
func createDir() {
	/*	err := os.Mkdir("test", os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}*/
	//创建一系列的文件目录
	err2 := os.MkdirAll("a/b/c", os.ModePerm)
	if err2 != nil {
		fmt.Println(err2)
	}
}

//删除目录或者文件
func remove() {
	//删除单个文件
	/*err := os.Remove("a.txt")
	if err != nil {
		fmt.Println(err)
	}*/
	//删除某个目录
	err2 := os.RemoveAll("a")
	if err2 != nil {
		fmt.Println(err2)
	}
}

//直接读写文件,本质也是先Open
func readFile() {
	b, _ := os.ReadFile("FileOpreation/OsFileOperation/test.txt")
	fmt.Printf("b:%v\n", string(b[:]))
}

//本质是先OpenFile,OpenFile(name, O_WRONLY|O_CREATE|O_TRUNC, perm)
func writeFile() {
	os.WriteFile("FileOpreation/OsFileOperation/test.txt", []byte("nimade "), os.ModePerm)

}

//执行了Open,务必要记得关闭,即defer file.close
