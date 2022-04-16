package main

import (
	"fmt"
)

/* 方法使得我们给Data扩展行为,而interface是我们能够实现多态
多态意味着:一段代码会根据他所在操作的数据而改变行为,解耦是对行为的关注,而data驱动着行为

什么时候一个Data应该具有行为呢?
	需要实现多态的时候(一段代码处理多种数据时)
*/
//接口定义的是行为(一定是动词),接口是抽象的,一些行为的声明,要理解接口不是真实的
//接口类型没有确切的值
type reader interface {
	read(b []byte) (int, error) //由调用者分配,然后将其分享
	//read(n int) ([]byte, error) 为什么这个行为很糟糕?返回[]byte将使得在方法内部分配内存(make切片),并在return时会造成底层数组的复制分配,不能在调用栈有指针
}

type file struct {
	name string
}

//file用值语义实现reader接口!
func (file) read(b []byte) (int, error) {
	s := "假装读取文件"
	copy(b, s)
	return len(b), nil
}

type pipe struct {
	name string
}

//pipe
func (pipe) read(b []byte) (int, error) {
	s := "假装读书"
	copy(b, s)
	return len(b), nil
}

//需要的是任何数据,只要他具有reader中包含的行为,而不是要一个指针类型的值!!
func retrieve(r reader) error {
	data := make([]byte, 100)
	len, err := r.read(data) //通过接口执行read,先到table查询接口存储的数据类型的方法在哪里,再调用,体现间接性
	if err != nil {
		return err
	}
	fmt.Println(string(data[:len]))
	return nil
}

func main() {
	//---------------用值语义实现接口-------------
	r := file{name: "leilei"} //创建file和pipe
	p := pipe{name: "yanzghen"}

	//值语义调用retrieve,以下调用分别创建各自的副本来传递入函数
	//但是此刻在retrieve中,我们只有r(reader接口,他本身是没有值的,不是任何具体的东西)
	//当我们把数据副本传递进去时,接口和数据就有了一种关系,存储关系,我们将数据存进接口,存入接口时将使接口具体化!
	//接口由两部分组成一个itable(函数指针的矩阵和一个用于存储实现接口了的具体类型,table有两部分,第一部分是我们存在接口中的具体的数据的类型,第二部分是函数指针,指向该数据的方法)
	//在此调用retrieve时,r.read将会做一个table的查询,找到实现的位置,然后再调用该实现(方法),此刻就体现出了解耦的2个代价(间接性(接口调用方法是先到table查询),和分配(复制具体的值到接口))
	//但是这里也体现出了接口的解耦作用,接口的retrieve函数 能够接收任何类型,只要有read方法,并且会根据传入的数据不同,改变自己的行为(file就是读文件,pipe就是读书)

	retrieve(r) //reader接口储存r的file类型
	retrieve(p) //reader接口储存r的pipe类型
	//---------------------用指针语义实现接口-------------------
	user := User{
		name:  "leilei",
		email: "1356556043@qq.com",
	}
	ToNotify(&user) //接口是用指针语义实现的,而这里企图使用值语义调用,那么就是切换语义

	//Torecive(user)
}

type notifier interface {
	send()
	//recived()
}

type User struct {
	name  string
	email string
}

//实现notifier接口,使用指针语义(有一个指针语义就算),代表只能和接口分享数据(而不是拷贝)
func (u *User) send() {
	u.name = u.name + "sended"
	fmt.Printf("has sended:user:%v to: %v\n", u.name, u.email)
}

/* func (u User) recived() {
	fmt.Println("recived msg from", u.email)
} */

func ToNotify(n notifier) {
	n.send()
	
}



/* func Torecive(n notifier) {
	n.recived()
} */
//接口的itable有一个方法集
//操作的是指针时如上述的ToNotify(&user),指针接收者和值接收者的方法都可以调用,方法集完整,实现接口,就可正常调用
//但是操作的是值时(ToNotify(user)),指针语义的方法被排除在方法集之外,则方法集不完整,无法实现接口,无法调用;
//编译器这样设定是为了保证一致性,语义的一致性,两个主要原因

//1.调用指针接收者意味着共享,需要传入指针,但并非所有的值都能取地址,如常量,无法取地址代表无法分享,自然无法调用指针接收者的方法
//不能假设总是能取得值的地址,但是如果操作的是一个指针,那么就说明已经能够取得地址就说明他是能够被分享的或者其指向的地址的值是能够拷贝的,
//因此操作指针能够 调用指针接收者和值接收者的方法
//2.从行为来看,选择了指针语义,分享;选择值语义,则是拷贝;
//当选择了指针语义 指针接收者实现方法,就只能够分享(此时用值调用指针接收者的方法(改变语义),就意味着你想拷贝指针指向的值,这是被禁止的,不能假设那里是安全的)

//当用值语义实现接口时,我们可以拷贝,在必要的时候也可以进行共享,因此值语义实现的方法,可以传入值来拷贝,以及传入指针来共享;但是为了保证语义一致性,建议始终使用值语义
//而指针语义实现的接口,调用时只能够分享(传入指针)
