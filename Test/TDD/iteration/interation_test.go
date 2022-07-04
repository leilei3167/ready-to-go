package interation

import (
	"fmt"
	"testing"
)

func TestRepeat(t *testing.T) {
	repeat := Repeat("a", 5)
	expected := "aaaaa"

	if repeat != expected {
		t.Errorf("expected '%s' but got '%s'", expected, repeat)
	}
}

//go test -bench=. 运行基准测试
func BenchmarkRepeat(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Repeat("a", 5)
	}

}

func ExampleRepeat() {
	repeat := Repeat("a", 5)
	fmt.Println(repeat)
	// Output: aaaaa
}
