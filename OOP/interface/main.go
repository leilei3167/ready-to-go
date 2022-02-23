package main

import "fmt"

type dog struct {
	color string
}
type animal interface {
	bark()
}

func (d dog) bark() {
	fmt.Println("汪汪汪")
}

func main() {
	var a animal
	wangcai := dog{color: "red"}
	wangcai.bark()
	a = wangcai
	a.bark()
}
