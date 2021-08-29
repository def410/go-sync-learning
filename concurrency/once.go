/*
	sync.Once 允许我们对传入的函数只做一次执行操作，可以实现单例模式、用于初始化资源等场景。
	但如果其中调用的函数产生错误，初始化资源失败，我们无从得知；且我们也无法知道资源是否已经初始化。
	因此，需要自己实现一个扩展的 Once
 */
package concurrency

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	done uint32
	sync.Mutex
}

func (o *Once) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) != 0 {
		return nil
	}
	return o.doSlow(f)
}

func (o *Once) doSlow(f func() error) error {
	o.Lock()
	defer o.Unlock()
	var err error
	if o.done == 0 { // 再次检查
		err = f()
		if err == nil { // 函数执行成功才更新 flag
			atomic.StoreUint32(&o.done, 1)
		}
	}
	return err
}

func (o *Once) Done() bool {
	return atomic.LoadUint32(&o.done) == 1
}

