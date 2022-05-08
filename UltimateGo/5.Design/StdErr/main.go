package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

var Errtest = errors.New("invalid input!")

func main() {
	err := first()
	if errors.Is(err, Errtest) { //if err==Errtest 将返回false
		//errors.Is将检查整个错误链
		log.Println("equal!")
	} else {
		log.Println("not euqal!")
	}
	err = TestLogin()
	if e, ok := err.(*LoginError); ok && e.Err == Errtest {
		log.Println("断言成功!")
	}
}

func first() error {
	return fmt.Errorf("first:it is wrong %w", Errtest) //用%w来包裹错误,而不是破坏错误

}

type LoginError struct {
	LogTime  time.Time
	UserID   string
	TryTimes int
	Err      error
}

func (l *LoginError) Error() string {
	return fmt.Sprintf("%v:%v login err,have try %d times!:%v", l.LogTime, l.UserID, l.TryTimes, l.Err)
}

func TestLogin() error {
	var testlogin = LoginError{
		LogTime:  time.Now(),
		UserID:   "123",
		TryTimes: 5,
		Err:      Errtest,
	}
	return &testlogin
}
