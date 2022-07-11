package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	file, err := os.Open("/home/lei/code/ready-to-go/MDtest.md")
	if err != nil {
		log.Fatal(err)
	}

	count := lineCounter(file)
	fmt.Println(count)

}

//count the number of line in a file
func lineCounter(src io.Reader) int {
	buf := make([]byte, 32*1024)

	count := 0
	sep := []byte{'\n'}

	for {
		c, err := src.Read(buf)
		count += bytes.Count(buf[:c], sep)

		switch {
		case err == io.EOF:
			return count

		case err != nil:
			return -1
		}
	}

}
