package main

import (
	"fmt"
	"os"
)

func main() {
	//判断是否是rootyonghu
	if os.Geteuid() > 0 {
		panic("permission denied!")

	}
	fmt.Println("yeahhh!")

}
