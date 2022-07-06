package generic

import "testing"

func TestAssertFunc(t *testing.T) {
	t.Run("asserting on integers", func(t *testing.T) {
		AssertEqual(t, 1, 1)
		AssertNotEqual(t, 1, 2)
	})

	t.Run("asserting on strings", func(t *testing.T) {
		AssertEqual(t, "hello", "hello")
		AssertNotEqual(t, "hello", "Grace")
	})

}

/*
泛型相比于空接口:
	接收空接口作为参数的函数,类型检查将毫无作用,对于使用者来讲可能出现person和airport比较的错误,这些在编译期
都是无法报错的(需要运行的时候通过反射来实现检查),而泛型,给到了像空接口类似的灵活性,但同时又限制了传入其中的类型的范围(保证了类型安全)

*/
/* //v1 1.18之前
func AssertEqual(t *testing.T, got, want interface{}) {
	t.Helper()

	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}

}

func AssertNotEqual(t *testing.T, got, want interface{}) {
	t.Helper()

	if got == want {
		t.Errorf("didn't want %+v", got)
	}

} */

//v2 使用泛型,需要在函数名称后提供类型参数 ,[ ]中声明类型为comparable的
func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()

	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}

}

func AssertNotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()

	if got == want {
		t.Errorf("didn't want %+v", got)
	}

}

/*

以下两个函数所执行的操作是否一致?

func GenericFoo[T any](x, y T)

func InterfaceyFoo(x, y interface{})

对于泛型的版本来说,x,y都必须是<相同的T类型>,即x,y传入的类型必须相同,而空接口版本对此没有限制
在函数返回泛型类型时,调用者可以按照原样进行使用,无需额外进行返回接口类型值时的类型断言


*/

func AssertTrue(t *testing.T, got bool) {
	t.Helper()

	if !got {
		t.Errorf("got %+v, want true", got)
	}

}

func AssertFalse(t *testing.T, got bool) {
	t.Helper()

	if got {
		t.Errorf("got %+v, want false", got)
	}

}

// 对于栈来讲,只是类型的不同而已,我们却要写两份几乎相同的代码,利用泛型就可以解决这个问题
func TestStack(t *testing.T) {
	t.Run("integer stack", func(t *testing.T) {
		myStackOfInts := new(Stack[int])

		// check stack is empty
		AssertTrue(t, myStackOfInts.IsEmpty())

		// add a thing, then check it's not empty
		myStackOfInts.Push(123)
		AssertFalse(t, myStackOfInts.IsEmpty())

		// add another thing, pop it back again
		myStackOfInts.Push(456)
		value, _ := myStackOfInts.Pop()
		AssertEqual(t, value, 456)
		value, _ = myStackOfInts.Pop()
		AssertEqual(t, value, 123)
		AssertTrue(t, myStackOfInts.IsEmpty())

		// can get the numbers we put in as numbers, not untyped interface{}
		myStackOfInts.Push(1)
		myStackOfInts.Push(2)
		firstNum, _ := myStackOfInts.Pop()
		secondNum, _ := myStackOfInts.Pop()
		AssertEqual(t, firstNum+secondNum, 3)
	})
}
