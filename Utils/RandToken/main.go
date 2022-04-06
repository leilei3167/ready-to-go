package main

import (
	"fmt"
	"math/rand"
)

func main() {

	fmt.Println(Rand1(64))
	fmt.Println(Rand2(64))
	fmt.Println(Rand3(64))

}

//https://zhuanlan.zhihu.com/p/90830253
/* 最简单方法 */
func Rand1(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	//	rand.Seed(time.Now().UnixNano())
	res := make([]rune, n)

	for i, _ := range res {
		res[i] = letters[rand.Intn(len(letters))]

	}
	return string(res)

}

/* 优化,改进为字符串常量,len(Bytes)也是常量,效率提升;
rand.Intn也是间接用了rand.Int63,直接使用效率更高;
但为了防止超过索引,使用取余 */
const (
	Bytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits

)

func Rand2(n int) string {
	b := make([]byte, n)
	for i := range b {
		//防止Int63()的值超过索引
		b[i] = Bytes[rand.Int63()%int64(len(Bytes))]

	}
	return string(b)
}

//用掩码进行替换

func Rand3(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(Bytes) {
			b[i] = Bytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
