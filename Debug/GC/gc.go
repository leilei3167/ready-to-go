package main

import (
	"os"
	"runtime/trace"
)

func main() {
	trace.Start(os.Stderr)
	defer trace.Stop()

	keepalloc()

}

func keepalloc() {
	for i := 0; i < 100000; i++ {
		go func() {
			select {}
		}()

	}

}
