package main

import (
	"log"
	"os"
)

var cmd string

//虽然cwd在外部已经声明过，但是:=语句还是将cwd和err重新声明为新的局部变量。因为内部声明的cwd将屏蔽外部的声明，因此上面的代码并不会正确更新包级声明的cwd变量。
func init() {
	cmd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	//var err error
	//cmd, err = os.Getwd() //正确打印
	log.Println(cmd)
}

func main() {
	//因为init中使用的是:=,函数会屏蔽外部的声明
	log.Printf("in main:%q", cmd)
}
