package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().Unix())
	number := rand.Intn(2)

	if number == 1 {
		time.Sleep(time.Second * 10)
		fmt.Fprintln(w, "慢响应")

	}

	fmt.Fprintln(w, "快响应")

}

func main() {

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)

}
