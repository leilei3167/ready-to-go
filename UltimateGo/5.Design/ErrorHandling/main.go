package main

import "log"

/* 完整性是软件工程的一个重要部分,而错误处理就是完整性的核心
在go中错误处理不是一个需要后面再处理的例外,而是一个需要主要关注的点

开发者有义务返回包含足够上下文的错误信息以便于用户能够作出正确的决定解决错误

处理一个错误要包含一下几方面:
	1.记录错误
	2.阻止传播错误
	3.决定是否要终止goroutine或者整个程序
*/

/* 在Go中 错误只是value,我们可以使他成为任何我们想要他们成为的东西,他们可以维持任何状态或者行为
type error interface {
 Error() string
}

内置的错误类型,虽然不可导出但是在任何地方只要实现了接口也能够被当成error类型使用

要记住,Error()string方法的实现是为了记录错误,和实现error接口!如果有任何的用户需要专门解析从这个方法返回的错误信息,那么是一个失败的错误处理!

以下是使用得最多的error value:注意关注指针语义是如何使用的
// http://golang.org/src/pkg/errors/errors.go
type errorString struct {
 s string
}
// http://golang.org/src/pkg/errors/errors.go
func (e *errorString) Error() string {
 return e.s}

 指针语义实现error接口
func New(text string) error {
 return &errorString{text}
}

*/

/* 上下文是错误处理的一切!
每一个错误必须提供足够的上下文,以使调用者能够对错误情况作出正确的决定

err!=nil 代表着,error接口如果储存着一个具体的类型,就说明发生了一个错误

if err := webCall(); err != nil {
 fmt.Println(err)
 return
}

func webCall() error {
 return New("bad request")
}

在这个例子中,上下文只是单单表示一个具体的错误存在,并不重点关注具体的错误值是什么

What if it’s important to know what error value exists inside the err interface
variable? Then error variables are a good option.
错误变量就会是一个好选择:

 错误变量都以Err开头,值为New创建的errorString
var (
 ErrBadRequest = errors.New("Bad Request")
 ErrPageMoved = errors.New("Page Moved")
)

func webCall(b bool) error {
 if b {
 return ErrBadRequest
 }
 return ErrPageMoved
}

明确定义了错误变量之后,将有助于调用者分辨出是哪些错误发生:
if err := webCall(true); err != nil {
 switch err {
 case ErrBadRequest:
 fmt.Println("Bad Request Occurred")
 return
 case ErrPageMoved:
 fmt.Println("The Page moved")
 return
 default:
 fmt.Println(err)
 return
 }
 }
 fmt.Println("Life is good")

 对于上述的例子,返回的上下文取决于返回的是哪一个错误变量
 那么如果一个错误变量包含的上下文信息不足怎么办?或者有些特殊的状态需要检查(如网络问题),在这种情况,就需要用户自定义具体的错误类型(仅需实现Error方法):

 type UnmarshalTypeError struct {
 Value string
 Type reflect.Type
}
func (e *UnmarshalTypeError) Error() string {
 return "json: cannot unmarshal " + e.Value +
 " into Go value of type " + e.Type.String()
}
注意到自定义的错误类型是有Error后缀的,同时注意指针语义的使用;再次强调,这个实现的目的是记录以及展示错误信息

自定义错误信息2(来自于标准库json包的实现)

type InvalidUnmarshalError struct {
 Type reflect.Type
}
func (e *InvalidUnmarshalError) Error() string {
 if e.Type == nil {
 return "json: Unmarshal(nil)"
 }
 if e.Type.Kind() != reflect.Ptr {
 return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
 }
 return "json: Unmarshal(nil " + e.Type.String() + ")"
}

关注如何返回的结构体;在return中,通过error接口 将其返回到调用处
func Unmarshal(data []byte, v interface{}) error {
 rv := reflect.ValueOf(v)
 if rv.Kind() != reflect.Ptr || rv.IsNil() {
 return &InvalidUnmarshalError{reflect.TypeOf(v)}
 }
 return &UnmarshalTypeError{"string", reflect.TypeOf(v)}
}

这里的上下文就更多是和错误接口内储存的值相关的信息,但是这种方式也需要一种取出错误的方法以便于测试(因为值存储在接口中,取出使用需要断言):
func main() {
 var u user
 err := Unmarshal([]byte(`{"name":"bill"}`), u)
 if err != nil {
 switch e := err.(type) {
 case *UnmarshalTypeError:
 fmt.Printf("UnmarshalTypeError: Value[%s] Type[%v]\n",
 e.Value, e.Type)
 case *InvalidUnmarshalError:
 fmt.Printf("InvalidUnmarshalError: Type[%v]\n", e.Type)
 default:
 fmt.Println(err)
 }
 return
 }
 fmt.Println("Name:", u.Name)
}
然而以上方法使得我的错误信息和代码耦合起来,我的错误信息一旦修改,整个代码将会崩溃,错误处理最优美的方式就是需要从重大改变中解耦出来

如果错误的实例拥有一个方法集,那么我就能够使用一个接口来应对类型查询;如net包中,有大量的错误类型实现了不同的方法;
有一个叫Temporary,他允许用户测试网络错误是严重的还是说它自己就能应付的
type temporary interface {
 Temporary() bool
}
func (c *client) BehaviorAsContext() {
 for {
 line, err := c.reader.ReadString('\n')
 if err != nil {
 switch e := err.(type) {
 case temporary:
 if !e.Temporary() {
 log.Println("Temporary: Client leaving chat")
 return
 }
 default:
 if err == io.EOF {
 log.Println("EOF: Client leaving chat")
 return
 }
 log.Println("read-routine", err)
 }
 }
 fmt.Println(line)
 }
}

*/

//永远操作error接口!
type customError struct{}

func (c *customError) Error() string {
	return "Find the bug."
}

func fail() ([]byte, *customError) { //返回值必须是error的接口而不是实现error接口的具体类型!
	return nil, nil
}

//fail返回的是nil,为什么会触发Fatal?
func main() {
	var err error
	if _, err = fail(); err != nil {
		/*  there is a nil pointer of type customError stored inside
		the err variable. That is not the same as a nil interface value of type error. */
		log.Printf("err is %T", err) //err is *main.customError
		log.Fatal("Why did this fail?")
	}
	log.Println("No Error")
}
