package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)
//只有大写的字段才能被json编码
type User struct {
	Name string `json:"name"`
	Age  string `json:"age"`

	File map[string]interface{} `json:"file"`
}

func json1(w http.ResponseWriter, r *http.Request) {
	file := make(map[string]interface{})
	for i := 0; i < 7; i++ {
		s := "haha" + strconv.Itoa(i)
		file[s] = i

	}

	user := &User{
		Name: "leilei",
		Age:  "123",
		File: file,
	}

	r.Header.Set("Content-Type", "application/json")

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func main() {

	http.HandleFunc("/", json1)

	http.ListenAndServe(":8090", nil)

}
