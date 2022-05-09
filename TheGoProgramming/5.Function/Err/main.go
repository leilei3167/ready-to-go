package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

//错误处理,几种常见的方式
/*
1.向上传递错误
resp, err := http.Get(url)
if err != nil{
    return nil, err
}
传递错误应该添加清晰的错误链(可使用errors包的warp,或者标准库中新添加的%w)
if err != nil {
    return nil, fmt.Errorf("parsing %s as HTML: %v", url,err)
}
呈现的错误信息应该类似:genesis: crashed: no parachute: G-switch failed: bad relay orientation
一般而言，被调用函数f(x)会将调用信息和参数信息作为发生错误时的上下文放在错误信息中并返回给调用者，调用者需要添加一些错误信息中不包含的信息，比如添加url到html.Parse返回的错误中。

2.重试;如果错误的发生是偶然性的，或由不可预知的问题导致的。一个明智的选择是重新尝试失败的操作。在重试时，我们需要限制重试的时间间隔或重试的次数，防止无限制的重试

3.输出错误并终止程序;此种策略只应该在main中执行,对于库函数,应该只向上传递错误;
  log.Fatalf("Site is down: %v\n", err)

4.打印错误信息,程序继续执行;
if err := Ping(); err != nil {
    log.Printf("ping failed: %v; networking disabled",err)
}

5.直接用_忽略错误
*/

// LoginError 1.自定义错误传递,再os包中非常常见
type LoginError struct { //实现了error接口可作为error返回
	Name string
	Time string
	Err  error
}

func (l *LoginError) Error() string { //相当于包裹底层错误(预定义错误,添加更多信息)
	return l.Name + ":" + l.Time + ":" + l.Err.Error()
}
func (l *LoginError) Unwarp() error {
	return l.Err
}

var ErrLogin = errors.New("Login Time Out!")

func f1() error {
	//do smoething  ...if err!=nil ...
	return &LoginError{
		Name: "login",
		Time: fmt.Sprintf("%v", time.Now()),
		Err:  ErrLogin,
	}
}

//2.出现错误重试
func WaitForServer(url string) error {
	const timeout = time.Second * 6
	deadline := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(deadline); tries++ { //使用Before来确保每次循环在deadline之前
		client := http.Client{

			Timeout: time.Second * 2,
		}
		_, err := client.Head(url)
		if err == nil {
			return nil
		}
		//否则重试
		log.Printf("server not Responding(%s);retrying...", err)
		sleep := time.Second << uint(tries) //位运算,每次重试间隔翻倍
		log.Printf("going to sleep: %v", sleep)
		time.Sleep(sleep)
	}
	return fmt.Errorf("server %s failed to respond after %s", url, timeout)

}

func main() {
	err := f1()
	log.SetPrefix("Err:")
	log.SetFlags(0)
	if err != nil {
		switch err.(type) {
		case *LoginError:
			log.Println(err)
		default:
			log.Println("unkonw Err:", err)
		}
	}
	err = WaitForServer("https://www.youtube.com")
	if err != nil {
		log.Println(err)
	}

}
