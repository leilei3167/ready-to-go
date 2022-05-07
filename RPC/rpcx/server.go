package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/smallnest/rpcx/server"
)

type Test struct{}

type Arg struct {
	Name string
}
type Res struct {
	Msg string
}

func (t *Test) Try(ctx context.Context, arg *Arg, res *Res) error {

	for {
		select {
		case <-ctx.Done():
			log.Println("G:", runtime.NumGoroutine())
			res.Msg = "调用方超时退出!!"
			log.Println(res.Msg)

			return ctx.Err()
		default:
			for i := 0; i < 6; i++ {
				go dosomething(ctx, arg.Name+strconv.Itoa(i))

			}

		}
	}

}

func dosomething(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			//	log.Printf("退出!")
			return
		default:
			time.Sleep(time.Second)
			fmt.Println("working!!!")
		}

	}

}

func main() {
	s := server.NewServer()

	s.RegisterName("test", new(Test), "")

	s.Serve("tcp", ":8080")

}
