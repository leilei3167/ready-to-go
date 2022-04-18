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



 */

func main() {

}
