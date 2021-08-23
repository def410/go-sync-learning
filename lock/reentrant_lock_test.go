package lock

import "testing"

func ReentrantMutexTest(t *testing.T) {
	m := &RecursiveMutex{}
	m.Lock()
	m.Lock()
	m.Lock()
	if m.recursion != 2 {
		t.Error("reentrant mutex: Wrong number of reentries")
	}
	gid := goid()
	if m.owner != int64(gid) {
		t.Error("reentrant mutex: owner record error")
	}
	m.Unlock()
	m.Unlock()
	m.Unlock()
	if m.recursion != 0 {
		t.Error("After the lock is released, the number of reentrants is not set to 0")
	}
	if m.owner != -1 {
		t.Error("After releasing the lock, the owner reset error")
	}
	
}
