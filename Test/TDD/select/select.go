package sel

import (
	"fmt"
	"net/http"
	"time"
)

//同时发送请求,返回最快获得结果的,10s内都没返回 则返回错误

/*
httptest来模拟http请求


*/
/* //v1 存在以下几个问题:

func Racer(a, b string) (winner string) {
	A := time.Now()
	http.Get(a)
	aD := time.Since(A)

	B := time.Now()
	http.Get(b)
	bD := time.Since(B)

	if aD < bD {
		return a
	} else {
		return b

	}
}
//存在重复代码 进行重构:
*/

/* func Racer(a, b string) (winner string) {
	aDuration := measureResponseTime(a)
	bDuration := measureResponseTime(b)

	if aDuration < bDuration {
		return a
	}

	return b
}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
} */

//v3 使用select来提高效率
var tenSecondTimeout = 10 * time.Second

func Racer(a, b string) (winner string, error error) {
	return ConfigurableRacer(a, b, tenSecondTimeout)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		close(ch) //关闭后select将读取到0值
	}()
	return ch
}
