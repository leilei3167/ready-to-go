package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	a := "123456789"
	b := make([]byte, 3)

	n := copy(b, []byte(a))

	b = b[:n]

	fmt.Printf("a的长度:%d 真实:%#v,len:%d,cap:%d\n", len(a), b, len(b), cap(b))

	p := Person{
		Name: "大家扫",
		Age:  18,
		Job: JobInfo{
			Company: "腾讯",
			//	Salary:  10000,
			Work: map[string]interface{}{"职责": "研发", "工作地点": "北京"},
		},
	}

	byt, _ := json.Marshal(p)

	var p1 interface{}
	json.Unmarshal(byt, &p1)

	fmt.Printf("%#v\n", p1)

}

type Person struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`

	Job JobInfo `json:"job,omitempty"`
}
type JobInfo struct {
	Company string                 `json:"company,omitempty"`
	Salary  int                    `json:"salary,omitempty"`
	Work    map[string]interface{} `json:"work,omitempty"`
}
