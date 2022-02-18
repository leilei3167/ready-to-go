package main

import (
	"log"
	"net/http"
	"text/template"
)

type User struct {
	Name   string
	gender string
	Age    int
}

func index(w http.ResponseWriter, r *http.Request) {
	//定义模板

	//解析模板
	t, err := template.ParseFiles("./hello.html")
	if err != nil {
		log.Fatalln(err)
	}
	//渲染模板,传入结构体
	u1 := User{
		Name:   "雷磊",
		gender: "男", //小写不对外开放
		Age:    18,
	}
	m1 := map[string]interface{}{
		"Name":   "雷磊1",
		"gender": "男1",
		"Age":    181,
	}
	//如果要传入多个复杂的复合类型,可以将参数设置为key为string,value为空接口的map即可
	t.Execute(w, map[string]interface{}{
		"u1": u1,
		"m1": m1,
	}) //将结构体替换到模板中, .就代表了u1, {{.Name}}
}
func main() {

	http.HandleFunc("/", index)
	http.ListenAndServe(":9090", nil)
}
