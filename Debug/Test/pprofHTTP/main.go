package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

var datas []string

func Add(str string) string {
	data := []byte(str)
	sData := string(data)
	datas = append(datas, sData)

	return sData
}

//访问localhost:6060/debug/pprof即可
//使用 go tool pprof -http=:8080 profile 即可看到火焰图
func main() {
	go func() {
		for {
			log.Println(Add("https://github.com/EDDYCJY"))
		}
	}()

	http.ListenAndServe("0.0.0.0:6060", nil)
}
