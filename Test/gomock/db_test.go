package db

import (
	"errors"
	"github.com/stretchr/testify/require"
	mockdb "ready-to-go/Test/gomock/mock"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetFromDB(t *testing.T) {
	//创建mock控制器,创建mock的接口,用接口执行其某个方法,以及期望的返回值
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()         //断言方法是否被正确调用
	m := mockdb.NewMockDB(ctrl) //如果Get的参数为Tom 则返回错误
	m.EXPECT().Get(gomock.Eq("Tom")).Return(100, errors.New("not exist"))
	/*以上称之为打桩,有明确的参数和返回值是最简单的打桩方式
		参数:
		Eq(value) 表示与 value 等价的值。
		Any() 可以用来表示任意的入参。
		Not(value) 用来表示非 value 以外的值。
		Nil() 表示 None 值
		返回值:
		Return 返回确定的值
		Do  Mock方法被调用时，要执行的操作，忽略返回值。
		DoAndReturn 可以动态地控制返回值。
		调用次数(Times()):
		Times() 断言 Mock 方法被调用的次数。
		MaxTimes() 最大次数。
		MinTimes() 最小次数。
		AnyTimes() 任意次数（包括 0 次）。
	调用顺序(InOrder):

	*/

	//测试GetFromDB方法,将m作为DB传入
	v := GetFromDB(m, "Tom")
	require.Equal(t, v, -1)

	/*如何编写可Mock的代码?
		写可测试的代码和写好测试是同等重要的,要让代码可mock,要注意以下几点:
		1.mock 作用的是接口，因此将依赖抽象为接口，而不是直接依赖具体的类。
		2.不直接依赖的实例，而是使用依赖注入降低耦合性。

	PS:依赖注入:给予调用方它所需要的事物。 “依赖”是指可被方法调用的事物。依赖注入形式下，调用方不再直接指使用“依赖”，取而代之是“注入” 。
	“注入”是指将“依赖”传递给调用方的过程。在“注入”之后，调用方才会调用该“依赖”。传递依赖给调用方，而不是让让调用方直接获得依赖，
	这个是该设计的根本需求。

	如果GetFromDB方法中,未将DB作为参数传入,而在函数中新建,那就无法进行mock测试

	在BankDemo项目中,一开始是使用的传递Store结构体来构建Server来实现路由中与数据库交互的,但为了便于mock,将其
	改为了接口,在后续工作中最好将数据库交互的一系列方法抽象为接口,便于mock

	*/

}
