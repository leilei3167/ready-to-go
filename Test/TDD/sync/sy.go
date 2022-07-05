package sy

import "sync"

type Counter struct {
	value int
	//mu使用值语义(使用mutex的地址将必须new)传递结构体的地址即可 不需要单独初始化互斥锁
	mu sync.Mutex //也有使用匿名嵌套的,但是注意匿名嵌套会使得其集成所有公开的方法(容易造成混乱,比如counter将会多出Lock等方法)
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}
func (c *Counter) Value() int {
	return c.value
}

func NewCounter() *Counter {
	return &Counter{}
}
