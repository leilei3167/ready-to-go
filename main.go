package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	ip, err := net.LookupHost("www.2szy.cn")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", ip)

	ips, err := net.LookupIP("www.qq.com")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", ips)

	addr, err := net.LookupAddr("121.14.77.221")

	fmt.Printf("%#v\n", addr)

}
