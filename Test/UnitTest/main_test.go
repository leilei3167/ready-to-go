package UnitTest

//包名必须以_test结尾
import (
	"reflect"
	"testing"
)

//必须以Test开头
//用go test -v运行
func TestSplit(t *testing.T) { // 测试函数名必须以Test开头，必须接收一个*testing.T类型参数
	got := Split("a:b:c", ":")         // 程序输出的结果
	want := []string{"a", "b", "c"}    // 期望的结果
	if !reflect.DeepEqual(want, got) { // 因为slice不能比较直接，借助反射包中的方法比较
		t.Errorf("expected:%v, got:%v", want, got) // 测试失败输出错误提示
	}
}

func TestSplitWithMore(t *testing.T) {
	got := Split("abcd", "bc")
	want := []string{"a", "d"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected:%v, got:%v", want, got)
	}
}

//通过单元测试获得结果之后需修改bug,修改之后必须运行所有的单元测试(回归测试)

//子测试,以上是两种最简单的测试,只有一组数据,支持在测试函数中使用t.Run执行一组测试用例，
//这样就不需要为不同的测试数据定义多个测试函数了。
func TestSplitChildren(t *testing.T) {
	t.Run("case1", func(t *testing.T) {
		got := Split("a:b:c", ":")         // 程序输出的结果
		want := []string{"a", "b", "c"}    // 期望的结果
		if !reflect.DeepEqual(want, got) { // 因为slice不能比较直接，借助反射包中的方法比较
			t.Errorf("expected:%v, got:%v", want, got) // 测试失败输出错误提示
		}
	})
	t.Run("case2", func(t *testing.T) {
		got := Split("a:b:c,d", ",")       // 程序输出的结果
		want := []string{"a:b:c", "d"}     // 期望的结果
		if !reflect.DeepEqual(want, got) { // 因为slice不能比较直接，借助反射包中的方法比较
			t.Errorf("expected:%v, got:%v", want, got) // 测试失败输出错误提示
		}
	})

}

//表格驱动测试,是一种编写更清晰测试的一种方式和视角,表格的每一条都是一个完整的测试用例
//通常表格是匿名结构体切片
func TestSplitAll(t *testing.T) {
	//设置需测试表格
	tests := []struct {
		name  string
		input string
		sep   string
		want  []string
	}{
		{"case1", "a:b:c", ":", []string{"a", "b", "c"}},
		{"case2", "a+b+c/d", "+c", []string{"a+b", "/d"}},
		{"case3", "a::c", ":", []string{"a", "c"}},
		{"case4", "雷磊你好", "磊", []string{"雷", "你好"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Split(test.input, test.sep)
			if !reflect.DeepEqual(test.input, test.want) {
				t.Errorf("%s err!want:%s got:%s\n", test.name, test.want, got)
			}
		})
	}

}
