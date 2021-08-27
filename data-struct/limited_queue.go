/*
	利用 sync.Cond 实现一个容量有限的 Queue
*/
package data_struct

import "sync"

type LimitedQueue struct {
	*sync.Cond
	s    []interface{}
	size int
}

func New(n int) LimitedQueue {
	return LimitedQueue{sync.NewCond(&sync.Mutex{}), make([]interface{}, 0, n), n}
}

// 追加元素到队列尾
func (q *LimitedQueue) Add(v interface{}) {
	q.L.Lock()
	// 队列满了就等待，唤醒后检查条件是否满足
	for len(q.s) == q.size { // 这里如果调用 q.Len() 会因重入而导致死锁
		q.Wait()
	}
	// 条件满足（有空位）则添加元素到队尾
	q.s = append(q.s, v)
	q.L.Unlock()
	// 通知其他 groutine 队列中数据有变化（元素 + 1）
	q.Cond.Broadcast()
}

// 从队列头取出元素
func (q *LimitedQueue) Pop() interface{} {
	q.L.Lock()
	// 队列为空则等待
	for len(q.s) == 0 {
		q.Wait()
	}
	// 有元素则取出
	v := q.s[0]
	q.s = q.s[1:]
	q.L.Unlock()
	// 通知其他 groutine 队列中数据有变化（空位 + 1）
	q.Cond.Broadcast()
	return v
}

func (q *LimitedQueue) Len() int {
	q.L.Lock()
	defer q.L.Unlock()
	return len(q.s)
}
