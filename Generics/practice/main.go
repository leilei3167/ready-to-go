package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func main() {
	/*
		调用一个泛型函数:
			实例化(instantiation):提供类型参数(此处为int),分为两步:1.编译器会将整个泛型函数或类型中替换成给定的类型参数;2.编译器
		验证每个类型参数是否满足自身的约束条件

	*/
	x := Min[int](1, 2)
	fmt.Println(x)
	x1 := Min(2.222, 333.2) //类型参数可以省略,由编译器自动推断
	fmt.Println(x1)
}

/*
一.Type Parameters 类型参数:
	函数和类型现在被允许拥有类型参数,类型参数被方括号包裹[ ]
*/

func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

/*
类型参数可以和类型一起使用:
在这里，通用类型Tree存储了类型参数T的值。通用类型可以有方法，比如本例中的Lookup。为了使用一个泛型，它必
须被实例化；Tree[string]是一个用类型参数string来实例化Tree的例子。
*/

type Tree[T any] struct {
	left, right *Tree[T]
	value       T
}

func (t *Tree[T]) Lookup(x T) *Tree[T] {
	fmt.Println("hello world")
	return nil
}

var stringTree Tree[string] //如果要初始化一个泛型,需要给定类型参数

/*
二.Type Sets 类型集:
	对于一个普通的函数,每一个参数都有一个类型,该类型定义了一组值,如int类型的参数,那么允许的参数值集合就
是int的集合;同样的,类型参数列表中每个类型参数都有一个类型,因为类型参数本身就是一个类型，所以类型参数的类
型定义了类型的集合。这种元类型被称为类型约束。

	在Min中的类型参数中,其类型约束是constraints.Ordered,这些类型的值可以被排序(可用<,==,>等比较),
	该约束确保只有可排序的值的类型才能被传递到Min函数中,也就是说Min的函数体中,该类型参数一定是可排序的

必须要注意的一点是:
	类型约束必须是一个接口!

现在接口的定义是一个类型集(type set)而不是之前的方法集(method set),即,实现了这些方法的类型的集合
类型集的观点比方法集的观点更有优势：我们可以明确地将类型添加到集合中，从而以新的方式控制类型集。

接口的语法得到了拓展,以下的someType接口顶一个包含int,string,bool这三种类型的类型集,也就是说:
这个接口只被int,string,bool所满足
'|'代表的是联合(或),对于类型约束,我们通常不关心一个特定的类型,如字符串,表达式~string意味着底层类型为
string的所有类型的集合。这包括string类型本身，以及所有用**定义声明的类型，如type MyString string **。


*/

type someType interface {
	int | string | bool
}

//当作类型约束使用的接口可以向上面一样带有名称,也可以用字面量的形式
type S[T someType] struct {
	Name string
}

type S1[T ~[]E, E interface{}] struct { //表示 S1[T interface{~[]E,E interface{}}]
	Name T
}

/*
三.Type inference 类型推理

3.1 函数参数类型推断:
x := Min[int](1, 2)
	fmt.Println(x)
	x1 := Min(2.222, 333.2) //类型参数可以省略,由编译器自动推断
	fmt.Println(x1)

这种从函数的参数类型中推断出参数类型的推理，被称为函数参数类型推理。
类型推理只适用于参数列表的推断,对于没有参数但却有返回值的函数,类型推断不适用


3.2 约束类型推理:

约束类型推理是从类型参数约束中推断出类型参数。当一个类型参数有一个定义在另一个类型参数上的约束时，它就会
被使用。当这些类型参数中的一个的类型参数是已知的，该约束被用来推断另一个的类型参数。

通常的情况是，当一个约束对某些类型使用~类型的形式时，该类型是用其他类型参数写的。我们在Scale的例子中看到
了这一点。S是~[]E，它是~后面的类型[]E用另一个类型参数来写。如果我们知道S的类型参数，我们就可以推断出E的
类型参数。S是一个片断类型，而E是该片断的元素类型。



*/
type Point []int

func ScaleAnd(p Point) {
	r := Scale(p, 2)
	fmt.Println(r)
}

func Scale[S ~[]E, E constraints.Integer](s S, c E) []E {
	r := make([]E, len(s))
	for i, v := range s {
		r[i] = v * c
	}
	return r
}
