package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//获取请求header
/* Header 是一个map[string][]string的map */
func header(w http.ResponseWriter, r *http.Request) {
	h := r.Header
	fmt.Fprintln(w, "直接整个读取header:", h)
	fmt.Fprintln(w, "通过字段读取header:", r.Header["Accept-Encoding"])
	fmt.Fprintln(w, "通过Get读取header:", r.Header.Get("Accept-Encoding"))

}

//获取请求主体
func body(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	fmt.Fprintln(w, "ReadAll body中的内容:", string(data))
	//或者根据长度来读取,body只能被读取一次
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	fmt.Fprintln(w, "用body的Read方法:", string(body))

}

//处理Form
func process(w http.ResponseWriter, r *http.Request) {
	//当url查询字段和key和Form中的key重复且对应的值不同时,直接访问Form将会同时的到两个值,并且Form的总是在前
	//From为: map[leilei:[fromForm fromquery]]
	r.ParseForm()
	fmt.Fprintln(w, "From为:", r.Form)
	//如果只想要获取Form中的Key对应的值

	fmt.Fprintln(w, "From为:", r.PostForm)
	//From为: map[leilei:[fromForm]]
	//要注意Form只支持application/x-www-form-urlencoded编码,如果要取得form-data编码的表单数据:
	fmt.Fprintln(w, "FormValue:", r.FormValue("leilei")) //直接返回字符串结果

}

func file(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("test.txt")
	if err == nil {
		data, err := ioutil.ReadAll(file)
		if err == nil {
			fmt.Fprintln(w, string(data))
		}

	}

}

//当访问response,但我还没开发好时 回应状态码
func response(w http.ResponseWriter, r *http.Request) {
	//WriteHeader是写入状态码,要注意WriteHeader之后就不能对Header进行操作
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "还没搭建好...")

}

type Post struct {
	User   string
	Treads []string
}

func Json(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	post := &Post{
		User:   "leilei",
		Treads: []string{"lei", "dsa", "cxzcxz"},
	}
	//编码
	data, _ := json.Marshal(post)
	w.Write(data)

}

//未设置过期时间,默认为会话cookie,在浏览器关闭之后cookie会消失
func cookie(w http.ResponseWriter, r *http.Request) {
	c1 := http.Cookie{
		Name:     "first_cookie",
		Value:    "Go Web",
		HttpOnly: true,
	}
	c2 := http.Cookie{
		Name:     "second_cookie",
		Value:    "Go Web222222",
		HttpOnly: true,
	}
	w.Header().Add("Set-Cookie", c1.String())
	w.Header().Add("Set-Cookie", c2.String())
	//一个更简便的方式
	c3 := http.Cookie{
		Name:     "3 third cookie",
		Value:    "power by SetCookie",
		HttpOnly: true,
	}

	http.SetCookie(w, &c3)

}
func getcookie(w http.ResponseWriter, r *http.Request) {
	//直接从Header中获取,返回的是字符串切片
	h := r.Header["Cookie"]
	fmt.Fprintln(w, "get cookie:", h)
	//用Cookie,获取指定cookie
	c1, err := r.Cookie("first_cookie")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "找不到cookie")
	}
	//获取所有cookie,和直接Header["Cookie"]结果一致
	cs := r.Cookies()
	fmt.Fprintln(w, c1)
	fmt.Fprintln(w, cs)

}

func setMessage(w http.ResponseWriter, r *http.Request) {
	//如果Value中有空格等特殊字符,则需要进行编码
	msg := []byte("hello world!!!!")
	c := http.Cookie{
		Name:  "flash",
		Value: base64.URLEncoding.EncodeToString(msg),
	}
	http.SetCookie(w, &c)

}

func getMessage(w http.ResponseWriter, r *http.Request) {
	//获取cookie
	c, err := r.Cookie("flash")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprintln(w, "there is no cookie yet")
		}

	} else {
		//设置同名cookie 覆盖设置,将MaxAge设置为-1表示已过期,删除这个cookie,过期时间也设置为过去的时间
		rc := http.Cookie{
			Name:    "flash",
			MaxAge:  -1,
			Expires: time.Unix(1, 0),
		}
		http.SetCookie(w, &rc)
		//解码获得cookie中的value
		val, _ := base64.URLEncoding.DecodeString(c.Value)
		fmt.Fprintln(w, "解码cookie:", string(val))

	}

}

func main() {

	http.HandleFunc("/header", header)       //获取请求头
	http.HandleFunc("/body", body)           //获取请求体
	http.HandleFunc("/process", process)     //获取表单信息
	http.HandleFunc("/file", file)           //获取被上传到服务器的文件
	http.HandleFunc("/response", response)   //构建Response
	http.HandleFunc("/Json", Json)           //Json回应
	http.HandleFunc("/cookie", cookie)       //设置cookie
	http.HandleFunc("/getcookie", getcookie) //获取cookie

	//用cookie实现闪现消息
	http.HandleFunc("/setflash", setMessage)
	http.HandleFunc("/getflash", getMessage)//设置cookie后 getcookie之后刷新cookie即消失 

	http.ListenAndServe(":8081", nil)
}
