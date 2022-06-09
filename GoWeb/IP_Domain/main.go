package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ip := "182.61.25.124"

	hostname, err := net.LookupAddr(ip)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hostname)

}
