package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var wg sync.WaitGroup
var threads int
var client http.Client

func main() {
	var target, dict string //必填
	var verbose bool
	var wait, timeout int

	flag.StringVar(&target, "t", "", "目标域名")
	flag.StringVar(&dict, "d", "", "字典目录")
	flag.IntVar(&wait, "w", 0, "请求之间的等待时间")
	flag.IntVar(&threads, "T", 10, "并发数 : 10")
	flag.IntVar(&timeout, "to", 3, "连接超时")
	flag.BoolVar(&verbose, "v", false, "详细模式")
	flag.Parse()

	if target == "" || dict == "" {
		fmt.Println("You didn't provide enough arguments. Refer to README.md to have the usage detail.")
		return
	}

	client.Timeout = time.Duration(timeout) * (time.Second)
	start := time.Now()

	list := getList(dict)
	fmt.Printf("\nTARGET : %s\n", target)
	fmt.Printf("DICT : %s\n", dict)
	fmt.Println("START TIME : " + time.Now().Format("15:04:05"))
	fmt.Printf("THREADS : %d\n", threads)
	fmt.Printf("-- Threads init\n\n")

	for i := 0; i < threads; i++ {
		//将任务分片
		begin := (len(list) / threads) * i
		end := (len(list) / threads) * (i + 1)
		wg.Add(1)
		go checkURL(list[begin:end], target, verbose, wait)
		logrus.Infof("thread %d created", i)

	}
	fmt.Printf("\n-- Scan started\n\n")

	// Wait for all the goroutines to end
	wg.Wait()

	elapsedTime := time.Since(start)
	fmt.Printf("\n-- Scan terminated in %v\n", elapsedTime)

}

// contact sends a request to a specified target
func contact(target string) (int, error) {
	resp, err := client.Get(target) //example:https://www.baidu.com/contact/index
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

// displayResults only processes the statusCode to display the result in a specific color
func displayResult(statusCode int, target, url string, v bool) {
	if statusCode >= 400 && statusCode <= 499 && v == true {
		color.Red("%v : %s is not present\n", statusCode, target+url)
	} else if statusCode >= 200 && statusCode <= 299 {
		color.Green("%v : %s\n", statusCode, target+url)
	} else if statusCode >= 500 && statusCode <= 599 && v == true {
		color.Magenta("%v : %s respond internal server error\n", statusCode, target+url)
	}
}
func checkURL(givenList []string, target string, verbose bool, wait int) {
	defer wg.Done()

	for _, url := range givenList {
		if strings.HasPrefix(url, "/") {
			url = "/" + url
		}

		statusCode, err := contact(target + url)
		if err != nil && verbose == true {
			logrus.Warnf("an error occured : %v\n", err)
		}
		displayResult(statusCode, target, url, verbose)
		if wait != 0 {
			time.Sleep(time.Duration(wait) * time.Millisecond)
		}
	}

}

//打开字典目录,按行解析为字符串切片
func getList(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var list []string

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "#") || scanner.Text() == "" {
			continue
		}
		list = append(list, scanner.Text())
	}
	file.Close()
	return list
}
