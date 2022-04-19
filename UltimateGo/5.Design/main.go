package main

/* 类型断言 */
import "fmt"

type Mover interface {
	Move()
}
type Locker interface {
	Lock()
	Unlock()
}
type MoveLocker interface {
	Mover
	Locker
}

type bike struct{}

func (bike) Move() {
	fmt.Println("Moving the bike")
}
func (bike) Lock() {
	fmt.Println("Locking the bike")
}
func (bike) Unlock() {
	fmt.Println("Unlocking the bike")
}

func main() {
	var ml MoveLocker
	var m Mover
	//因为bike实现了所有三个接口,可以存入ml接口
	ml = bike{}
	//并不是在将ml接口本身的值赋值给m,而是将ml中的值(bike)存入m
	m = ml
	fmt.Println(m)

	//然而当把m重新赋值给ml时确报错;因为编译器只能确定在m中的值是可以Move的,但是不确定其是否能够Lock
	//ml = m
	//cannot use m (variable of type Mover) as MoveLocker value in assignment: Mover does not implement MoveLocker (missing method Lock)

	/* 类型断言
			有一种方法能够在运行时测试赋值是否合法,并让其实现,就要使用类型断言
			类型断言使得runtime能够提问,is there a value of the given type stored inside the interface?
			b:=m.(bike)
			ml=b
			代表着,在这条code执行的时候,m中是否储存着bike这样的值
			如果有,将会将b的副本将可以存入ml
		使用的更广泛的是comma ok的形式
		b,ok:=m.(bike)

		In this form, if ok is true, there is a bike value stored inside of the interface. If ok is
	false, then there isn’t and the program does not panic. The variable b however is
	still of type bike, but it is set to its zero value state.
	*/
	b, ok := m.(bike)
	if !ok {
		fmt.Printf("%T\n", b)
	}
	fmt.Printf("%T\n", b) 

	/* Interface pollution
	大多数人开发软件使用接口,而不是去发现接口;
	应该现实现一个具体的解决方案,然后再在其中去发现,哪些地方是要求多态的,然后再使用接口解耦
	几个典型的接口污染:
		1.一个包声明了一个能够和被他具体类型全部实现的接口
		2.接口类型可导出,但是实现他的具体类型不可导出
		3.工厂函数返回的接口(带有不可导出的值)
		4.可以被删除但却对api没有任何影响的接口
		5.没有解耦作用的接口

	一些指导意见:
		使用接口:
			当api需要提供一个实现的细节时
			当api有多个实现需要维护时
			当api可以改变,并且确定需要解耦

		质疑一个接口:
			当他存在的目的完全是为了写可测试的代码时(优先写可用的api)
			当他没有解耦作用时
			当接口对代码的改善作用并不明确时
	*/


}
