package main

import (
	"fmt"
	"log"
	"net/url"
)

func main() {
	toCheck := "//mat1.gtimg.com/www/icon/favicon2.ico"

	u, err := url.Parse(toCheck)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v\n", u)

	fmt.Println(u.IsAbs())

}
