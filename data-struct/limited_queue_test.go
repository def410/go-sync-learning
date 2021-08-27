package data_struct

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestLimitedQueue(t *testing.T) {
	q := New(5)
	for i := 0; i < 5; i++ {
		q.Add(i)
	}
	if q.Len() != 5 {
		t.Error("The number of elements in the queue is incorrect")
	}

	for i := 0; i < 5; i++ {
		q.Pop()
	}
	if q.Len() != 0 {
		t.Error("The number of elements in the queue is incorrect")
	}
}

func ExampleLimitedQueue() {
	q := New(3)
	fmt.Println(q.Len(), cap(q.s))

	q.Add(1)
	q.Add(2)
	q.Add(3)
	fmt.Println(q.Len(), cap(q.s))

	for i := 0; i < 3; i++ {
		fmt.Println(q.Pop())
	}

	// Output:
	// 0 3
	// 3 3
	// 1
	// 2
	// 3
}

func TestLimitedQueueInConcurrency(t *testing.T) {
	q := New(1)
	for i := 0; i < 5; i++ {
		go func(n int) {
			time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
			t.Logf("worker[%d] add %d", n, n)
			q.Add(n)
		}(i)
	}

	time.Sleep(1 * time.Second)
	for i := 0; i < 5; i++ {
		t.Logf("worker[-1] pop %d", q.Pop())
	}

}
