package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open failed") //在原始err上包装错误,warp将会附上堆栈信息(+v打印),而withmessage则不会有堆栈信息
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}
	return buf, nil
}

func ReadConfig() ([]byte, error) {
	home := os.Getenv("HOME")
	config, err := ReadFile(filepath.Join(home, ".settings.xml"))
	return config, errors.WithMessage(err, "could not read config") //如果err是nil 则返回nil
}

func main() {
	_, err := ReadConfig()
	if err != nil {
		fmt.Printf("直接打印错误信息:%v\n\n", err) //could not read config: open failed: open /home/lei/.settings.xml: no such file or directory
		fmt.Println("-------使用cause--------")
		fmt.Printf("使用Cause处理:,类型: %T err值:%v\n", errors.Cause(err), errors.Cause(err)) //cause打印的是最底层的失败原因,不带任何附加包装的信息
		fmt.Println("--------堆栈信息-------\n\n")
		//	fmt.Printf("stack trace:\n%+v\n", err) //+v可打印堆栈追踪
		os.Exit(1)
	}
}
