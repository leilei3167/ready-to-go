package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		_, err := w.Write([]byte("hello world"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "internal err")
		}

	})

	http.ListenAndServe(":8080", nil)

}
