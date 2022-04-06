package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	Ago()
	fmt.Println(SplitStrToMap("type:process|ip:127.0.0.1", "|", ":"))
}

//获取当前执行文件的绝对路径
func GetPath() string {
	//Abs函数返回path代表的绝对路径，如果path不是绝对路径，会加入当前工作目录以使之成为绝对路径
	//Dir返回路径除去最后一个路径元素的部分，即该路径最后一个元素所在的目录
	//Dir参数此处获取执行文件的名称
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	dir1 := strings.Replace(dir, "\\", "/", -1)
	fmt.Println(dir)
	fmt.Println(dir1)
	return dir1
}

//处理字符串,将传入的字符串按指定的符号分割成map
//type:process|ip:127.0.0.1
func SplitStrToMap(str string, sep1 string, sep2 string) map[string]string {
	sublist := strings.Split(str, sep1)
	var resultmap = map[string]string{}
	for _, pair := range sublist {
		z := strings.SplitN(pair, sep2, 2)
		if len(z) == 2 {
			key := strings.TrimSpace(z[0])
			value := strings.TrimSpace(z[1])
			key = strings.ToLower(key)
			resultmap[key] = value
		}
	}
	return resultmap
}

//获取从现在开始 一定时间之前的时间
func Ago() {
	OneMin := time.Now().Add(time.Minute * time.Duration(-1))
	fmt.Println(OneMin)
	TenMin := time.Now().Add(time.Minute * 10 * time.Duration(-1))
	fmt.Println(TenMin)
}
