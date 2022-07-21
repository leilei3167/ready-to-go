package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		form := map[string][]string(r.PostForm)
		log.Printf("Got form: %v,len:%d	", form, len(form))

		for k, v := range form {
			fmt.Println(k, v)
		}
		fmt.Fprintln(w, "done")
	}))

	http.ListenAndServe(":8080", nil)
}
