package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func main() {
	ctx := context.TODO()

	var (
		maxWorkers = int64(runtime.GOMAXPROCS(0))
		sem        = semaphore.NewWeighted(maxWorkers)
		out        = make([]int, 32)
	)
	//传入的是一个新的ctx,而不是和一开始的同一个
	//不要将此ctx传递给下游
	//WithContext创建的g在出现错误时会执行cancle
	group, _ := errgroup.WithContext(context.Background())
	for i := range out {
		//注意区分errgroup的ctx和acquire中的ctx并不是同一个
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to Acquire %v\n", err)
			break
		}
		group.Go(func() error {
			go func(i int) {

				defer sem.Release(1)
				out[i] = collatzSteps(i + 1)

			}(i)
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(out)
}

func collatzSteps(n int) (steps int) {
	if n <= 0 {
		panic("nonpositive input")
	}

	for ; n > 1; steps++ {
		if steps < 0 {
			panic("too many steps")
		}

		if n%2 == 0 {
			n /= 2
			continue
		}

		const maxInt = int(^uint(0) >> 1)
		if n > (maxInt-1)/3 {
			panic("overflow")
		}
		n = 3*n + 1
	}

	return steps
}
