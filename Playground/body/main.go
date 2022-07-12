package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		log.Printf("read body len:%d string: %s", len(b), string(b))

		//body只会被读取一次,如果多次读取,需要重新保存

		r.Body = io.NopCloser(bytes.NewBuffer(b))

		Crete(w, r)

	})
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.ListenAndServe()

}

func Crete(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	log.Printf("Crete read body: %d", len(b))
	Crete2(w, r)
}

func Crete2(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	log.Printf("Crete2 read body: %d", len(b))

}
