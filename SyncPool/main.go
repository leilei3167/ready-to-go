package main

import (
	"bytes"
	"io"
	"os"
	"sync"
	"time"
)

/*
使用Pool能够实现对象的复用,使得程序在运行过程中减少内存申请分配,极大的减轻GC的压力,在需要频繁创建
昂贵对象时特别的有用
sync.Pool 本质用途是增加临时对象的重用率，减少 GC 负担。划重点：临时对象。
所以说，像 socket 这种带状态的，长期有效的资源是不适合 Pool 的。

注意事项:
1.sync.Pool 本质用途是增加临时对象的重用率，减少 GC 负担;高负载时会动态扩容，存放在池中的对象如果不活跃了会被自动清理。；
2.不能对 Pool.Get 出来的对象做预判，有可能是新的（新分配的），有可能是旧的（之前人用过，然后 Put 进去的）；
3.不能对 Pool 池里的元素个数做假定，存在池子的对象可能在任何时刻被自动移除，我们对此不能做任何预期
4.sync.Pool 本身的 Get, Put 调用是并发安全的，sync.New 指向的初始化函数会并发调用，里面安不安全只有自己知道；
5.当用完一个从 Pool 取出的实例时候，一定要记得调用 Put，否则 Pool 无法复用这个实例，通常这个用 defer 完成；
*/

var bufPool = sync.Pool{
	New: func() any {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(bytes.Buffer)
	},
}

// timeNow is a fake version of time.Now for tests.
func timeNow() time.Time {
	return time.Unix(1136214245, 0)
}

func Log(w io.Writer, key, val string) {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset() //拿到对象时必须先清空
	// Replace this with time.Now() in a real logger.
	b.WriteString(timeNow().UTC().Format(time.RFC3339))
	b.WriteByte(' ')
	b.WriteString(key)
	b.WriteByte('=')
	b.WriteString(val)
	w.Write(b.Bytes())
	bufPool.Put(b)
}

func main() {
	Log(os.Stdout, "path", "/search?q=flowers")
}
