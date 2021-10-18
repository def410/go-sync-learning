/*
	有时我们希望在获取锁不成功时直接返回，而不是阻塞，要实现这种功能，可以基于 sync.Mutex 来扩展出一个 TryLock 方法。
*/

package lock

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// 这些变量在 monitor_mutex.go 中定义过
/*
const (
	mutexLocked = 1 << iota // 表示锁被持有
	mutexWoken // 标记是否有被唤醒的 goroutine
	mutexStarving // 标记是否处于饥饿状态
	mutexWaiterShift = iota // 标记 记录“等待者数量”的 bits 的起始位置
)
*/

type TryLockMutex struct {
	sync.Mutex
}

// 操纵 sync.Mutex 的 state 字段来实现 TryLock 方法
func (m *TryLockMutex) TryLock() bool {
	// 如果成功获得锁，则返回 true
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		return true
	}
	// 如果锁处于饥饿、唤醒、锁定状态，返回 false
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
		return false
	}
	// 否则当前的锁的状态为“未锁定”，并且有其他 goroutine 竞争这把锁。则当前 goroutine 也尝试竞争。
	new := old | mutexLocked
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), old, new)
}
