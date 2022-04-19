/* 必须要明确的一点,疯狂的记录详细的日志没有任何意义 */
package main

import (
	"fmt"

	"github.com/pkg/errors"
)

//自定义一个错误,实现error接口即可
type AppError struct {
	State int
}

func (c *AppError) Error() string {
	return fmt.Sprintf("App Error, State: %d", c.State)
}

//模拟错误的分层传递以及包装
func main() {
	if err := firstCall(10); err != nil {
		switch v := errors.Cause(err).(type) {
		case *AppError:
			fmt.Println("Custom App Error:", v.State)
		default:
			fmt.Println("Default Error")
		}
		fmt.Printf("%v\n", err)
	}
}

func firstCall(i int) error {
	if err := secondCall(i); err != nil {
		return errors.Wrapf(err, "secondCall(%d)", i)
	}
	return nil
}
func secondCall(i int) error {
	return &AppError{99}
}
