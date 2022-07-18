package poker

import (
	"fmt"
	"os"
	"time"
)

type BlindAlerter interface {
	ScheduledAlertAt(duration time.Duration, amout int)
}

//当暴露一个只有一个函数的接口时,惯例是同时暴露一个实现该接口的xxxFunc 函数类型,并创建方法使该函数类型实现该接口
//而在该方法中则是直接用接收者执行,这样就能够使的调用者,无需创建某个结构体来实现接口,而是直接传入一个函数(再强转为之前暴露的函数类型)

type BlindAlerterFunc func(duration time.Duration, amount int)

func (a BlindAlerterFunc) ScheduledAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

// StdOutAlerter 函数签名和BlindAlerterFunc一致,这样的化直接强转就可作为BlindAlerter接口的依赖注入实例使用
func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
	})
}
