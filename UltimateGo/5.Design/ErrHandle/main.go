/* 必须要明确的一点,疯狂的记录详细的日志没有任何意义
日志会造成大量的开销,日志的作用是帮助快速定位错误




*/
package main

import (
	"fmt"

	"github.com/pkg/errors"
)

//自定义一个错误,实现error接口即可
type AppError struct {
	State int
	Msg   string
}

func (c *AppError) Error() string {
	return fmt.Sprintf("App Error, State: %d,Msg:%s", c.State, c.Msg)
}

//模拟错误的分层传递以及包装
func main() {
	if err := firstCall(10); err != nil {
		//类型选择
		switch v := errors.Cause(err).(type) { //从error接口中选择其中存储的具体实例是什么?直接用err.(type)也是可行的
		case *AppError:
			fmt.Println("Custom App Error:", v.State) //Custom App Error: 99
		default:
			fmt.Println("Default Error")
		}
		fmt.Printf("%+v\n", err) //+v打印堆栈调用
	}
}

func firstCall(i int) error {
	if err := secondCall(i); err != nil {
		return errors.Wrapf(err, "secondCall(%d)err", i) //
	}
	return nil
}
func secondCall(i int) error {
	return &AppError{99, "SecondCall error!!!!!"}
}
