package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"
)

/* 6.1调度模式
当go程序启动时,runtime会向机器申请可以并行使用的线程数,这取决于程序可用核心数量,对于每一个可以并行使用的线程,runtime将会创建M,然后在之上绑定一个P,PM代表了计算运行能力和Go程序的执行上下文
此外 还会创建一个初始G来管理选定PM上指令的执行(Go调度器),就像M管理硬件上指令的执行一样,G管理M上指令的执行;这在操作系统之间创造了新的抽象层,将执行控制转移到了应用层级
*/

/* 6.2基础
runtime.GOMAXPROCX(1)设置Go程序为单线程执行,只会有一个M和P来执行所有的G,当设为0时默认和逻辑核心数相同
此外他同时是一个环境变量,想要在容器环境运行Go程序时需要考虑此项配置和容器环境进行匹配

 waitgroup
WaitGroup是一个解决并发编排的工具,一般情况下我们需要某个协程等待其他被自己开启的协程做完工作,当我们知道要开启的协程数量时,我们应该使用Add一次性加载完协程数量,在协程内部执行Done,当运行的是流式服务时,可以一次Add(1)一个,在主协程中调用main可使其阻塞等待,知道waitgroup计数器归零
不要将waitgroup当作参数在函数之间传递,要确保wg的add done在同一视线中,这样有助于减少bug和死锁
*/

/* 6.3抢占式调度 Preemptive Scheduler
了解是如何做到抢占式的,这意味着我们无法预测何时会发生上下文的切换,并且每次运行程序时都会发生变化:

func init() {
	runtime.GOMAXPROCS(1)

}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	//使每个协程做大量的工作,多到无法在一个时间片内完成的工作
	go func() {
		printHashes("A")
		wg.Done()
	}()
	go func() {
		printHashes("B")
		wg.Done()
	}()
	fmt.Println("Waiting To Finish")
	wg.Wait()
	fmt.Println("\nTerminating Program")

}

//做大量的IO工作,可能使其上下文切换
func printHashes(prefix string) {
	for i := 1; i <= 50000; i++ {
		num := strconv.Itoa(i)
		sum := sha1.Sum([]byte(num))
		fmt.Printf("%s: %05d: %x\n", prefix, i, sum)
	}
	fmt.Println("Completed", prefix)
}

调度器不会运行一个G占用太长的时间,当一个G执行系统调用时会执行上下文切换 让其他的G执行

*/

/* 6.4 DataRace
当多个G尝试同时接近同一个内存地址时,并且至少有一个G执行写入时,就会发生data race,当发生时,执行结果将变得不可靠,这类bug非常难以排查,因为datarace造成的结果总是随机的(并发调度的随机性)

6.5 DataRace 例子:

var counter int

func main() {
	const grs = 2
	var wg sync.WaitGroup
	wg.Add(grs)
	for g := 0; g < grs; g++ {
		//两个协程同时对counter进行修改,但能够得到正确结果4
		go func() {
			for i := 0; i < 2; i++ {
				value := counter
				value++//背后是一个读取,修改,写入的过程,操作系统很容易在其执行期间切换上下文
				//	log.Println("logging") //加入一个log却只能得到2这个错误结果,为什么?
				//因为log涉及到了系统调用,在错误的时间 调度器发生了上下文切换
				counter = value //并发写入
			}
			wg.Done()
		}()

	}

	wg.Wait()
	fmt.Println("Counter:", counter)
}
--------------------------------------------------------------
6.6 DATa Race检测:
run,build,test命令都可以进行竞态检测,但是要注意,使用build命令加 -race时将会使程序性能降低20%

$go build -race
但是一般-race会和test搭配更多一起使用,这里暂时先用build来演示
结果将会展示WARNING:DATA RACE 并打印追踪信息,会打印冲突的协程的信息,会非常便于定位发生race的位置

定位到了Race发生的地点,那么如何修复呢?可以使用两种工具:Atomic和Mutexes

---------------------------------------------------------------
6.7 Atomic
原子操作提供在硬件层面的同步,非常适用于计数器或者快速切换的机制,Wait Group就使用了原子操作
	原子操作只能在具体的 准确的数据类型上操作,如int32 int64等
	atomic包中的所有函数都是操作的地址,同步只会发生在地址级别,如果不同的协程调用相同的函数,但是操作的却是不同的地址的话,就无法同步
只需修改:
	var counter int =>var counter inte32
	counter++ =>atomic.AddInt32(&counter,1)
删除value := counter 	counter = value
--------------------------------------------------------------
6.8Mutexes
那么如果我并不想删除我的代码该怎么办呢?原子操作就不行了
使用互斥锁
调度器将同一时间将只会使得一个拥有锁的G执行;互斥锁不是一个队列!第一个
调用Lock的G不一定是第一个获得锁的G

非常重要的一点:使用锁会增加压力
确保在同一函数中上锁以及解锁,避免死锁的发生
错误示范,在同一函数内对同一mutexes多次解锁上锁,实际上根本就丧失了同步性,并且使得race也无法检测
go func() {
 for i := 0; i < 2; i++ {
 var value int
 mu.Lock() <-- Bad Use Of Mutex
 {
 value = counter
 }
 mu.Unlock()
 value++
 mu.Lock() <-- Bad Use Of Mutex
 {
 counter = value
 }
 mu.Unlock()
 }
 wg.Done()
 }()

 -----------------------------------------------------------
 6.9 RW Mutexes
第二种类型的锁是读写锁,因为读操作不会对数据竞态产生威胁(只有写会),读写锁将允许多个G同时读取同一块内存,当一个写请求发生时,就会停止读取,直到写入完成


var data []string
var rwMutex sync.RWMutex

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	//10个写入协程
	go func() {
		for i := 0; i < 10; i++ {
			writer(i)
		}
		wg.Done()
	}()
	//8个读取
	for i := 0; i < 8; i++ {
		go func(id int) {
			for {
				reader(id)
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("Program Complete")
}

//写入时
func writer(i int) {
	rwMutex.Lock() //锁定为写入状态,禁止其他协程读取和写入
	{
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		fmt.Println("****> : Performing Write")
		data = append(data, fmt.Sprintf("String: %d", i))
	}
	rwMutex.Unlock()
}

//Rlock锁定为读取状态,只其他线程允许读取 不允许写入
func reader(id int) {
	rwMutex.RLock()
	{
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		fmt.Printf("%d : Performing Read : Length[%d]\n", id, len(data))
	}
	rwMutex.RUnlock() //解锁,允许写入
}

---------------------------------------------------------------
6.10 Channel
非常重要的一个观点:不要把Channel当作数据结构,而是看作是一个发送信号通讯的机制;如果摆在我面前的问题不能用通讯解决,或者我没有提到信号两个字,那么就应该质疑,我是否该使用channel

当我们想到signaling时,应该将注意力放到3个方面:
	1.对于发送信号的G,是否需要确保他发送的信号已被接受?
		可能答案总是yes,但是要明白每个在信号层面的决定都会有成本的;此处的成本就是未知的延迟,sender不会知道它需要等待多久接收者才能收到信号,等待就会产生阻塞延迟,而且是未知时长的延迟

	2.我是否需要在信号中传输数据?如果要传输,那么传输就是一对一的,如果一个新的G也想要传输,那么就必须再发送第二个信号;如果数据不需要和信号一起发送,那么信号就可以是G之间的一对一或一对多的发送,不含数据的信号主要用于取消或者关闭,这些是由关闭channel实现的

	3.第三就是要考虑channel的状态;
	nil状态:读写nil都会阻塞

	make之后的channel可用状态,make的channel有两种形式:
		无缓冲chan:
			发送和接收方必须同时就绪
		有缓冲chan:
			缓冲没有满时,写入将不会阻塞;缓冲没有空时,读取将不会阻塞

	关闭状态:Close Channel并不是为了释放内存,而是为了改变状态;向关闭的chan写入将会panic,在关闭的chan读取将会立即返回

有了这些信息,我就可以专注于channel的模式,对于signal的关注非常重要!
	对于在信号层面是否需要确保接收方收到,取决于是否关注阻塞延迟;
	是否需要在信号传递数据,取决于是否处理取消,接下来将这些语法转换为语义

*/

/*
   6.11 Channel Patterns
   There are 7 Channel patterns that are important to understand since they provide the building blocks to signaling

*/

//6.11.1 Wait For Result
//此模式是后面的基础模式,G处理某些任务并且返回处理结果,允许使用带缓冲的chan存放工作;
func waitForResult() {
	//看作是带保证的信号级别
	ch := make(chan string)
	go func() {
		//假装处理了一些事,模拟一些随机的阻塞延迟
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		//返回结果
		ch <- "data"
		fmt.Println("child : sent signal")
	}()
	//阻塞获取结果,不知道会等待多久,所以阻塞延迟是未知的
	d := <-ch
	fmt.Println("parent : recv'd signal :", d) //绝对不要用打印来尝试确定运行的先后顺序
	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

//6.11.2 Fan Out/In
//运用了wait for result模式
//要记住此种模式在服务类应用中非常危险(如web),因为每个请求本身就是一个G,要特别注意规模
func fanOut() {
	//带缓冲的通道的缓冲尺寸必须经过谨慎考虑
	children := 2000
	ch := make(chan string, children)
	//此处不需要保证,因此使用带缓冲的chan,将有助于减少阻塞延迟(无缓冲必须等待处理完一个才能发送)

	//模拟生产2000个协程处理2000个某种工作
	for c := 0; c < children; c++ {
		go func(child int) {
			//模拟随机的处理时间
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			ch <- "data"
			fmt.Println("child : sent signal :", child)
		}(c)
	}

	//循环获取,使用计数器来确保每个协程都返回结果
	for children > 0 {
		d := <-ch
		children--
		fmt.Println(d)
		fmt.Println("parent : recv'd signal :", children)
	}
	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

//必须要记住fan out在运行的服务中是很危险的,假设一个http服务器能够用5wG处理5w个请求,而我假设要对其中10个G使用fan out来处理请求时,这些G的增长将是乘法级别增长的,大量的G可能会耗尽资源导致宕机

//6.11.3 Wait For Task 线程池的基础模式
func waitForTask() {
	//无缓冲Chan代表在信号层面有保证,意味着会有阻塞延迟
	ch := make(chan string)
	//worker不断的等待work
	go func() {
		d := <-ch //未知的阻塞时间
		fmt.Println("child : recv'd signal :", d)
	}()
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	//发布一个任务
	ch <- "data"
	fmt.Println("parent : sent signal")
	//和协程中的打印的执行时间是不确定的!!!!!!!!
	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

//6.11.4 Pooling 在Go中,池更注重资源使用的效率而不是操作系统那样提高cpu处理效率
func pooling() {
	ch := make(chan string)
	g := runtime.GOMAXPROCS(0)
	//开启和逻辑核心数(硬件线程)相同的协程处理任务(pool大小限制)
	for c := 0; c < g; c++ {
		go func(child int) {
			//使用for range可同时具备监听任务和监听退出信号的功能(一旦关闭则会break)
			for d := range ch {
				fmt.Printf("child %d : recv'd signal : %s\n", child, d)
			}
			fmt.Printf("child %d : recv'd shutdown signal\n", child)
		}(c)
	}

	//生产100个work,放入ch
	const work = 100
	for w := 0; w < work; w++ {
		ch <- "data"
		fmt.Println("parent : sent signal :", w)
	}
	//传入完成后关闭ch,G中的range将立刻能感知到
	close(ch)
	fmt.Println("parent : sent shutdown signal")
	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

//6.11.5 Drop 此模式在提供service时是非常重要的一个模式,它能在某个请求的负载过大,或者请求超出容量时丢弃请求
func drop() {
	const cap = 100 //容量
	ch := make(chan string, cap)

	go func() { //故意只设置一个worker,将使容量快速达到上限
		for p := range ch {
			fmt.Println("child : recv'd signal :", p)
		}
	}()

	const work = 2000
	//带有default的select永远不会阻塞
	for w := 0; w < work; w++ {
		select {
		case ch <- "data": //模拟放入work
			fmt.Println("parent : sent signal :", w)
			//实现Drop模式的关键点就在于default的设置
		default: //当容量满(即ch被阻塞时,将会执行丢弃或者其它想要设置的处理逻辑,比如返回一个500错误码或者服务正忙稍后重试等)
			fmt.Println("parent : dropped data :", w)
		}
	}

	close(ch)
	fmt.Println("parent : sent shutdown signal")
	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

//6.11.6 Cancellation
//用于告诉一个函数我愿意等待的时间
func cancellation() {
	rand.Seed(time.Now().UnixMicro())
	//设置带超时控制的context
	duration := 150 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	ch := make(chan string, 1)
	go func() {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		ch <- "data" //如果ch没有缓冲1,那么当主协程超时时,此协程将永久阻塞成为内存泄漏!
	}()

	select {
	case d := <-ch:
		fmt.Println("work complete", d)
	case <-ctx.Done(): //达到超时时间
		fmt.Println("work cancelled")
	}

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

//6.11.7 Fan Out/In Semaphore 信号量提供一个机制控制在某一时间同时工作的G的数量,但是仍然会为每一份工作创建G(不控制G的数量,但控制同时工作的G的数量)
func fanOutSem() {
	//2000个worker
	children := 2000
	ch := make(chan string, children)
	//同时工作的Worker限制在硬件线程数
	g := runtime.GOMAXPROCS(0)
	sem := make(chan bool, g) //这种形式会增加延时,仔细衡量取舍
	//开2000个worker
	for c := 0; c < children; c++ {
		go func(child int) {
			//执行的worker最大只有硬件线程数量个
			sem <- true
			{
				t := time.Duration(rand.Intn(200)) * time.Millisecond
				time.Sleep(t)
				ch <- "data"
				fmt.Println("child : sent signal :", child)
			}
			<-sem
		}(c)
	}

	for children > 0 {
		d := <-ch
		children--
		fmt.Println(d)
		fmt.Println("parent : recv'd signal :", children)
	}
	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

//6.11.8 Bounded Work Pooling 使用一个协程池来执行固定数量的work
func boundedWorkPooling() {
	//some works
	work := []string{"paper", "paper", "paper", "paper", 2000: "paper"}
	//number of workers
	g := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(g)
	ch := make(chan string, g)
	for c := 0; c < g; c++ {
		//开启逻辑核心数个worker
		go func(child int) {
			defer wg.Done()
			//监听ch
			for wrk := range ch {
				fmt.Printf("child %d : recv'd signal : %s\n", child, wrk)
			}
			fmt.Printf("child %d : recv'd shutdown signal\n", child)
		}(c)
	}
	//发送任务(生产者)
	for _, wrk := range work {
		ch <- wrk
	}
	close(ch)
	wg.Wait()
	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

//6.11.9 Retry TimeOut
//特别实用,比如我们在链接数据库,或者ping某个位置失败时,我们希望执行重试多少次
func retryTimeout(ctx context.Context, retryInterval time.Duration,
	check func(ctx context.Context) error) {

	for {
		//先执行用户要执行的,没有错误直接返回
		fmt.Println("perform user check call")
		if err := check(ctx); err == nil {
			fmt.Println("work finished successfully")
			return
		}
		//执行有错的话先检查ctx.Err 是否过期(即最外层 用户设置的WithTimeout)
		fmt.Println("check if timeout has expired")
		if ctx.Err() != nil {
			fmt.Println("time expired 1 :", ctx.Err())
			return
		}
		//设置计时器,同时监听外部ctx控制和时间脉冲
		fmt.Printf("wait %s before trying again\n", retryInterval)
		t := time.NewTimer(retryInterval)
		select {
		case <-ctx.Done():
			fmt.Println("timed expired 2 :", ctx.Err())
			t.Stop()
			return
		case <-t.C:
			fmt.Println("retry again") //计时器到则重新开始循环(重新执行)
		}
	}
}

//6.11.10 Channel Cancellation
func channelCancellation(stop <-chan struct{}) { //空结构体作为信号
	//创建能被cancle取消的ctx
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//开启协程,专门监听外来的stop信号和ctx的结束信号(因为本程序结束时此协程也应该结束)
	go func() {
		select {
		case <-stop: //发起者的关闭信号,关闭所有和ctx相关联的执行
			log.Println("收到调用者退出信号")
			cancel()
		case <-ctx.Done(): //自己创建的ctx的退出信号
			log.Println("任务完成执行退出...")
		}
	}()

	func(ctx context.Context) error {
		time.Sleep(time.Second * 5)
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"https://www.ardanlabs.com/blog/index.xml",
			nil,
		)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()

		return nil
	}(ctx)
}

func main() {
	//waitForResult()
	//fanOut()
	//drop()
	//cancellation()
	fanOutSem()
	//调用处设置一个最大的等待时长(期间内将会每秒间隔进行重试)
	/* Context */
	//context.BackGround 和Todo:都可以作为基本的父COntext,但是TODO可以表示后续增加某些控制(只是现在还没想清楚)
	//context的传递用的是值语义!! 也就是说函数间传递是复制,最大的原因是因为每个函数调用可能都会修改CTX(增加新功能),这样就能够实现在某条执行路径执行结束不会影响上游的执行

	ctx, cancle := context.WithTimeout(context.Background(), time.Second*3)
	defer cancle()
	retryTimeout(ctx, time.Second, func(ctx context.Context) error {
		conn, err := net.DialTimeout("tcp", "124.223.174.63:8080", time.Second*2)
		if err != nil {
			return errors.Wrap(err, "链接出错...")
		}
		fmt.Println("连接成功")
		conn.Close()
		return nil
	})

	//测试监听控制退出的G
	stop := make(chan struct{})
	go channelCancellation(stop)
	time.Sleep(time.Second * 3)
	stop <- struct{}{}
	time.Sleep(time.Second * 3)
}
