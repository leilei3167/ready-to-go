package main

import (
	"fmt"
	"log"

	"github.com/go-ini/ini"
)

func main() {
	//读取配置文件
	cfg, err := ini.Load("go-ini/my.ini")
	if err != nil {
		log.Fatal("读取配置文件失败!", err)
	}

	//典型读取,默认分区可使用空字符串
	fmt.Println("app mode:", cfg.Section("").Key("app_mode").String())
	//其余字段需加分区,表示从paths分区找data选项对应的值,输出为字符串
	fmt.Println("data:", cfg.Section("paths").Key("data").String())

	//可以限制候选值的操作,如果不在候选区内,则返回设定的默认值
	fmt.Println("server protocol:",
		cfg.Section("server").Key("protocol").In("http", []string{"tcp", "https"}))
	//会优先选择选项内的
	fmt.Println("http port:",
		cfg.Section("server").Key("http_port").In("8080", []string{"9090", "9999"}))

	//自动类型转换
	fmt.Printf("Port Number:(%[1]T)(%[1]d)\n", cfg.Section("server").Key("http_port").MustInt())
	fmt.Printf("Enforce Domain: (%[1]T) %[1]v\n", cfg.Section("server").Key("enforce_domain").MustBool(false))

	//修改配置文件
	cfg.Section("").Key("app_mode").SetValue("production")
	cfg.SaveTo("go-ini/my.ini")

}
