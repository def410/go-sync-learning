/*
	sync.Mutex 结构中的 state 字段是一个 int32 类型的值，它是一个复合字段，用不同的位表示多种信息。
	state 字段从右往左看：
		第 1 位：标记锁是否被持有
		第 2 位：标记是否有唤醒的 groutine
		第 3 位：标记是否处于饥饿状态
		剩余：标记等待当前锁的 groutine 数量
	sync.Mutex 并没有暴露这些信息，当它们有时很有用，比如我们想知道当前有多少 groutine 因这把锁被阻塞（监控等待者数量）、当前锁是否被持有等信息。
	因此我们可以自己实现这些功能，对 sync.Mutex 进行扩展。
*/

package lock

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// 定义与 sync 包中对应的常量，方便后续使用
const (
	mutexLocked      = 1 << iota // 表示锁被持有，值为 1（二进制）
	mutexWoken                   // 标记是否有被唤醒的 groutine，值为 10
	mutexStarving                // 标记是否处于饥饿状态
	mutexWaiterShift = iota      // 标记 记录“等待者数量”的 bits 的起始位置
)

type MonitorMutex struct {
	sync.Mutex
}

// 返回持有锁的 groutine 和 等待锁的 groutine 的总数
func (m *MonitorMutex) Count() int {
	// m.Mutex 在内存中起始位置，就是指向的就是 sync.Mutex 的 state 字段的指针
	n := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	// 右移，获取等待者数量
	n = n >> mutexWaiterShift
	// 如果当前的锁被持有，则 + 1
	n += n & mutexLocked
	return int(n)
}

// 判断锁是否处于某个状态，获取到的状态只是一个瞬态的值，函数返回前可能就发生了变化
func (m *MonitorMutex) Is(f func(n int32) bool) bool {
	n := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return f(n)
}

// 被持有？
func Locked(n int32) bool {
	n = n & mutexLocked
	return n == mutexLocked
}

// 有被唤醒的 groutine？
func Woken(n int32) bool {
	n = n & mutexWoken
	return n == mutexWoken
}

// 处于饥饿状态？
func Starving(n int32) bool {
	n = n & mutexStarving
	return n == mutexStarving
}
