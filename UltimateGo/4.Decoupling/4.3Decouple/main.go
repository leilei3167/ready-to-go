package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"
)

/* 一个简单示例,理解解耦的过程和意义
1.一开始只关注手头的问题(业务的逻辑),将解耦,约定什么的抛到脑后,先将代码能够跑通业务逻辑,到时候再考虑,哪些可以"解耦",以便于可以应对变化,更具灵活性;将业务分解成一个一个的小问题,一次解决,最终最上层的代码就是将这些小问题串联,实现业务逻辑

*/
//业务模拟,将Xenia系统的数据迁移至系统...
//分解成小的任务步骤:连接到Xenia;获取到数据;连接到新的系统;放入之前拉取的数据;

/*一,构建原始层,Xenia有状态(连接),设置数据结构
这一层不要关注性能优化,这些只是底层的业务,在具备完整基准测试或者分析逻辑前,说快慢都是在猜测!
要关注的是解决问题的方法,解决这个问题,最简单,最干净的方法是什么??

*/
//模拟数据
type Data struct {
	SomeData string
	Line     string
}

//Xenia是我们要拿取数据的地方
type Xenia struct {
	Host    string        //主机
	Timeout time.Duration //超时
}

//Pull提供从Xenia获取数据的方法(模拟)
func (*Xenia) Pull(d *Data) error { //指针语义,我们想要向下分享
	switch rand.Intn(10) {
	case 1, 9:
		return io.EOF

	case 5:
		return errors.New("Error reading data from Xenia...")

	default:
		d.Line = "Data"
		fmt.Println("In:", d.Line)
		return nil

	}

}

//Pillar是要迁移过去的系统
type Pillar struct {
	Host    string        //主机
	Timeout time.Duration //超时

}

//Store提供向Pillar存入数据的方法
func (*Pillar) Store(d *Data) error { //指针语义,向下传递分享,尽量减小内存分配
	fmt.Println("Out:", d.Line)
	return nil

}

//---------------------以上解决了原始层4个小目标,假设我们完成了这些小api的单元测试-----------------
/* 二,Lower Level
此层代码
由原始层的两个函数表示,建立在原始代码层之上,这些功能专注于Xenia和Pillar的具体类型;接收数据的集合

*/
func Pull(x *Xenia, data []Data) (int, error) {
	for i := range data {
		if err := x.Pull(&data[i]); err != nil {
			return i, err
		}

	}

	return len(data), nil
}

func Store(p *Pillar, data []Data) (int, error) {
	for i := range data {
		if err := p.Store(&data[i]); err != nil {
			return i, err
		}
	}
	return len(data), nil

}

/* 三,high level
此层构建与PUll和Store之上,用于移动所有待处理数据

*/
func Copy(sys *System, batch int) error {
	data := make([]Data, batch)
	for {
		i, err := Pull(&sys.Xenia, data)
		if i > 0 {
			if _, err := Store(&sys.Pillar, data[:i]); err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}

}

//最初的想法是组成一个知道如何提取和储存的系统
type System struct {
	Xenia
	Pillar
}

func main() {
	sys := System{
		Xenia: Xenia{
			Host:    "localhost:8000",
			Timeout: time.Second,
		},
		Pillar: Pillar{
			Host:    "localhost:9000",
			Timeout: time.Second,
		},
	}
	if err := Copy(&sys, 3); err != io.EOF {
		fmt.Println(err)
	}
}

/* ----------至此,原型完成编写,假设实现基础功能----------- */

/* -------------------开始用接口解耦------------------ */

/* 一,首先要明白,在上面的程序中 哪些是可能会改变的;在这个例子中,最有可能会改变的是System自身,今天是Xenia Pillar,明天可能扩展更多的系统,为了应对这个变化,我们需要修改具体的函数 使之成为多态函数 */

//将Pull和Store的参数修改为,知道如何pull和store的事物(不管是什么类型)
/* func Pull(p Puller, data []Data) (int, error) {
	for i := range data {
		if err := p.Pull(&data[i]); err != nil {
			return i, err
		}
	}
	return len(data), nil
}
func Store(s Storer, data []Data) (int, error) {
	for i := range data {
		if err := s.Store(&data[i]); err != nil {
			return i, err
		}
	}
	return len(data), nil
}

type Puller interface {
 Pull(d *Data) error
}
type Storer interface {
 Store(d *Data) error
}

//只要同时知道Pull和Store的数具体数据都可作为参数
type PullStorer interface {
 Puller
 Storer
}

//永远传递的都不是接口本身的值(接口是没有意义的),而是储存在接口中的具体的值
func Copy(ps PullStorer, batch int) error {
 data := make([]Data, batch)
 for {
 i, err := Pull(ps, data)
 if i > 0 {
 if _, err := Store(ps, data[:i]); err != nil {
 return err
 }
 }
 if err != nil {
 return err
 }
 }
}

//现在System实际是实现了PullStorer接口的,在未来也可以构建不同的System

type System1 struct {
 Xenia
 Pillar
}
type System2 struct {
 Alice
 Bob
}

//至此,实现了应对不同的系统类型实现同样的功能,但是声明所有可能的System的组合形式是不可能的,高可维护性需要一个更好的解决方法

*/

/*
将System修改为接口
type System struct {
 Puller
 Storer
}
这样将使得System实现PullStorer接口,以任何可能的数据组合
至此,解耦完成,整个代码逻辑与System的类型组合方式将没有关系,只需要具备所需的方法即可


*/

/*
可改进的地方:
	1.PullerStorer接口可以删除,Copy的参数依旧改为*System
	func Copy(sys *System, batch int) error
	2.这样修改可使Copy函数更精准,并且同样多态,这样System的结构体可以删除
	func Copy(p Puller, s Storer, batch int) error
*/

//最终要明白的一点:抽象的目的不是模糊，而是创造一个新的语义层次，在这个层次上可以绝对精确


//------------------------------------------------------

