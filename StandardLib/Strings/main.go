package main

import (
	"fmt"
	"strings"
)

func main() {
	a := "123-123-432"
	b := strings.SplitN(a, "-", 2)
	fmt.Printf("%#v,len:%d,cap:%d", b, len(b), cap(b))

}
