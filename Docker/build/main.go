package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second * 2)
	var a time.Time
	for v := range ticker.C {
		fmt.Printf("Hello World! the Time is [%v]\n", time.Now())
		a = v
	}
	log.Panicln(a)

}
