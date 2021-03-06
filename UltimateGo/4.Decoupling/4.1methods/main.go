package main

import "fmt"

/*
在调用方法时,go只关心数据,如*user和user都是数据,不关心是指针还是值,go会自动转换来满足接收者
不要使用a:=&user{...}这类表达,而应该直接创建值,在传参时使用&a来表示传入地址



方法可以有指针接收者和值接收者,如何选择取决于你的语义,切勿混合(有例外) ,要明白每个选择背后的成本

语义一致性就是一切!!!

在Go中,实际只有3中类型:

build-in类型: 如 int bool等等,只使用值语义!!!!不要使用任何指针语义 如*int *bool *string,包括作为结构体字段时(除非在极其特殊的情况下,如和关系型数据库交互时)
reference-type引用类型:必须使用值语义!!!严禁使用a:=&[]string,*[]string 之类的表达,除非是在unmarshal 或者decode时
defined-type 自定义类型:需要自己决定要使用的语义,不清楚时就使用指针语义!!!如果确定能使用值语义,就使用值语义!!发现错误再重构;方法选择指针接收者或值接收者,取决于你选择的语义
而不是是否要在方法内部修改值!!
如os包中的*file,很多方法并没有修改file的字段!
方法,api,代码 必须尊重语义,语义一致性永远是第一位,值接收者还是指针接收者,取决于语义



指针语义改变一个已有的数据而值语义是拷贝后在副本基础上修改,创造一个
思考:1.假如有一个代表时间的值,我想要在上面加5秒钟,是在其本身修改增加5秒,还是说将其复制,在其副本上加五秒后返回?
2.假设一个人结构体,有名字年龄身份等字段,如果我要给他其别名,是在原身修改名字,还是重新复制他 仅仅修改名字后再返回副本?
很明显 后者并不符合我们思考的直觉,编码时应该考虑到这一点;在处理自定义类型时,要符合人的直觉,不是所有的东西都是能修改的!

在标准库或者第三方包中,着重的看工厂函数的返回值,那代表了开发者所选择的语义类型,返回指针就说明应该尽量避免对其指向的地址的拷贝(如range切片等操作导致)

*/

type User struct {
	Age int
}

func (u *User) AddAge(n int) { //值语义
	u.Age += n

}

func (u User) ReduceAge(n int) User {
	return User{
		Age: u.Age - n,
	}

}

type test User

func main() {
	//	var a test
	var b User

	//a.AddAge() test虽然基于User 但无法使用其方法

	b.AddAge(5)
	fmt.Printf("point:%v\n", b.Age)
	c := b.ReduceAge(10)
	fmt.Printf("reduce:%v\n", c)
	f1 := b.AddAge //值语义,f不能引用d,而是用的是b的副本
	f1(5)

	f2 := b.ReduceAge //指针语义
	f2(10)

}
