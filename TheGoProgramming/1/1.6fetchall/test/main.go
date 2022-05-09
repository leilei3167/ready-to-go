package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	url := "https://www.baidu.com"
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	n, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	time := time.Since(start).Seconds()
	log.Printf("time:%.2fs,read %d byte!", time, n)
}
