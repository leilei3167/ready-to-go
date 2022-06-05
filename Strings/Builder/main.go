package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const charset = "abcdefhijklmnopq123456789"

func main() {
	rand.Seed(time.Now().UnixNano())
	var builder strings.Builder
	for i := 0; i < 100; i++ {
		builder.WriteByte(charset[rand.Intn(len(charset))]) //高效拼接字符串
	}
	fmt.Println(builder.String())
	builder.Reset() //清空字符串的缓存

	var test = []string{"dsa", "dsacxz", "everything", "mcixz"}

	//拼接大量次数
	for i := 0; i < 100; i++ {
		for _, t := range test {
			builder.WriteString(t)
		}
		builder.WriteString("\n")
	}
	fmt.Println("拼接100次字符串切片:", builder.String())

}
