package main

//如何优雅的处理error?https://learnku.com/go/t/33210
/*

作者将error分为三种类型
1.Sentinel errors:
	如io.EOF,代表提前约定好的Err,使用起来不太灵活(必须使用==来确定是否为某个错误),直接通过Errorf添加上下文将使得底层Error被改变;
exm:
	err := readfile(“.bashrc”)
    if strings.Contains(error.Error(), "not found") {
        // handle error
    }
注:err.Error()函数是给程序员看的(用于打印),不要将其用于判断或其他

Sentinel errors 最大的问题在于它在定义 error 和使用 error 的包之间建立了依赖关系。
比如要想判断 err == io.EOF 就得引入 io 包，当然这是标准库的包，还 Ok。如果很多用户自定义的包都定义了错误，那我就要引入很多包，
来判断各种错误。麻烦来了，这容易引起循环引用的问题。因此,减少模仿标准库中的此类err的用法


2.Error Types
	指的是实现了Error接口的各种类型,它的一个重要的好处是，类型中除了 error 外，还可以附带其他字段，从而提供额外的信息，例如出错的行数等

exm:

	// PathError records an error and the operation and file path that caused it.
type PathError struct {
    Op   string
    Path string
    Err  error
}
	以上例子来自于os包,PathError能够在保存Err同时,附加上路径和操作名称等信息;使用这样的 error 类型,返回时直接返回实例的地址,外层调用者需要使用类型断言来判断错误：

exm:

	func underlyingError(err error) error {
    switch err := err.(type) {
    case *PathError: //类型选择,err中是否保存的是哪些实例
        return err.Err
    case *LinkError:
        return err.Err
    case *SyscallError:
        return err.Err
    }
    return err
}

但是这又不可避免地在定义错误和使用错误的包之间形成依赖关系，又回到了前面的问题。
即使 Error types 比 Sentinel errors 好一些，因为它能承载更多的上下文信息，但是它仍然存在引入包依赖的问题。因此，也是不推荐的。至少，不要把 Error types 作为一个导出类型。


3.Opaque errors
	即黑盒错误,或者不公开错误;你知道错误发生了,但是不知道其内部具体的发生细节

exm:

	func fn() error {
    x, err := bar.Foo()
    if err != nil {
        return err
    }

    // use x
    return nil
}

	在调用Foo时,如果发生错误,直接返回,否则继续执行;对于调用者,只关注调用成功与否;不成功则返回,成功则继续执行,无需关心错误具体信息!
当然，在某些情况下，这样做并不够用。例如，在一个网络请求中，需要调用者判断返回的错误类型，以此来决定是否重试。这种情况下，作者给出了一种方法：
就是说，不去判断错误的类型到底是什么，而是去判断错误是否具有某种行为，或者说实现了某个接口。

exm:

type temporary interface {
    Temporary() bool
}

func IsTemporary(err error) bool {
    te, ok := err.(temporary)
    return ok && te.Temporary()
}

发生错误时,调用IsTemporary,判断该Err是否具有某种行为,如果是,则就是某个错误

二,优雅的处理错误
使用pkg/errors的Warp来无损的添加上下文包裹错误,检查错误时使用Cause取出底层的错误类型,并用%+v来在最上层打印错误日志和堆栈追踪

三,错误只处理一次 Only handle errors once

	错误处理指的是:查看错误信息,并做出一个决定

exm:
  _, err := w.Write(buf)
    if err != nil {
        // annotated error goes to log file
        log.Println("unable to write:", err)

        // unannotated error returned to caller  return err
        return err
    }

以上例子是错误例子:进行了2次错误处理:1.打印错误日志 2.返回错误

第一次处理是将错误写进了日志，第二次处理则是将错误返回给上层调用者。而调用者也可能将错误写进日志或是继续返回给上层。这样就造成日志重复冗余;
正确做法应该是每一层用Warp添加上下文信息后直接返回

一定要记住,打印一定要在最上层进行,错误处理无论如何只进行一次;

总结:
errors 就像对外提供的 API 一样，需要认真对待。
将 errors 看成黑盒，判断它的行为，而不是类型。
尽量不要使用 sentinel errors。
使用第三方的错误包来包裹 error（errors.Wrap），使得它更好用。
使用 errors.Cause 来获取底层的错误。

而在go的标准库中,fmt.Errof使用%w关键字也可以无损的包裹error,并且使用Is或者As及Unwarp来获取错误链上的错误;
%w的实现其实是自定义错误类型的方式;
这里需要注意的是，嵌套可以有很多层，我们调用一次errors.Unwrap函数只能返回最外面的一层error，如果想获取更里面的，
需要调用多次errors.Unwrap函数。最终如果一个error不是warpping error，那么返回的是nil

As可用于替换原先的错误类型断言;Is可以用于替换类似 if err==...之类的判断

PS:官方包还是无法保留堆栈信息,因此还是使用pkg/errors比较好




*/
func main() {

}
