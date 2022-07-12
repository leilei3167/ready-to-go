package main

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"time"
)

func main() {
	if err := GetErr(); err != nil {
		if errors.Is(err, ErrTest) {
			fmt.Printf("Got Sentinal Err:%v\n", err)
		} else {
			log.Fatal("Got Err:", err)
		}
	}

	if err := GetTypeErr(); err != nil {

		var er *SomeError
		if errors.As(err, &er) {
			fmt.Printf("Got Type Err:%v\n", err)
		} else {
			log.Fatal("Got Err:", err)
		}

	}

}

var ErrTest = errors.New("test error")

type SomeError struct {
	Time   time.Time
	Caller string
}

func (e *SomeError) Error() string {
	return fmt.Sprintf("%v: %v", e.Time, e.Caller)
}

func GetErr() error {
	return fmt.Errorf("warp err: %w", ErrTest)
}

func GetTypeErr() error {
	return fmt.Errorf("warp type err:%w", &SomeError{Time: time.Now(), Caller: string(debug.Stack())})
}
