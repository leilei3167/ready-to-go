package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex
var count int

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	//打印Request中的method,URL,Proto
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	//Header是一个map[string][]string
	for k, v := range r.Header {
		//%q双引号围绕的字符串
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	//获取查询字段需要Parse
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	//r.Form是map[string][]string
	//Form contains the parsed form data, including both the URL
	// field's query parameters and the PATCH, POST, or PUT form data.
	// This field is only available after ParseForm is called.
	// The HTTP client ignores Form and uses Body instead.
	for k, v := range r.Form {
		//如果只想获取请求的实体如POST中的表单信息,使用r.PostForm
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

// counter echoes the number of calls so far.
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)

	mu.Unlock()
}
