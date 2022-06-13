package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	pwd, _ := os.Getwd() // 获取到当前目录，相当于python里的os.getcwd()
	fmt.Println("当前的操作路径为:", pwd)
	//文件路径拼接
	f1 := filepath.Join(pwd, "test", "1.txt")
	fmt.Println("文件的路径为:", f1)
	//文件的目录名
	fmt.Println("文件的目录名:", filepath.Dir(f1))
	//文件的文件名
	fmt.Println("文件的文件名:", filepath.Base(f1))
	//文件的绝对路径
	adspath, _ := filepath.Abs("evn/3.txt")
	fmt.Println("文件的绝对路径为:", adspath)
	//拆分路径
	dirname, filename := filepath.Split(f1)
	fmt.Println("目录名为:", dirname, "文件名为", filename)
	//扩展名相关
	fmt.Println("f1的扩展名为:", filepath.Ext(f1))

	//通过os.Stat()函数返回的文件状态，如果有错误则根据错误状态来判断文件或者文件夹是否存在
	fileinfo, err := os.Stat(f1)
	if err != nil {
		fmt.Println(err.Error())
		if os.IsNotExist(err) {
			fmt.Println("file:", f1, " not exist！")
		}
	} else {
		//判断路径是否为文件夹
		fmt.Println(fileinfo.IsDir())
		fmt.Println(!fileinfo.IsDir())
		fmt.Println(fileinfo.Name())
	}
	fmt.Println("当前工作目录:", getCurrentDir())
}

//当前工作目录
func getCurrentDir() string {
	_, fileName, _, _ := runtime.Caller(1)
	path, _ := filepath.Split(fileName)
	/* 	aPath := strings.Split(fileName, "/")
	   	dir := strings.Join(aPath[0:len(aPath)-1], "/") */
	return path[:len(path)-1]
}
