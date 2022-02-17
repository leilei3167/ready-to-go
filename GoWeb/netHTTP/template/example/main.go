package main

import (
	"log"
	"net/http"
	"text/template"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	//解析已定义的模板(事先创建好了html文件
	t, err := template.ParseFiles("./hello.html") //相对地址, .代表编译后可执行文件的地址
	if err != nil {
		log.Fatalln(err)
	}
	//渲染模板
	err = t.Execute(w, "雷磊") //将t渲染,并将data的数据填充进{{.}}中
	if err != nil {
		log.Fatalln(err)
	}
}
func main() {

	http.HandleFunc("/", sayHello)
	http.ListenAndServe(":9090", nil)
}
