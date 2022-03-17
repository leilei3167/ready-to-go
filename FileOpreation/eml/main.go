package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

func main() {
	s, _ := os.Open("test.eml")
	defer s.Close()
	r := bufio.NewReader(s)
	buf := make([]byte, 2048)
	for {
		n, err := r.Read(buf) //n是每次读取的字节数,为0时说明读取完毕
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

	}

	re := regexp.MustCompile(`^Sub.*`)

	//t := "SubjecL:leilei"
	t1 := re.FindString(string(bs))
	fmt.Println(t1)

}
