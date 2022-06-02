package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var homePath = os.Getenv("PWD")
var savePath = "  .temp  /"

func main() {
	fmt.Println(Abs(savePath))
	testRel()
}

func Abs(s string) string {
	s = strings.TrimSpace(s) //去除空格
	if runtime.GOOS == "windows" {
	homePath=
	}
	path := filepath.Join(homePath, s) //将多个元素以路径分隔的形式连接
	if filepath.IsAbs(path) {
		return path
	}
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func testRel() {
	//home/lei
	fmt.Println(filepath.Rel(os.Getenv("HOME"), os.Getenv("PWD"))) ///home/lei/code/ready-to-go
}

