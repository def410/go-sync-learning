/*
	基于 sync.Mutex 实现一个可重入的互斥锁，只有拥有锁的 goroutine 能做解锁操作。
	通过 goroutine id 识别每一个 goroutine
*/

package lock

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type RecursiveMutex struct {
	sync.Mutex
	owner     int64 // 拥有锁的 goroutine 的 id
	recursion int32 // 拥有锁的 goroutine 的 重入次数
}

func (m *RecursiveMutex) Lock() {
	id := int64(goid())
	// 如果请求锁的 goroutine 是当前拥有锁的 goroutine，则只记录重入次数，并返回
	if atomic.LoadInt64(&m.owner) == id {
		m.recursion++
		return
	}
	// 否则当前请求锁的 goroutine 阻塞，直到获得锁
	m.Mutex.Lock()
	// 记录 goroutine id
	atomic.StoreInt64(&m.owner, id)
	// 重入次数设为 1
	m.recursion = 1
}

func (m *RecursiveMutex) Unlock() {
	id := int64(goid())
	// 当前 goroutine 未拥有锁，则不允许解锁
	if atomic.LoadInt64(&m.owner) != id {
		panic(fmt.Sprintf("RecursiveMutex: A goroutine that does not own the lock attempts to unlock"))
	}
	// 如果是拥有锁的 goroutine，则将重入次数减 1
	m.recursion--
	// 如果是最后一次释放锁，则清除属主信息,并释放锁
	if m.recursion == 0 {
		atomic.StoreInt64(&m.owner, -1)
		m.Mutex.Unlock()
	}
}

// 通过 runtime.Stack 方法获取当前 goroutine 的栈帧信息，取出其中包含的 goroutine id
func goid() int {
	buf := make([]byte, 32)
	l := runtime.Stack(buf, false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:l]), "goroutine "))
	id, err := strconv.Atoi(idField[0])
	if err != nil {
		panic(fmt.Sprintf("goid: can not to get the id of current goroutine: err"))
	}
	return id
}
