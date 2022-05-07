package main

import (
	"log"
	"os"
)

func main() {
	file, err := os.Open("../small.txt")
	if err != nil {
		log.Fatal(err)
	}
	info, _ := file.Stat()
	log.Printf("%+v\n", info)

}
/* type FileInfo interface {
    Name() string       // 文件的名字
    Size() int64        // 普通文件返回值表示其大小；其他文件的返回值含义各系统不同
    Mode() FileMode     // 文件的模式位
    ModTime() time.Time // 文件的修改时间
    IsDir() bool        // 等价于Mode().IsDir()
    Sys() interface{}   // 底层数据来源（可以返回nil）
} */