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
}
