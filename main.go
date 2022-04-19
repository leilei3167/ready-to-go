package main

import (
	"fmt"
	"log"
)

func main() {
	StartPort := 12
	EndPort := 17
	var ports []int
	for i := StartPort; i < EndPort+1; i++ {
		ports = append(ports, i)
	}
	fmt.Println(ports)
	n := EndPort - StartPort + 1 //端口数
	var ports1 = make([]int, n)
	for i := 0; i < n; i++ {
		ports1[i] = StartPort
		StartPort++
	}
	fmt.Println(ports1)

	for k, v := range ports {
		if ports1[k] != v {
			log.Fatalln("err")
		}

	}
	fmt.Println("success")

}
