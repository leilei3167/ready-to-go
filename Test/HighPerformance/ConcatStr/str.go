package benchmark

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
)

//使用benchmark测试一下常用的字符串拼接方式的性能差异,最终+性能最差,实际工作中尽量避免
//在知道长度的情况下,使用builder得到最高性能,在实际工作中,生产固定位数的随机token应使用builder

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//生成随机字符串来测试
func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

//直接用+号拼接
func plusConcat(n int, str string) string {
	s := ""
	for i := 0; i < n; i++ {
		s += str
	}
	return s
}

//sprintf
func sprintfConcat(n int, str string) string {
	s := ""
	for i := 0; i < n; i++ {
		s = fmt.Sprintf("%s%s", s, str)
	}
	return s
}

//使用builder来拼接
func builderConcat(n int, str string) string {
	var builder strings.Builder
	//如果事先知道长度 可使用grow
	builder.Grow(n * len(str))
	for i := 0; i < n; i++ {
		builder.WriteString(str)
	}
	return builder.String()
}

//bytes.buffer
func bufferConcat(n int, s string) string {
	buf := new(bytes.Buffer)
	//buf.Grow(n * len(s)) //即使设置了容量 性能也不如builder
	for i := 0; i < n; i++ {
		buf.WriteString(s)
	}
	return buf.String()
}

//用字节切片
func byteConcat(n int, str string) string {
	buf := make([]byte, 0)
	for i := 0; i < n; i++ {
		buf = append(buf, str...)
	}
	return string(buf)
}

//如果事先知道长度
func preByteConcat(n int, str string) string {
	buf := make([]byte, 0, n*len(str)) //事先知道容量会减少分配
	for i := 0; i < n; i++ {
		buf = append(buf, str...)
	}
	return string(buf)
}

func GenerateToken(n int) string {
	var builder strings.Builder
	builder.Grow(n) //如按字节生产字符串

	for i := 0; i < n; i++ {
		builder.WriteByte(letterBytes[rand.Intn(len(letterBytes))])
	}
	return builder.String()
}
