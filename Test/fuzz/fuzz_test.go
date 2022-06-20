package fuzz

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

//https://go.dev/doc/fuzz/#glossary 文档
func Fuzz_isEqual(f *testing.F) {
	f.Fuzz(func(t *testing.T, a []byte, b []byte) {
		isEqual(a, b)

	})
}

//即使是表驱动的单元测试,也无法提供尽可能多的测试用例,始终都会有一些不可预测的结果
//在模糊测试中无法预测预期输出，因为您无法控制输入
//但是可以在模糊测试中验证 Reverse 函数的一些属性

func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus,在不适用-fuzz标签的情况下,go test将会至少执行这些部分
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rev, err1 := Reverse(orig)
		if err1 != nil {
			return
		}
		doubleRev, err2 := Reverse(rev)
		if err2 != nil {
			return
		}
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}

//当模糊测试出现错误时,在testdata/fuzz/FuzzReverse 目录查看出错输入的值,并且再执行go test时会执行错误的case,利于回归测试

//模糊测试默认永远允许,直到被中断,-fuzztime 30s 可设置执行多长的时间

func FuzzSum(f *testing.F) {
	rand.Seed(time.Now().UnixNano())

	f.Add(10)
	f.Fuzz(func(t *testing.T, n int) {
		n %= 20
		var vals []int64
		var expect int64
		var buf strings.Builder
		buf.WriteString("\n")
		for i := 0; i < n; i++ {
			val := rand.Int63() % 1e6
			vals = append(vals, val)
			expect += val
			buf.WriteString(fmt.Sprintf("%d,\n", val))
		}

		assert.Equal(t, expect, Sum(vals), buf.String()) //将生成的数字写入到buf中记录,出错时会打印出来,然后对此编写新的单元测试
	})
}

func TestSumFuzzCase1(t *testing.T) {
	vals := []int64{
		622262,
		504498,
		328884,
		634493,
		96783,
		166009,
		862846,
		966464,
		760583,
		300000,
		616714,
		505432,
		983735,
		296814,
	}
	assert.Equal(t, int64(7645517), Sum(vals))
}
