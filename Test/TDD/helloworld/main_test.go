package main

import "testing"

func TestHello(t *testing.T) {

	assertCorrectMessage := func(t testing.TB, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}

	t.Run("to a person", func(t *testing.T) {
		got := Hello("Chris", "")
		want := "Hello, Chris"
		assertCorrectMessage(t, got, want)
	})

	t.Run("empty string", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World"
		assertCorrectMessage(t, got, want)
	})

	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Elodie", spanish)
		want := "Hola, Elodie"
		assertCorrectMessage(t, got, want)
	})

	t.Run("in French", func(t *testing.T) {
		got := Hello("Lauren", french)
		want := "Bonjour, Lauren"
		assertCorrectMessage(t, got, want)
	})

}

/*
测试驱动开发:
	每次改动都应该运行测试函数,根据错误的结果进行检查,再重构使得代码输出无错误,不断的循环

t.Errorf 使测试失败并打印结果

t.Run 开启一个子测试,可以在这个子测试中运行其他测试,用于区分不同的case

在自测试中可能会出现重复代码,可以引入一个辅助函数来避免大量重复代码
对于辅助函数，接受 testing.TB 是一个好主意，它是 *testing.T 和 *testing.B 都满足的接口，因此您可以从测试或基准中调用辅助函数。
t.Helper告诉编译器这是一个辅助函数,在测试结果出错时他会将错误定位至自测试出错的位置,而不是在辅助函数中


1.写一个测试
2.使得编译器通过
3.编写使得测试能够运行的最小代码量,运行抛出错误并检查错误
4.编写足够的代码,使得测试能够通过
5.重构代码(期间要不断的运行测试,确保之前的用例不被破坏)

测试驱动开发就是不断的重复以上的过程!

*/
