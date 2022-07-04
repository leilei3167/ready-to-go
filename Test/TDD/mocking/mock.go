package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

/*
假设你需要实现一个倒数计时的函数

虽然这是一个非常简单的函数,但是我们仍然需要使用迭代的,测试驱动的方式来实现他

迭代是什么意思?就是说 我们确保每次尽量实现最小的步骤,将需求拆分为尽量小的需求是
非常重要的技能

像这个需求我们可以这样进行拆分:
1.打印3
2.打印3,2,1 和Go
3.每一行打印间间隔一秒



*/
func main() {
	sleeper := &ConfigurableSleeper{1 * time.Second, time.Sleep}
	Countdown(os.Stdout, sleeper)
}

/* //v1
func Countdown(out io.Writer) {
	fmt.Fprintf(out, "3")
} */

//v2
const (
	finalWord      = "Go!"
	countdownStart = 3
)

/*
func Countdown(out io.Writer) {
	for i := countdownStart; i > 0; i-- {
		fmt.Fprintln(out, i)
		time.Sleep(time.Second)
	}
	fmt.Fprint(out, finalWord)
} */

/*直到v2版本,已经实现了所需要的功能,但是他仍然有几个大问题:
1.我们的单元测试要花费3s!!
  -在软件开发中,应该反复强调快速反馈的重要性
  -缓慢的测试破坏开发者的生产力
  -如果要增加更多的测试,每一次测试都需要花费3s,那将变得非常糟糕
2.我们还没有测试我们函数的一个重要属性

我们必须要提取对Sleeping的依赖,以便于我们可以在测试中控制它,我们可以使用mock模拟,然后通过依赖注入来使用他
而不是真正的依赖time.Sleep!!然后我们可以SPY on the calls(窥视这些调用),并且对其使用断言

*/

//v3 使函数依赖接口,而不是硬编码的time.Sleep
type Sleeper interface {
	Sleep()
}

type SpySleeper struct {
	Calls int
}

func (s *SpySleeper) Sleep() {
	s.Calls++ //到时候可以对其进行断言
}

//创建一个默认的Sleeper供主函数调用,而在测试函数中可以自己构建依赖(此处只是使计数器加一,而不是休眠)
//在有了configurableSleeper之后,我们可以安全删除
/* type DefaultSleeper struct{}

func (d *DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
} */

func Countdown(out io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		fmt.Fprintln(out, i)
		sleeper.Sleep()
	}

	fmt.Fprint(out, finalWord)
}

/*
v3版本仍然存在问题:
	1.我们只确保了他睡眠了3次,但是没有确保顺序!

	func Countdown(out io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		sleeper.Sleep()
	}

	for i := countdownStart; i > 0; i-- {
		fmt.Fprintln(out, i)
	}

	fmt.Fprint(out, finalWord)
}
以上代码也能通过单元测试!
*/

//v4
//实现Writer接口和Sleeper接口
type SpyCountdownOperations struct {
	Calls []string
}

func (s *SpyCountdownOperations) Sleep() {
	s.Calls = append(s.Calls, sleep)
}

func (s *SpyCountdownOperations) Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, write)
	return
}

const write = "write"
const sleep = "sleep"

//使得睡眠时间可以被配置

type ConfigurableSleeper struct {
	duration time.Duration
	sleep    func(time.Duration)
}

func (c *ConfigurableSleeper) Sleep() {
	c.sleep(c.duration)
}

type SpyTime struct {
	durationSlept time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
	s.durationSlept = duration
}

/*
如果mocking code变得非常复杂且有大量需要mock的东西,参考如下建议:
	1.你正在测试的东西承载了太多的功能,以至于他需要mock非常多的依赖
		-将这个功能模块再拆分得更小
	2.依赖过于细致
		-考虑如何将其中得一些依赖关系整合到一个有意义得模块中
	3.你的测试可能过分得关注细节了
		-测试应该更关注行为,而不是实现细节!
大多数时候,过多得需要mock得点往往是因为你使用了不恰当得抽象!

非常难以编写单元测试得代码往往设计得比较糟糕,设计良好得代码应该是非常易于编写测试的!

当你在重构某些部分时,如果你需要修改大量的测试函数,那就说明你在写测试时过于关注了细节,而不是行为,一定要
注意,关注行为而不是细节

以下是编写测试时需要遵循的过程和规则:
	1.重构的本质应该是 代码改变,但是最终的行为不改变;在进行一些重构时,理论上来说测试代码是不需要改变的,
	当你需要重构时,问问自己以下几个问题:
		-我是在测试我想要的行为,还是实现的细节?
		-如果我重构了这段代码,我是否需要对测试代码进行大量的修改?
	2.虽然Go允许测试私有函数,但是不建议这样做;因为私有函数实际就是公有函数的实现细节,尽量避免测试函数和
	私有函数相关联
	3.如果一个测试函数需要3个或以上的mock,那就是一个警告,请重新思考代码的设计!
	4.谨慎使用spy类型来监视测试代码,他将使得测试代码和代码实现紧密耦合;在需要关注实现细节时再使用
*/

/*
Gomock框架的使用:
https://geektutu.com/post/quick-gomock.html

测试中的名词:
http://xunitpatterns.com/Test%20Double.html

Test Double(即测试替身的类型):

Dummy objects: are passed around but never actually used. Usually they are just used to fill parameter lists.

Fake objects: actually have working implementations, but usually take some shortcut which makes them not suitable for production (an InMemoryTestDatabase is a good example).

Test Stubs: provide canned answers to calls made during the test, usually not responding at all to anything outside what's programmed in for the test.

Test Spies: are stubs that also record some information based on how they were called. One form of this might be an email service that records how many messages it was sent.

Mocks: are pre-programmed with expectations which form a specification of the calls they are expected to receive. They can throw an exception if they receive a call they don't expect and are checked during verification to ensure they got all the calls they were expecting.




*/
