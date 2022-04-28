package parse

import (
	"flag"
	"os"
)

type Config struct {
	Name  string
	Age   string
	Class string
	Major string
}

func Banner() {
	banner := `
   ___                              _    
  / _ \     ___  ___ _ __ __ _  ___| | __ 
 / /_\/____/ __|/ __| '__/ _` + "`" + ` |/ __| |/ /
/ /_\\_____\__ \ (__| | | (_| | (__|   <    
\____/     |___/\___|_|  \__,_|\___|_|\_\   
                      version: ` + "test" + `
`
	print(banner)

}

func Parse(c *Config) {
	Banner()
	//var是直接将参数绑定到字段或者变量中
	flag.StringVar(&c.Age, "a", "", "年龄")
	flag.StringVar(&c.Name, "n", "admin", "名字")
	flag.StringVar(&c.Class, "c", "", "班级")
	flag.StringVar(&c.Major, "m", "", "专业")
	flag.Parse() //必须要解析
	//解析之后，flag的值可以直接使用。如果你使用的是flag自身，它们是指针；如果你绑定到了某个变量，它们是值,此种情况是值
	validate(c)

}

//解析输入是否符合要求
func validate(c *Config) {
	if c.Age == "" || c.Class == "" || c.Major == "" || c.Name == "" {
		flag.Usage()
		os.Exit(1)
	}

}
