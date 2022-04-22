package benchmark

import (
	"fmt"
	"strings"
	"testing"
)

func benchmark(b *testing.B, f func(int, string) string) {
	//生成10个长度的字符串
	var str = randomString(10)
	//拼接一万次
	for i := 0; i < b.N; i++ {
		f(10000, str)
	}
}

func BenchmarkPlusConcat(b *testing.B)    { benchmark(b, plusConcat) }
func BenchmarkSprintfConcat(b *testing.B) { benchmark(b, sprintfConcat) }
func BenchmarkBuilderConcat(b *testing.B) { benchmark(b, builderConcat) }
func BenchmarkBufferConcat(b *testing.B)  { benchmark(b, bufferConcat) }
func BenchmarkByteConcat(b *testing.B)    { benchmark(b, byteConcat) }
func BenchmarkPreByteConcat(b *testing.B) { benchmark(b, preByteConcat) }

//builder buffer []string 都是成倍数扩容,实际底层都是切片,但是builder性能更高
//而+则是一次性申请内存
//测试一下builder内部的扩容
//16 32 64 128 256 512 896 1408 2048 3072 4096 5376 6912 9472 12288 16384 21760 28672 40960 57344 73728 98304 131072

func TestBuilderConcat(t *testing.T) {
	var str = randomString(10)
	var builder strings.Builder
	cap := 0
	for i := 0; i < 10000; i++ {
		if builder.Cap() != cap {
			fmt.Print(builder.Cap(), " ")
			cap = builder.Cap()
		}
		builder.WriteString(str)
	}
}

func TestGenerateToken(t *testing.T) {
	fmt.Println(GenerateToken(100))

}
