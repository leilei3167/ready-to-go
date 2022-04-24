package slice

import (
	"math/rand"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestDelete(t *testing.T) {
	a := []string{"a", "b", "c", "d"}
	got := Delete(a, 2)                // 程序输出的结果
	want := []string{"a", "b", "d"}    // 期望的结果
	if !reflect.DeepEqual(want, got) { // 因为slice不能比较直接，借助反射包中的方法比较
		t.Errorf("expected:%v, got:%v", want, got) // 测试失败输出错误提示
	}

}

//用于测试的辅助函数(生产随机数),打印内存
func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func printMem(t *testing.T) {
	t.Helper() //声明为辅助函数
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm) //记录内存消耗
	t.Logf("%.2f MB", float64(rtm.Alloc)/1024./1024.)
}

//这种写法特别适合比较,benmark也比较时候这种写法,将待测试的函数签名传入
func testLastChars(t *testing.T, f func([]int) []int) {
	t.Helper()
	ans := make([][]int, 0)
	for k := 0; k < 100; k++ {
		origin := generateWithCap(128 * 1024) // 1M
		ans = append(ans, f(origin))
		runtime.GC() //主动GC的话copy的使用量会降到更低
	}
	printMem(t)
	_ = ans

}

//内存消耗高达100M
func TestLastCharsBySlice(t *testing.T) { testLastChars(t, lastNumsBySlice) }

//内存消耗仅有1M出头
func TestLastCharsByCopy(t *testing.T) { testLastChars(t, lastNumsByCopy) }
