package rand

import (
	"math/rand"
	"time"
)

//创建一个包内的全局变量,避免使其和标准库的rand进行隔离,避免其他调用时影响seed
//如果我们使用rand.Seed(time.Now())之后,其他调用者调用了rand.seed(1)那么,我们的seed将会被覆盖! 因此
//此处我们单独构建一个rand.Rand
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

//StringWithCharset 从给定的字符集中,生成length长度的随机字符串
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]

	}
	return string(b)

}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func String(length int) string {
	return StringWithCharset(length, charset)
}
