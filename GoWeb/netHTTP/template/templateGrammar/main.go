package main

import (
	"html/template"
	"net/http"
)

//模板中的变量赋值
func fun1(w http.ResponseWriter, r *http.Request) {
	tem, _ := template.ParseFiles("test1.html")
	tem.Execute(w, nil)

}

//模板中的if选择
func fun2(w http.ResponseWriter, r *http.Request) {
	tem, _ := template.ParseFiles("test2.html")

	tem.Execute(w, nil)

}

//模板中的循环range
func fun3(w http.ResponseWriter, r *http.Request) {
	tem, _ := template.ParseFiles("test3.html")
	s := []string{"你好", "我好", "她也好"}
	tem.Execute(w, s)

}

//模板中的pipeline
func fun4(w http.ResponseWriter, r *http.Request) {
	tem, _ := template.ParseFiles("test4.html")
	tem.Execute(w, nil)

}

//模板中的with...end
func fun5(w http.ResponseWriter, r *http.Request) {
	tem, _ := template.ParseFiles("test5.html")
	s := []string{"你好", "我好", "她也好"}
	tem.Execute(w, s)

}

//difine和template 模板嵌套和莫版内互相引用
func fun6(w http.ResponseWriter, r *http.Request) {
	//要将涉及到的全部解析
	tem, _ := template.ParseFiles("index.html", "header.html", "footer.html")
	tem.Execute(w, nil)

}
func main() {

	http.HandleFunc("/var", fun1)
	http.HandleFunc("/if", fun2)
	http.HandleFunc("/loop", fun3)
	http.HandleFunc("/pipeline", fun4)
	http.HandleFunc("/with", fun5)
	http.HandleFunc("/define", fun6)
	http.ListenAndServe(":8080", nil)
}
