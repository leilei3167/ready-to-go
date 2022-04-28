package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"ready-to-go/Flag/test/parse"
	"strings"
)

func main() {
	var test parse.Config
	parse.Parse(&test) //获取输入
	log.Printf("解析参数后:%#v", test)
	//获取当前的工作路径
	log.Println("简化后的路径:", getPath())

}

func getPath() string {
	//1.搜索可执行文件的路径,如file中有斜杠，则只在当前目录搜索;此处会直接返回./main
	file, _ := exec.LookPath(os.Args[0]) //Args[0]代表程序的名字
	log.Println("lookpath查询程序名获得的路径为:", file)
	pwd, _ := os.Getwd() // 获取到当前目录，相当于python里的os.getcwd()
	log.Println("Getwd获得的路径为:", pwd)

	//2.返回./main的绝对路径,会包含test
	path1, _ := filepath.Abs(file)
	log.Println("绝对路径:", path1)
	//3.Dir返回路径除去最后一个路径元素的部分
	filename := filepath.Dir(path1) //不如直接使用os.Getwd
	log.Println("Dir之后:", filename)
	var path string
	if strings.Contains(filename, "/") {
		tmp := strings.Split(filename, `/`)
		tmp[len(tmp)-1] = ``
		path = strings.Join(tmp, `/`)
	} else if strings.Contains(filename, `\`) { //如果是windows
		tmp := strings.Split(filename, `\`)
		tmp[len(tmp)-1] = ``
		path = strings.Join(tmp, `\`)
	}

	return path

}
