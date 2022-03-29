package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

/* 用于读取配置文件,用的是ini包
https://ini.unknwon.io/docs/intro/getting_started
用于程序初始化


*/
//先定义从配置文件的选项的数据类型
var (
	Cfg *ini.File //句柄

	RunMode string

	HTTPPort int

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PageSize  int
	JwtSecret string
)

//初始函数,读取配置文件
func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatal("Fail to Parse 'conf/app.ini': ", err)
	}
	//分别加载各个分区,并获取值
	LoadBase()
	LoadApp()
	LoadServer()

}

func LoadBase() {
	//Muststring在获取不到值时会使用设置的默认值
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")

}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}
	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	//转换为时间
	ReadTimeout = time.Duration(Cfg.Section("server").Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = Cfg.Section("server").Key("WRITE_TIMEOUT").MustDuration(60) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}
